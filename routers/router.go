package routers

import (
	"github.com/bofrobber/fixinvest/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
