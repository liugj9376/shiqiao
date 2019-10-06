package routers

import (
	"github.com/astaxie/beego"
	"shiqiao/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")

	//active
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	//login
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")

}
