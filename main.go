package main

import (
	"github.com/astaxie/beego"
	_ "shiqiao/models"
	_ "shiqiao/routers"
)

func main() {
	beego.Run()
}
