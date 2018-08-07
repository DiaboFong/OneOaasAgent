package controllers

//基础服务类

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/astaxie/beego"
)

// 400   Bad Request（错误请求） 服务器不理解请求的语法。
// 401   Unauthorized（未授权） 请求要求身份验证。 对于需要登录的网页，服务器可能返回此响应。
// 403   Forbidden（禁止） 服务器拒绝请求。
// 404   Not Found（未找到） 服务器找不到请求的网页。
//
// 500   Internal Server Error（服务器内部错误）  服务器遇到错误，无法完成请求。
// 501   Not Implemented（尚未实施） 服务器不具备完成请求的功能。
// 502   Bad Gateway（错误网关） 服务器作为网关或代理，从上游服务器收到无效响应。
// 504   Gateway Timeout（网关超时）  服务器作为网关或代理，但是没有及时从上游服务器收到请求。

var (
	Finish200                   = errors.New("请求成功")
	ErrorBadRequest400          = errors.New("请求不正确")
	ErrorBadJson400             = errors.New("请求的Json格式不正确")
	ErrorBadParam400            = errors.New("请求的参数未通过验证")
	ErrorUnauthorized401        = errors.New("该请求未授权")
	ErrorForbidde403            = errors.New("该请求被禁止")
	ErrorDataNotFound404        = errors.New("服务器找不到对应的数据")
	ErrorPageNotFound404        = errors.New("服务器找不到该请求")
	ErrorInternalServerError500 = errors.New("服务器内部发生错误，请联系管理员")
	ErrorNotImplemented501      = errors.New("该功能未实现")
	ErrorBadGateway502          = errors.New("终端网关未响应")
	ErrorGatewayTimeout504      = errors.New("终端网关超时")
)

type JsonResult struct {
	Code    int         `json:"code"`
	Content interface{} `json:"content,omitempty"`
	Msg     string      `json:"msg"`
}

var tokenSecrityKey = "nR5cCI6IkpXVCJ9.eoJhZG1pbi"

//控制器公共结构
type BaseController struct {
	beego.Controller
	ControllerName string
	ActionName     string
	UserId         int64
	UserName       string
	PageSize       int
}

/*
func (this *BaseController) Prepare() {
	this.EnableRender = false
	util.InitLog()

	controllerName, actionName := this.GetControllerAndAction()
	this.ControllerName = strings.ToLower(controllerName[0:len(controllerName)])
	this.ActionName = strings.ToLower(actionName)

	// skip
	if controllerName == "AuthController" && actionName == "Login" {
		return
	}
	if controllerName == "EventController" && actionName == "AddEvent" {
		return
	}

	// 认证
	if auth := this.Auth(); !auth {
		util.ConsoleLog.Error("API请求认证失败")
        		this.ResponseErrorJSON(401, ErrorUnauthorized_401.Error())
		return
	}
}*/

//登录状态验证
/*func (this *BaseController) Auth() bool {
	authorizationString := this.Ctx.Request.Header.Get("Authorization")
	if len(authorizationString) == 0 {
		return false
	}

	userid, passed := models.CheckAuthentication(authorizationString)
	this.UserId = userid

	return passed
}*/

// 输出json
func (ctx *BaseController) JsonResult(out interface{}) {
	ctx.Data["json"] = out
	ctx.ServeJSON()
	ctx.StopRun()
}

// 处理响应信息
func (ctx *BaseController) ResponseJSON(code int, msg string, other ...map[string]interface{}) {
	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	if other != nil && len(other) > 0 {
		for k, v := range other[0] {
			if k == "code" || k == "msg" {
				continue
			}
			response[k] = v
		}
	}
	ctx.Data["json"] = response
	ctx.ServeJSON()
}

func (ctx *BaseController) DataJSON(arry []interface{}) {
	response := map[string]interface{}{
		"data": arry,
	}
	ctx.Data["json"] = response
	ctx.ServeJSON()
}

//处理响应错误信息
func (ctx *BaseController) ResponseDataJSON(code int, msg string, content interface{}) {
	response := map[string]interface{}{
		"code":    code,
		"msg":     msg,
		"content": content,
	}

	ctx.Data["json"] = response
	ctx.ServeJSON()
}

//处理响应错误信息
func (ctx *BaseController) ResponseErrorJSON(code int, msg string, other ...map[string]interface{}) {
	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	if other != nil && len(other) > 0 {
		for k, v := range other[0] {
			if k == "code" || k == "msg" {
				continue
			}
			response[k] = v
		}
	}
	ctx.Data["json"] = response
	ctx.ServeJSON()
}

//处理请求响应,以JSON格式返回响应
func (ctx *BaseController) ResponseSuccessJSON(result bool, msg string, id int64) map[string]interface{} {
	return map[string]interface{}{
		"result": result,
		"msg":    msg,
		"id":     id,
	}
}

func (ctx *BaseController) ResponseSuccessJSONEx(result bool, msg string, path string) map[string]interface{} {
	return map[string]interface{}{
		"result": result,
		"msg":    msg,
		"path":   path,
	}
}

// 处理响应JSON格式
func (ctx *BaseController) TotalRowsJSON(arry []interface{}) {
	response := map[string]interface{}{
		"total": 100,
		"rows":  arry,
	}
	ctx.Data["json"] = response
	ctx.ServeJSON()
}

func (ctx *BaseController) TotalRowsPageJSON(count int64, arry []interface{}) {
	var response map[string]interface{}
	response = map[string]interface{}{
		"total": count,
		"rows":  arry,
	}
	ctx.Data["json"] = response
	ctx.ServeJSON()
}

// 渲染模版
func (ctx *BaseController) Display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "error/error.html"
	}
	ctx.TplName = tplname
}

// 重定向
func (ctx *BaseController) GoRedirect(url string) {
	ctx.Redirect(url, 302)
}

// 是否POST提交
func (ctx *BaseController) IsPost() bool {
	return ctx.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (ctx *BaseController) GetClientIp() string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(ctx.Ctx.Request.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	//[::1]:61721 localhost有IPV6会被解析，正常情况为127.0.0.1:61856
	s := strings.Split(ctx.Ctx.Request.RemoteAddr, ":")
	if len(s) == 4 {
		return "[::1]"
	}
	return s[0]
}

func SizeFormat(size float64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB"}
	n := 0
	for size > 1024 {
		size /= 1024
		n += 1
	}

	return fmt.Sprintf("%.2f %s", size, units[n])
}
