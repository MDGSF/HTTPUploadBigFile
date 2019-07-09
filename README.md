# HTTPUploadBigFile

upload big file, html5, http, file slice

支持上传大文件，测试 50G 文件稳定上传，测试环境：`chrome+ubuntu16.04`。

## 使用

```
$ go run main.go
2019/07/09 15:40:47.143 [I] [asm_amd64.s:1337]  http server Running on http://127.0.0.1:8080
```

访问网页： [http://127.0.0.1:8080](http://127.0.0.1:8080) ，选择一个大文件上传。

## 配置文件

`conf/app.conf`

```
appname = HTTPUploadBigFile
httpaddr = 127.0.0.1
httpport = 8080
runmode = dev
RuntimeUploadDataDirectory = /tmp/httpbigfile/runtime/uploaddata
```

`RuntimeUploadDataDirectory` 是上传之后的文件存放路径。
