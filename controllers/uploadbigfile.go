package controllers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/MDGSF/HTTPUploadBigFile/modules/setting"
	"github.com/astaxie/beego"
)

// TUploadBigFileController 大文件上传控制器
type TUploadBigFileController struct {
	beego.Controller
}

// Post 接收 chunk
func (c *TUploadBigFileController) Post() {
	c.Ctx.Request.ParseMultipartForm(32 << 20)
	r := c.Ctx.Request

	beego.Info("TUploadBigFileController, c.Ctx.Input.Params() =", c.Ctx.Input.Params())
	beego.Info("TUploadBigFileController, c.Ctx.Request.Form =", c.Ctx.Request.Form)
	beego.Info("TUploadBigFileController, c.Ctx.Request.PostForm =", c.Ctx.Request.PostForm)
	beego.Info("TUploadBigFileController, len(r.MultipartForm.File) =", len(r.MultipartForm.File))
	for k := range r.MultipartForm.File {
		beego.Info("k =", k)
	}

	chunkfileHeaders := r.MultipartForm.File["chunk_data"]
	chunkfileHeader := chunkfileHeaders[0]
	chunkfile, err := chunkfileHeader.Open()
	if err != nil {
		beego.Error(err)
		return
	}
	defer chunkfile.Close()

	localFilePath := filepath.Join(
		setting.RuntimeUploadDataDirectory,
		c.Ctx.Request.Form.Get("file_name")+"_"+c.Ctx.Request.Form.Get("chunk_index"))
	localFile, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error(err)
		return
	}
	defer localFile.Close()
	io.Copy(localFile, chunkfile)

	out := make(map[string]interface{})
	out["errno"] = 200
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}
