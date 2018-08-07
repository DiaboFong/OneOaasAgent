package controllers

import (
	"encoding/base64"
	"fmt"
	"iagent/models"
	"iagent/util"
	"strings"
)

type InstallAgentController struct {
	BaseController
}

// @Title 安装agent
// @Description 安装agent
// @Success 200 {map} map[string]interface{}
// @Param	agentip		query 	string	true		"agent的ip地址"
// @Param	username		query 	string	true		"用户"
// @Param	password		query 	string	true		"密码"
// @Param	islist		query 	string	true		"是否为列表"
// @Failure 400 添加用户失败
// @router /install [post]
func (i *InstallAgentController) InstallAgent() {
	/* 功能:安装agent
	   1.以单个ip，用户，密码安装agent
	   2.以多个ip，每个ip对应的用户，密码进行安装
	*/
	var ips []string
	ip := strings.TrimSpace(i.GetString("ip", ""))
	port, _ := i.GetInt("port", 22)
	username, err := base64.StdEncoding.DecodeString(i.GetString("username", ""))
	if err != nil {
		i.ResponseJSON(400, "用户名格式不正确......")
	}
	password, err := base64.StdEncoding.DecodeString(i.GetString("password", ""))
	if err != nil {
		i.ResponseJSON(400, "用户名格式不正确......")
	}

	if strings.Contains(ip, ",") {
		ips = strings.Split(ip, ",")
	} else {
		ips = append(ips, ip)
	}

	//将待安装IP放入到缓存中去，实现队列的功能
	var cachedata util.CacheData
	for _, ip := range ips {
		if ip == "" {
			continue
		}
		var ipmap util.IpMap
		ipmap.IP = ip
		ipmap.Info.Port = port
		ipmap.Info.User = strings.TrimSpace(string(username))
		ipmap.Info.Pwd = strings.TrimSpace(string(password))
		cachedata.Data = append(cachedata.Data, ipmap)
	}
	cachedata.Count = len(cachedata.Data)

	var cache util.Cache
	err = cache.Set(cachedata)
	if err != nil {
		fmt.Println(err)
		panic("设置存储失败,请检查数据类型")
	}
	//v, err := cache.Get()
	//v1, _ := json.Marshal(v)
	i.ResponseJSON(200, "任务已全部下发成功,等待执行结果......")
}

// @Title 卸载agent
// @Description 卸载agent
// @Success 200 {map} map[string]interface{}
// @Param	ip		query 	string	true		"agent的ip地址"
// @Failure 400 添加用户失败
// @router /remove [post]
func (i *InstallAgentController) RemoveAgent() {
	var agent models.AgentLog
	ip := strings.TrimSpace(i.GetString("ip", ""))
	/*var ips []string
	if strings.Contains(ip, ",") {
		ips = strings.Split(ip, ",")
	} else {
		ips = append(ips, ip)
	}*/
	if ip == "" {
		i.ResponseJSON(400, "IP地址不允许为空......")
		return
	}
	v := agent.GetAgentLogByIp(ip)
	user1, err1 := base64.StdEncoding.DecodeString(v.User)
	pwd1, err2 := base64.StdEncoding.DecodeString(v.Password)

	if err1 != nil {
		i.ResponseJSON(400, "内部错误:\tuser数据类型不正确......"+err1.Error())
		return
	}

	if err2 != nil {
		i.ResponseJSON(400, "内部错误:\tpassword数据类型不正确......"+err2.Error())
		return
	}
	user := string(user1)
	pwd := string(pwd1)
	val, err := util.ClientCmd(v.Ip, v.Port, user, pwd, "bash /tmp/.zabbix/zabbix-agent.sh remove")
	agent.DelAgentLog(v.Ip)
	if err != nil {
		i.ResponseJSON(400, string(val)+"卸载任务命令执行失败......"+err.Error())
		return
	} else {
		i.ResponseJSON(200, string(val)+"卸载任务命令执行成功......")
	}
}

func (i *InstallAgentController) GetAgentStatus() {
	//获取agent的状态
}

// @Title  获取安装日志
// @Description 获取安装日志
// @Success 200            {map} map[string]interface{}
// @Failure 400 获取日志失败
// @router /log [get]
func (i *InstallAgentController) GetAgentLog() {
	var agetnlog models.AgentLog

	logs, err := agetnlog.GetAgentLog()
	if err != nil {
		i.ResponseJSON(400, "获取日志失败")
		return
	}
	ret := map[string]interface{}{
		"content": logs,
	}
	i.ResponseJSON(200, "操作成功", ret)
}

// @Title  从IP获取安装日志
// @Description 从IP获取安装日志
// @Success 200            {map} map[string]interface{}
// @Failure 400 获取日志失败
// @router /logbyip [get]
func (i *InstallAgentController) GetAgentLogByIp() {
	ip := strings.TrimSpace(i.GetString("ip", ""))
	if ip == "" {
		i.ResponseJSON(400, "ip不能为空")
		return
	}
	var agetnlog models.AgentLog
	agetnlog = agetnlog.GetAgentLogByIp(ip)
	ret := map[string]interface{}{
		"content": agetnlog,
	}
	i.ResponseJSON(200, "操作成功", ret)
}

func (i *InstallAgentController) GetCount() {
	var agentlog models.AgentLog

	today := agentlog.GetAgentToday()
	count := agentlog.GetAgentCount()
	connect := agentlog.GetAgentConnect()
	fail := count - connect

	ret := map[string]interface{}{
		"today":   today,
		"count":   count,
		"connect": connect,
		"fail":    fail,
	}
	i.ResponseJSON(200, "获取数据成功", ret)

}
