package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
)

const (
	MSGOK  = 0
	MSGERR = -1
)

// Controller base controller
type Controller struct {
	beego.Controller
	IsLogin      bool
	IsRoot       bool
	IsSuperAdmin bool
	IsAdmin      bool
}

// NestPreparer exec NestPrepare after base controller Prepare
type NestPreparer interface {
	NestPrepare()
}

// Prepare exec before Get, post, ...
func (c *Controller) Prepare() {

	c.Ctx.Request.ParseMultipartForm(32 << 20)

	if c.Ctx.Request.Form == nil {
		c.Ctx.Request.ParseForm()
	}

	c.StartSession()

	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

// GetParameterInt 获取客户端 int 参数
func (c *Controller) GetParameterInt(key string) int {
	value, err := c.GetInt(key)
	if err != nil {
		c.AjaxMsg(MSGERR, "invalid parameter "+key, http.StatusBadRequest)
	}
	return value
}

// GetParameterInts 获取客户端 int 参数
func (c *Controller) GetParameterInts(key string) []int {
	strValue := c.GetParameterStrings(key)
	result := make([]int, 0)
	for _, v := range strValue {
		iValue, err := strconv.Atoi(v)
		if err != nil {
			c.AjaxMsg(MSGERR, "invalid parameter "+key, http.StatusBadRequest)
		}
		result = append(result, iValue)
	}
	return result
}

// GetParameterString 获取客户端 string 参数
func (c *Controller) GetParameterString(key string) string {
	value := c.GetString(key)
	if len(value) == 0 {
		c.AjaxMsg(MSGERR, "invalid parameter "+key, http.StatusBadRequest)
	}
	return value
}

// GetParameterStrings 获取客户端 string 参数
func (c *Controller) GetParameterStrings(key string) []string {
	value := c.GetStrings(key)
	if len(value) == 0 {
		c.AjaxMsg(MSGERR, "invalid parameter "+key, http.StatusBadRequest)
	}
	return value
}

// AjaxMsg 返回
func (c *Controller) AjaxMsg(msgno int, msg interface{}, arg ...interface{}) {
	out := make(map[string]interface{})
	if msgno == MSGOK {
		out["success"] = true
	} else {
		out["success"] = false
	}
	out["result"] = msg
	c.Data["json"] = out

	if len(arg) > 0 {
		code := arg[0].(int)
		c.Ctx.ResponseWriter.WriteHeader(code)
	}

	c.ServeJSON()
	c.StopRun()
}

// AjaxList 返回 列表
func (c *Controller) AjaxList(msgno int, msg interface{}, count int, data interface{}, arg ...interface{}) {
	out := make(map[string]interface{})
	if msgno == MSGOK {
		out["success"] = true
	} else {
		out["success"] = false
	}
	out["result"] = msg
	out["count"] = count
	out["data"] = data
	c.Data["json"] = out

	if len(arg) > 0 {
		code := arg[0].(int)
		c.Ctx.ResponseWriter.WriteHeader(code)
	}

	c.ServeJSON()
	c.StopRun()
}

// AjaxMap 返回
func (c *Controller) AjaxMap(msgno int, msg interface{}, m map[string]interface{}, arg ...interface{}) {
	out := make(map[string]interface{})
	if msgno == MSGOK {
		out["success"] = true
	} else {
		out["success"] = false
	}
	out["result"] = msg

	for k, v := range m {
		out[k] = v
	}

	c.Data["json"] = out

	if len(arg) > 0 {
		code := arg[0].(int)
		c.Ctx.ResponseWriter.WriteHeader(code)
	}

	c.ServeJSON()
	c.StopRun()
}

// ServeDownloadFile 处理下载请求
func (c *Controller) ServeDownloadFile(fileName, filePathName string) {
	w := c.Ctx.ResponseWriter
	r := c.Ctx.Request
	//w.Header().Add("Content-Disposition", "attachment; filename="+url.QueryEscape(packagename+suffixName))
	w.Header().Add("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Add("Content-Description", "File Transfer")
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Transfer-Encoding", "binary")
	w.Header().Add("Expires", "0")
	w.Header().Add("Cache-Control", "must-revalidate")
	w.Header().Add("Pragma", "public")
	http.ServeFile(w, r, filePathName)
}
