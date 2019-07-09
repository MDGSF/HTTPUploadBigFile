package setting

import (
	"log"
	"os"

	"github.com/astaxie/beego/config"
)

var (
	// RuntimeUploadDataDirectory 存放上传文件的目录
	RuntimeUploadDataDirectory string
)

// LoadConfig 加载配置文件
func LoadConfig(filename string) {
	cnf, err := config.NewConfig("ini", filename)
	if err != nil {
		log.Fatal(err)
	}

	RuntimeUploadDataDirectory = cnf.DefaultString("RuntimeUploadDataDirectory", "/tmp/httpbigfile/runtime/uploaddata")
	os.MkdirAll(RuntimeUploadDataDirectory, 0777)
}
