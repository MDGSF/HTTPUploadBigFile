package controllers

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MDGSF/HTTPUploadBigFile/modules/setting"
	"github.com/astaxie/beego"
)

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

	randNumber := rand.Int()
	randNumberStr := strconv.Itoa(randNumber)
	curTimeStr := time.Now().Format("2006-01-02_15-04-05")
	curFileName := curTimeStr + "_" + fileName + "_" + fileSize + "_" + randNumberStr
	curDir := filepath.Join(setting.RuntimeUploadDataDirectory, curFileName)
	err := os.MkdirAll(curDir, 0777)
	if err != nil {
		beego.Error(err)
		c.AjaxMsg(MSGERR, "inner error, create directory failed", http.StatusInternalServerError)
	}
	beego.Info("curDir =", curDir)

	c.AjaxMap(MSGOK, "upload big file initialize success", map[string]interface{}{
		"file_directory": curDir,
	})
}

// UploadOneChunk 接收 chunk
func (c *TUploadBigFileController) UploadOneChunk() {
	c.Ctx.Request.ParseMultipartForm(32 << 20)
	r := c.Ctx.Request

	beego.Info("TUploadBigFileController, c.Ctx.Input.Params() =", c.Ctx.Input.Params())
	beego.Info("TUploadBigFileController, c.Ctx.Request.Form =", c.Ctx.Request.Form)
	beego.Info("TUploadBigFileController, c.Ctx.Request.PostForm =", c.Ctx.Request.PostForm)
	beego.Info("TUploadBigFileController, len(r.MultipartForm.File) =", len(r.MultipartForm.File))
	for k := range r.MultipartForm.File {
		beego.Info("k =", k)
	}

	fileName := c.GetParameterString("file_name")
	chunkIndex := c.GetParameterString("chunk_index")
	curDir := c.GetParameterString("file_directory")
	curFileName := fileName + "_" + chunkIndex

	chunkfileHeaders := r.MultipartForm.File["chunk_data"]
	chunkfileHeader := chunkfileHeaders[0]
	chunkfile, err := chunkfileHeader.Open()
	if err != nil {
		beego.Error(err)
		return
	}
	defer chunkfile.Close()

	localFilePath := filepath.Join(curDir, curFileName)
	localFile, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error(err)
		return
	}
	defer localFile.Close()
	io.Copy(localFile, chunkfile)

	c.AjaxMsg(MSGOK, "upload chunk success")
}
