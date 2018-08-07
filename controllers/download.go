package controllers

import (
	"iagent/util"
)

type DownLoadController struct {
	BaseController
}

func (this *DownLoadController) Get() {
	filename := this.Ctx.Input.Param(":filename")
	if filename == "" {
		this.Abort("404")
	}
	filepath := "rpm/" + filename
	exists, err := util.PathExists(filepath)
	if err != nil {
		this.Abort("500")
	}

	if exists {
		this.Ctx.Output.Download(filepath)
	} else {
		this.Abort("404")
	}
}
