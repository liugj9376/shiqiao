package routers

import (
	"github.com/astaxie/beego"
	"shiqiao/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
