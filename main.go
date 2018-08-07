package main

import (
	_ "iagent/routers"

	"iagent/task"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.SetStaticPath("/dist", "static/dist")
	go task.Docron()
	beego.Run()
}
