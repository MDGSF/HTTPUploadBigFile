package routers

import (
	"github.com/MDGSF/HTTPUploadBigFile/controllers"
	"github.com/astaxie/beego"
)

func init() {
	uploadBigFileCtrl := &controllers.TUploadBigFileController{}
	beego.Router("/", uploadBigFileCtrl, "get:GetUploadPage")
	beego.Router("/api/v1/UploadBigFileInit", uploadBigFileCtrl, "post:UploadBigFileInit")
	beego.Router("/api/v1/UploadBigFileChunk", uploadBigFileCtrl, "post:UploadOneChunk")
}
