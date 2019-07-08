package routers

import (
	"github.com/MDGSF/HTTPUploadBigFile/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/UploadBigFile", &controllers.TUploadBigFileController{})
}
