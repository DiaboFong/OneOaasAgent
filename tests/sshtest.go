package main

import (
	"fmt"
	"iagent/util"
)

func main() {
	host := "10.211.55.9"
	user := "root"
	sshKeyPath := "/Users/itnihao/.ssh/id_rsa"
	//sshKeyStr := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsJTpS72alxVUYvKRJWIvF3DCRH+hu+lehqr/0Vcyj0SQkJfeWkNLTYz1GAYQwPNWqJ5MBCut9T28BrqZGQai4Ey6fWr6g6UX0yD1AH1b/XGab+K38ym2Vje8BzvrGA20iRpCqo2VBho3cBSZVSG7Pi2okmNrmOi6H2DSrjMGb4kSpPluvzTG+1A5gxxpv0cIGcvYPuJvXB2qE/Dx1UEwWX+93+5QolmSRyzAl+5oIphGLqaXPacb27++x+jO6Jg9G1FmyvJpS8n2wK6zc/vFWziIil/ftMX8Ordzn12Duyj8hLrVQFcOlJ1or5R/AXSF3i32OOs1yapYRpX8z0Ohx itnihao@itnihao.local"
	//var somebody person

	client := util.SshClient{
		IP:   host,
		User: user,
		Port: 22,
		//CertStr: sshKeyStr,
		CertPath: sshKeyPath,
	}
	var cmdlist = [2]string{""}
	s := cmdlist[:1]
	s = append(s, "ls /", "ls /etc")
	for _, i := range s {
		client.Connect()
		out := client.RunCmd(i)
		client.Close()
		fmt.Println(string(out))
		fmt.Println("------------------------------------")
	}

}
