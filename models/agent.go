package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego"
)

type AgentLog struct {
	Id         int64  `json:"id" orm:"pk;column(id);auto"`
	Ip         string `json:"ip" orm:"column(ip);size(50);null"`
	Log        string `json:"log" orm:"column(log);type(longtext);null" `
	UpdateTime int64  `json:"update_time" orm:"column(update_time);null" `
	CreateTime int64  `json:"create_time" orm:"column(create_time);null" `
	Status     string `json:"status" orm:"column(status);size(50);null"`
	Connect    string `json:"connect" orm:"column(connect);size(50);null"`
	User       string `json:"-" orm:"column(user);size(50);null"`
	Password   string `json:"-" orm:"column(password);size(256);null"`
	Port       int    `json:"-" orm:"column(port);size(6);null"`
}

func (a *AgentLog) TableEngine() string {
	return "INNODB"
}

func (a *AgentLog) TableName() string {
	return "agent_log"
}

func (a *AgentLog) GetAgentLogByIp(ip string) AgentLog {
	var log AgentLog
	//查询一个，不能使用eamil *Group
	err := Orm.QueryTable("agent_log").Filter("ip", ip).One(&log)
	if err != nil {
		beego.Debug("查询agent_log错误", err.Error())
	}
	return log
}
func (a *AgentLog) GetAgentToday() int64 {

	nTime := time.Now()
	//ToDay := nTime.Format("20060102 00:00:00")
	tStr := nTime.Format("20060102")
	t, _ := time.Parse(tStr, "20060102")
	count, err := Orm.QueryTable("agent_log").Filter("create_time__gt", t.Unix()).Count()
	if err != nil {
		return 0
	}
	return count
}

func (a *AgentLog) GetAgentConnect() int64 {
	count, err := Orm.QueryTable("agent_log").Filter("connect", "连接").Count()
	if err != nil {
		return 0
	}
	return count
}

func (a *AgentLog) GetAgentCount() int64 {
	count, err := Orm.QueryTable("agent_log").Count()
	if err != nil {
		return 0
	}
	return count
}
func (a *AgentLog) AddAgentLog() (int64, error) {
	a.UpdateTime = time.Now().Unix()
	a.CreateTime = time.Now().Unix()
	count, err := Orm.QueryTable("agent_log").Filter("ip", a.Ip).Count()
	if err != nil && count != 0 {
		return 0, err
	}
	if count == 1 {
		getid := a.GetAgentLogByIp(a.Ip)
		a.Id = getid.Id
		a.UpdateAgentLog()
		return 0, nil
	}
	if count > 1 {
		a.DelAgentLog(a.Ip)
	}
	id, err := Orm.Insert(a)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AgentLog) DelAgentLog(ip string) error {
	if ip == "" {
		return errors.New("删除的条件:IP地址不允许为空")
	}
	_, err := Orm.QueryTable("agent_log").Filter("ip", ip).Delete()
	if err != nil {
		beego.Error("删除IP的安装日志错误", err.Error())
		return err
	}
	return nil
}

func (a *AgentLog) UpdateAgentLog() (err error) {
	_, err = Orm.Update(a)
	if err != nil {
		return err
	}
	return nil
}

func (a *AgentLog) GetAgentLog() ([]*AgentLog, error) {
	var agentlog []*AgentLog
	_, err := Orm.QueryTable("agent_log").All(&agentlog)
	if err != nil {
		beego.Debug("【查询agentlog错误", err.Error(), "]")
	}
	//beego.Debug("【查询agentlog的记录数据是", num, "]条")

	return agentlog, nil
}
