package task

import (
	"encoding/base64"
	"fmt"
	"iagent/models"
	"iagent/util"
	"net"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
)

func runInstall() {
	var cache util.Cache
	val, err := cache.Get()
	if err != nil {
		//fmt.Println(err)
		return
	}
	if val.Count > 0 && val.Data != nil {
		//fmt.Println(val.Count)
		for _, v := range val.Data {
			connect := "未知"
			var status, err string
			var val_all, val1, val2, val3 string
			var err1, err2, err3 error

			val1, err1 = util.ClientCmd(v.IP, v.Info.Port, v.Info.User, v.Info.Pwd, "mkdir -p /tmp/.zabbix/")
			val_all = ""
			if err1 != nil {
				fmt.Println(err1)
				goto ADD
			} else {
				val_all = val1 + "\n"
				//fmt.Println(val1)
			}

			val2, err2 = util.ClientScp(v.IP, v.Info.Port, v.Info.User, v.Info.Pwd, "./rpm/zabbix-agent.sh", "/tmp/.zabbix/zabbix-agent.sh")
			if err2 != nil {
				fmt.Println(err2)
				goto ADD
			} else {
				val_all = val_all + val2 + "\n"
				//fmt.Println(val1)
			}
			//fmt.Println("err2:------------")

			val3, err3 = util.ClientCmd(v.IP, v.Info.Port, v.Info.User, v.Info.Pwd, "bash /tmp/.zabbix/zabbix-agent.sh install")
			//fmt.Println("val3:\t" + val3)
			if err3 != nil {
				fmt.Println(err3)
				goto ADD
			} else {
				val_all = val_all + val3
				//fmt.Println(val3)
				goto ADD
			}

		ADD:
			//ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
			//dial tcp 124.0.0.1:22: i/o timeout
			//1.认证失败 2.无法连接 3.安装已执行
			if strings.Contains(err, "unable to authenticate") {
				status = "认证失败"
			} else if strings.Contains(err, "i/o timeout") {
				status = "无法连接"
			} else {
				status = "脚本已执行"
			}

			//fmt.Println(v.IP, v.Info.Port, v.Info.User, v.Info.Pwd, val_all, status, connect)
			//fmt.Println(v.IP + "\t运行错误:\t" + err)
			AddAgentLog(v.IP, v.Info.Port, v.Info.User, v.Info.Pwd, val_all, status, connect)
			fmt.Println(v.IP + "\t安装已执行完毕")
		}
		UpdateAgentStatus()
	}
	//v1, _ := json.Marshal(val)
	//fmt.Println(string(v1))
}
func Docron() {
	s := gocron.NewScheduler()
	s.Every(5).Seconds().Do(runInstall)
	s.Every(5).Minutes().Do(UpdateAgentStatus)
	<-s.Start()
}

func AddAgentLog(ip string, port int, user string, pwd string, logtext string, status string, connect string) {
	var agent models.AgentLog
	agent.Ip = ip
	agent.Port = port
	agent.User = base64.StdEncoding.EncodeToString([]byte(user))
	agent.Password = base64.StdEncoding.EncodeToString([]byte(pwd))
	agent.Log = logtext
	agent.Status = status
	agent.Connect = connect
	agent.AddAgentLog()
}

func UpdateAgentStatus() {
	fmt.Println("定时任务运行:\t" + "UpdateAgentStatus")
	var agent models.AgentLog
	ips, err := agent.GetAgentLog()
	if err != nil {
		return
	}
	//检测所有IP，定期修改其tcp端口10050的连接状态
	for _, ip := range ips {
		getAgent := agent.GetAgentLogByIp(ip.Ip)
		if CheckIpPort(ip.Ip, 10050) {
			//已经为连接状态，不执行update
			if getAgent.Connect == "连接" {
				break
			}
			getAgent.Connect = "连接"
			getAgent.UpdateAgentLog()
		}
		//端口无法连接，但数据中记录为连接，故现在状态为离线
		if !CheckIpPort(ip.Ip, 10050) && getAgent.Connect == "连接" {
			getAgent.Connect = "离线"
			getAgent.UpdateAgentLog()
		}
	}

}

func CheckIpPort(ip string, port int) bool {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer conn.Close()
	return true
}
