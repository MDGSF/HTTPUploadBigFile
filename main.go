package main

import (
	"github.com/MDGSF/HTTPUploadBigFile/modules/setting"
	_ "github.com/MDGSF/HTTPUploadBigFile/routers"
	"github.com/astaxie/beego"
)

func main() {
	setting.LoadConfig("conf/app.conf")
	beego.Run()
}
