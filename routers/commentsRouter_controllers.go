package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"] = append(beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"],
		beego.ControllerComments{
			Method: "InstallAgent",
			Router: `/install`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"] = append(beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"],
		beego.ControllerComments{
			Method: "GetAgentLog",
			Router: `/log`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"] = append(beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"],
		beego.ControllerComments{
			Method: "GetAgentLogByIp",
			Router: `/logbyip`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"] = append(beego.GlobalControllerRouter["iagent/controllers:InstallAgentController"],
		beego.ControllerComments{
			Method: "RemoveAgent",
			Router: `/remove`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddUser",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "UpdateUser",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "DelUser",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAllUsers",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["iagent/controllers:UserController"] = append(beego.GlobalControllerRouter["iagent/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetUserById",
			Router: `/uid/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
