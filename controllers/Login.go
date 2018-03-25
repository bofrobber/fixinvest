package controllers

import (
	"github.com/astaxie/beego"
)

//login control，这个提供了基本的登录的

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
