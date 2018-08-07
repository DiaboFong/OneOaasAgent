// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"iagent/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/*", &controllers.IndexController{})
	beego.Router("/download/:filename", &controllers.DownLoadController{}, "get:Get")
	beego.Router("/api/v1/dashboard/count", &controllers.InstallAgentController{}, "get:GetCount")
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/agent",
			beego.NSInclude(
				&controllers.InstallAgentController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
