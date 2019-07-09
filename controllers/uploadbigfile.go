package controllers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/MDGSF/HTTPUploadBigFile/modules/setting"
	"github.com/astaxie/beego"
)

// GUploadFileMap 保存每一个大文件上传过程中需要的数据
var GUploadFileMap sync.Map

// TBigFileContext 大文件上下文
type TBigFileContext struct {
	// FileName 文件名
	FileName string

	// FileSize 文件大小，单位：字节
	FileSize int64

	// FileDirectory 保存 chunk 的目录
	FileDirectory string

	// SuccessReceivedChunkNum 成功收到的 chunk 数量
	SuccessReceivedChunkNum int
}

// NewBigFileContext 新建一个大文件上下文
func NewBigFileContext(
	fileName string,
	fileSize int64,
	fileDirectory string,
) *TBigFileContext {
	bigFileCtx := &TBigFileContext{}
	bigFileCtx.FileName = fileName
	bigFileCtx.FileSize = fileSize
	bigFileCtx.FileDirectory = fileDirectory
	return bigFileCtx
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TUploadBigFileController 大文件上传控制器
type TUploadBigFileController struct {
	Controller
}

// GetUploadPage 获取上传页面
func (c *TUploadBigFileController) GetUploadPage() {
	c.TplName = "index.tpl"
}

// UploadBigFileInit 上传大文件初始化
func (c *TUploadBigFileController) UploadBigFileInit() {
	fileName := c.GetParameterString("file_name")
	fileSize := c.GetParameterString("file_size")
	iFileSize := c.GetParameterInt64("file_size")

	randNumber := rand.Int()
	randNumberStr := strconv.Itoa(randNumber)
	curTimeStr := time.Now().Format("2006-01-02_15-04-05")
	curFileName := curTimeStr + "_" + fileName + "_" + fileSize + "_" + randNumberStr
	fileDirectory := filepath.Join(setting.RuntimeUploadDataDirectory, curFileName)
	err := os.MkdirAll(fileDirectory, 0777)
	if err != nil {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, create directory failed", http.StatusInternalServerError)
	}
	beego.Info("fileDirectory =", fileDirectory)

	GUploadFileMap.Store(fileDirectory, NewBigFileContext(
		fileName, iFileSize, fileDirectory,
	))

	c.AjaxMap(MSGOK, "upload big file initialize success", map[string]interface{}{
		"file_directory": fileDirectory,
	})
}

// UploadOneChunk 接收 chunk
func (c *TUploadBigFileController) UploadOneChunk() {
	r := c.Ctx.Request

	// beego.Info("TUploadBigFileController, c.Ctx.Input.Params() =", c.Ctx.Input.Params())
	// beego.Info("TUploadBigFileController, c.Ctx.Request.Form =", c.Ctx.Request.Form)
	// beego.Info("TUploadBigFileController, c.Ctx.Request.PostForm =", c.Ctx.Request.PostForm)
	// beego.Info("TUploadBigFileController, len(r.MultipartForm.File) =", len(r.MultipartForm.File))
	// for k := range r.MultipartForm.File {
	// 	beego.Info("k =", k)
	// }

	fileName := c.GetParameterString("file_name")
	fileSize := c.GetParameterInt64("file_size")
	chunkIndex := c.GetParameterInt("chunk_index")
	chunkSize := c.GetParameterInt64("chunk_size")
	chunkTotalNumber := c.GetParameterInt("file_chunk_total_number")
	fileDirectory := c.GetParameterString("file_directory")
	maxZeroPaddingNumber := len(c.GetParameterString("file_chunk_total_number"))
	fileNameFormatStr := "%s_%0" + strconv.Itoa(maxZeroPaddingNumber) + "d"
	curFileName := fmt.Sprintf(fileNameFormatStr, fileName, chunkIndex)

	beego.Info(fmt.Sprintf("UploadOneChunk, chunkIndex = %v, chunkTotalNumber = %v, chunkSize = %v, file_directory = %v",
		chunkIndex, chunkTotalNumber, chunkSize, fileDirectory))

	if len(r.MultipartForm.File) == 0 {
		beego.Error("len(r.MultipartForm.File) == 0")
		c.AjaxMsg(MSGERR, "invalid request, no chunk file data", http.StatusBadRequest)
	}
	chunkfileHeaders, ok := r.MultipartForm.File["chunk_data"]
	if !ok {
		c.AjaxMsg(MSGERR, "invalid request, no chunk_data", http.StatusBadRequest)
	}
	chunkfileHeader := chunkfileHeaders[0]
	chunkfile, err := chunkfileHeader.Open()
	if err != nil {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, open chunk file failed", http.StatusInternalServerError)
	}
	defer chunkfile.Close()

	localFilePath := filepath.Join(fileDirectory, curFileName)
	localFile, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, create local file failed", http.StatusInternalServerError)
	}
	defer localFile.Close()
	if written, err := io.Copy(localFile, chunkfile); err != nil || written != chunkSize {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, copy data to local file failed", http.StatusInternalServerError)
	}

	bigFileCtxI, ok := GUploadFileMap.Load(fileDirectory)
	if !ok {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, get context from global failed", http.StatusInternalServerError)
	}

	bigFileCtx := bigFileCtxI.(*TBigFileContext)
	bigFileCtx.SuccessReceivedChunkNum++

	if bigFileCtx.SuccessReceivedChunkNum == chunkTotalNumber {
		go MergeChunkToOneFile(
			fileDirectory,
			fileName,
			fileSize,
			chunkTotalNumber,
			fileNameFormatStr,
		)
	}

	c.AjaxMsg(MSGOK, fmt.Sprintf("upload chunk %v success", chunkIndex))
}

/*
MergeChunkToOneFile 把多个 chunk 合并为一个文件
@param fileDirectory: chunk 所在的目录
@param fileName: 文件名
@param fileSize: 文件大小，单位：字节
@param chunkTotalNumber: 一共有多少个 chunk
@param fileNameFormatStr: chunk 文件名格式
*/
func MergeChunkToOneFile(
	fileDirectory string,
	fileName string,
	fileSize int64,
	chunkTotalNumber int,
	fileNameFormatStr string,
) {
	bigFilePath := filepath.Join(fileDirectory, fileName)
	bigFile, err := os.OpenFile(bigFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Info(fmt.Sprintf("start merge file, %v", bigFilePath))

	for chunkIndex := 1; chunkIndex <= chunkTotalNumber; chunkIndex++ {
		chunkFileName := fmt.Sprintf(fileNameFormatStr, fileName, chunkIndex)
		chunkFilePath := filepath.Join(fileDirectory, chunkFileName)
		chunkFile, err := os.OpenFile(chunkFilePath, os.O_RDONLY, 0666)
		if err != nil {
			beego.Error(err)
			return
		}

		beego.Info(fmt.Sprintf("merging chunk %v/%v to file %v...", chunkIndex, chunkTotalNumber, fileName))

		if _, err := io.Copy(bigFile, chunkFile); err != nil {
			beego.Error(err)
			return
		}

		chunkFile.Close()
	}

	bigFile.Close()

	bigFileInfo, err := os.Stat(bigFilePath)
	if err != nil {
		beego.Error(err)
		return
	}

	if bigFileInfo.Size() != fileSize {
		beego.Error(fmt.Sprintf("bigFileInfo.Size() = %v, fileSize = %v", bigFileInfo.Size(), fileSize))
		return
	}

	GUploadFileMap.Delete(fileDirectory)

	beego.Info(fmt.Sprintf("merge file success, %v", bigFilePath))
}
