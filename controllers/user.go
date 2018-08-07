package controllers

import (
	"iagent/models"
	"iagent/util"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	BaseController
}

// @Title CreateUser 添加用户
// @Description create users 添加用户
// @Success 200 {map} map[string]interface{}
// @Param	username		query 	string	true		"用户名"
// @Param	password		query 	string	true		"密码"
// @Failure 400 添加用户失败
// @router / [post]
func (u *UserController) AddUser() {
	var user models.User
	username := strings.TrimSpace(u.GetString("username", ""))
	password := strings.TrimSpace(u.GetString("password", ""))
	user.Enable = true
	if username == "admin" {
		u.ResponseJSON(400, "参数错误,username不允许为admin")
		return
	}

	if username == "" {
		u.ResponseJSON(400, "参数错误,username不允许为空")
		return
	}

	if password == "" {
		u.ResponseJSON(400, "参数错误,password不允许为空")
		return
	}

	user.UserName = username
	user.Salt = util.GenerateSalt()
	user.Password = util.Strtomd5(util.GenrateHash(user.Salt, util.Strtomd5(password)))

	_, err := user.AddUser()
	if err != nil {
		if err.Error() == "用户名已经存在" {
			u.ResponseJSON(400, "参数错误,username不允许重复")
			return
		}
		u.ResponseJSON(500, ErrorInternalServerError500.Error())
		return
	}

	u.ResponseJSON(200, "添加用户成功")
}

// @Title  修改用户
// @Description update users 添加用户
// @Success 200            {map} map[string]interface{}
// @Param	userid		    query 	int64	true		"用户名"
// @Param	username		query 	string	true		"用户名"
// @Param	password		query 	string	true		"密码"
// @Failure 400 修改用户失败
// @router / [put]
func (u *UserController) UpdateUser() {
	var user models.User
	userid, err := u.GetInt64("userid", 0)
	username := u.GetString("username", "")
	password := u.GetString("password", "")
	sessionid := u.Ctx.GetCookie("sessionid")
	if err != nil {
		u.ResponseJSON(400, "userid参数不存在")
		return
	}
	getoneuser, err := user.GetUserById(userid)
	if sessionid != getoneuser.Sessionid {
		u.ResponseJSON(400, "session非法")
		return
	}
	//通过userid从数据库中取出的username 与API传入的username不一致
	if username != getoneuser.UserName {
		u.ResponseJSON(400, "userid和username不匹配")
		return
	}
	//API传入的username与数据库中的username不一致，不允许修改username
	if username != getoneuser.UserName {
		u.ResponseJSON(400, "参数错误,username不允许修改")
		return
	}

	//可以修改密码
	if password != "" {
		user.Password = util.Strtomd5(util.GenrateHash(getoneuser.Salt, util.Strtomd5(password)))
		u.ResponseJSON(400, "参数错误,password不允许为空")
		return
	}
	err = user.UpdateUser()
	if err != nil {
		u.ResponseJSON(400, "更新失败")
	}
}

// @Title  删除用户
// @Description update users 添加用户
// @Success 200            {map} map[string]interface{}
// @Param	userid		    query 	int64	true		"用户名"
// @Param	username		query 	string	true		"用户名"
// @Param	password		query 	string	true		"密码"
// @Failure 400 修改用户失败
// @router / [delete]
func (u *UserController) DelUser() {
	var user models.User
	userid, err := u.GetInt64("userid", 0)
	if err != nil {
		u.ResponseJSON(400, "userid参数不存在")
		return
	}

	getoneuser, err := user.GetUserById(userid)
	sessionid := u.Ctx.GetCookie("sessionid")
	if sessionid != getoneuser.Sessionid {
		u.ResponseJSON(400, "session非法")
		return
	}
	if err != nil || getoneuser.Userid == 0 {
		u.ResponseJSON(400, "传入userid不存在")
		return
	}

	err = user.DelUser(userid)
	if err != nil {
		u.ResponseJSON(400, "更新失败")
	} else {
		u.ResponseJSON(200, "删除用户成功")
	}
}

// @Title  获取所有用户
// @Description get all users 获取所有用户
// @Success 200            {map} map[string]interface{}
// @Failure 400 修改用户失败
// @router / [get]
func (u *UserController) GetAllUsers() {
	var user models.User
	users := user.GetUser()

	ret := map[string]interface{}{
		"content": users,
	}
	u.ResponseJSON(200, "操作成功", ret)
}

// @Title  通过用户id获取用户
// @Description get users from uid 通过用户id获取用户
// @Success 200            {map} map[string]interface{}
// @Param	userid		    query 	int64	true		"用户id"
// @Failure 400 该id的用户不存在
// @router /uid/:id [get]
func (u *UserController) GetUserById() {
	var user models.User
	userid, _ := strconv.ParseInt(u.Ctx.Input.Param(":userid"), 10, 64)
	getoneuser, err := user.GetUserById(userid)

	sessionid := u.Ctx.GetCookie("sessionid")
	if sessionid != getoneuser.Sessionid {
		u.ResponseJSON(400, "session非法")
		return
	}

	if err != nil {
		u.ResponseJSON(400, "该id的用户不存在")
		return
	}
	ret := map[string]interface{}{
		"content": getoneuser,
	}

	u.ResponseJSON(200, "操作成功", ret)
}

// @Title  用户登录
// @Description 用户登录
// @Success 200            {map} map[string]interface{}
// @Param	username		    query 	string	true		"用户名称"
// @Param	password		    query 	string	true		"用户密码"
// @Failure 400 用户名和密码不正确
// @router /login [post]
func (u *UserController) Login() {
	var user models.User
	username := u.GetString("username", "")
	password := u.GetString("password", "")
	if username == "" || password == "" {
		u.ResponseJSON(400, "用户名和密码不能为空")
		return
	}

	getoneuser := user.GetUserByName(username)
	sessionid := u.Ctx.GetCookie("sessionid")
	if sessionid == getoneuser.Sessionid && sessionid != "" {
		u.Ctx.Redirect(200, "/")
		return
	}
	if getoneuser.Password != util.Strtomd5(util.GenrateHash(getoneuser.Salt, util.Strtomd5(password))) {
		u.ResponseJSON(400, "用户名和密码不正确")
		return
	} else {
		getoneuser.Sessionid = util.GenerateSessionId()
		getoneuser.LoginIp = u.GetClientIp()
		getoneuser.LoginTime = time.Now().Unix()
	}
	if getoneuser.UpdateUser() != nil {
		u.ResponseJSON(400, "登录失败")
		return
	} else {
		ret := map[string]interface{}{
			"sessionid":  getoneuser.Sessionid,
			"login_ip":   getoneuser.LoginIp,
			"login_time": getoneuser.LoginTime,
		}
		u.Ctx.SetCookie("sessionid", getoneuser.Sessionid, 7*24*3600)
		u.ResponseJSON(200, "登录成功", ret)
		//u.Ctx.Redirect(200, "/")
	}
}
