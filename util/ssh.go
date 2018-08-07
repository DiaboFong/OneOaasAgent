package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

const defaultTimeout = 10 // second
type SshClient struct {
	IP       string
	Port     int
	User     string
	CertPath string
	CertStr  string
	Password string
	session  *ssh.Session
	client   *ssh.Client
}

func (sshclient *SshClient) readPublickKeyFile(file string) ssh.AuthMethod {
	var b []byte
	if sshclient.CertPath != "" {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil
		}
		b = f
	}
	if sshclient.CertPath == "" {
		b = []byte(sshclient.CertStr)
	}

	key, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

func (sshclient *SshClient) Connect() error {
	var sshConfig *ssh.ClientConfig
	var auth []ssh.AuthMethod
	if sshclient.CertPath != "" || sshclient.CertStr != "" {
		auth = []ssh.AuthMethod{sshclient.readPublickKeyFile(sshclient.CertPath)}
	} else {
		auth = []ssh.AuthMethod{ssh.Password(sshclient.Password)}
	}

	sshConfig = &ssh.ClientConfig{
		User: sshclient.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * defaultTimeout,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshclient.IP, sshclient.Port), sshConfig)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}

	sshclient.session = session
	sshclient.client = client
	return nil
}

func (sshclient *SshClient) Scp(src string, dst string) error {
	err := scp.CopyPath(src, dst, sshclient.session)
	if err != nil {
		return err
	}
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("no such file or directory: %s", src))
	} else {
		return nil
	}
}
func (sshclient *SshClient) RunCmd(cmd string) []byte {
	out, err := sshclient.session.Output(cmd)
	if err != nil {
		fmt.Println("运行命令出现错误了--------------------------------")
		fmt.Println(err, out)
	}
	//defer sshclient.Close()
	return out
}

func (sshclient *SshClient) Close() {
	sshclient.session.Close()
	sshclient.client.Close()
}

func ClientScp(ip string, port int, user string, pwd string, src string, dst string) (string, error) {
	client := SshClient{
		IP:       ip,
		User:     user,
		Port:     port,
		Password: pwd,
	}
	err := client.Connect()
	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	err = client.Scp(src, dst)
	if err != nil {
		//fmt.Println(err)
		return "", err
	} else {
		return fmt.Sprintf(ip + "\tcopy文件:\t" + src + "\t成功"), nil
	}
	client.Close()
	return "", nil
}

func ClientCmd(ip string, port int, user string, pwd string, cmd string) (string, error) {
	client := SshClient{
		IP:       ip,
		User:     user,
		Port:     port,
		Password: pwd,
	}
	err := client.Connect()
	if err != nil {
		return "", err
	}
	out := client.RunCmd(cmd)
	fmt.Println(ip + "\t当前运行的命令是\t" + cmd)
	client.Close()
	return string(out), nil
}
