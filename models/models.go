package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	//"oneoaas.com/oneoaas_event/util"
)

var (
	//public
	Orm orm.Ormer

	//db conf
	dbtype string
	dbuser string
	dbpass string
	dbhost string
	dbport string
	dbname string
	dsn    string

	//mysql
	maxIdle int
	maxConn int
)

//注册模型
func setModels() {
	orm.RegisterModel(
		new(User),
		new(AgentLog),
	)
}

func setDB() {
	dbtype = beego.AppConfig.String("dbtype")
	dbuser = beego.AppConfig.String("dbuser")
	dbpass = beego.AppConfig.String("dbpass")
	dbhost = beego.AppConfig.String("dbhost")
	dbport = beego.AppConfig.String("dbport")
	dbname = beego.AppConfig.String("dbname")
	switch dbtype {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpass, dbhost, dbport, dbname) + "&loc=" + url.QueryEscape("Local")
		break
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@:%s/%s?sslmode=%s&host=%s", url.QueryEscape(dbuser), url.QueryEscape(dbpass), dbport, dbname, "disable", dbhost)
		break
	case "sqlite3":
		dsn = fmt.Sprintf("./conf/iagent.db")
	default:
		beego.Error("不支持此类型数据库")
	}

	maxIdle = 30
	maxConn = 50
}

func setBeego() {
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
}

// 设定默认配置
func predefineConfig() {
	var nowTime int64 = time.Now().Unix()
	n, err := Orm.QueryTable(new(User)).Filter("username", "admin").Count()
	if err != nil {
		fmt.Println("数据库连接错误" + err.Error())
		os.Exit(1)
	}
	if n == 0 {
		Orm.Raw(`insert into user (userid,username,possword,salt,register_time,update_time,login_time,login_ip,sessionid,enable) 
		values (?,?,?,?,?,?,?,?,?,?)`, 1, "admin", "1de5e16b6dd8d919c0f5800f14786340", "[eo7y#b$UXnC.#56", nowTime, nowTime, nowTime, "127.0.0.1", "", 1).Exec()
	}
}

func init() {
	setDB()
	setModels()
	//setBeego()

	switch dbtype {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
	case "postgres":
		orm.RegisterDriver("postgres", orm.DRPostgres)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
		//orm.RegisterDataBase("default", "sqlite3", "./conf/slaver.db")
	default:
		beego.Error("不支持", dbtype, "数据库")
		os.Exit(1)
	}
	orm.RegisterDataBase("default", dbtype, dsn, maxIdle, maxConn)
	orm.Debug = false

	//orm.RegisterDataBase("default", dbtype, "oneoaas_event:oneoaas_event@tcp(127.0.0.1:3306)/oneoaas_event?charset=utf8&loc=Local", maxIdle, maxConn)
	//force=true 会先删除表,后重建表
	//verbose=true 显示执行信息
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Error("RunSyncdb错误,msg:%v", err)
		os.Exit(1)
	}
	Orm = orm.NewOrm()
	Orm.Using("default") // 默认使用 default，你可以指定为其他数据库
	predefineConfig()
}

func dbErrorParse(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}

//3|root|1de5e16b6dd8d919c0f5800f14786340|[eo7y#b$UXnC.#56|1528205487|1528205487|1|0||
//adminABCD
