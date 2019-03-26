// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"lambda/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/game",
			beego.NSInclude(
				&controllers.GameController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/list",
			beego.NSInclude(
				&controllers.ListController{},
			),
		),
		beego.NSRouter(
			"/token",
			&controllers.TokenController{},
		),
		beego.NSRouter(
			"/confirmation",
			&controllers.ConfirmationController{},
		),
	)

	beego.AddNamespace(ns)
}
