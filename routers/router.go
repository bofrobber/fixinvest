package routers

import (
	"github.com/astaxie/beego"
	"github.com/bofrobber/fixinvest/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

}
