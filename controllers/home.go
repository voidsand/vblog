package controllers

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.TplName = "home.html"
	c.Data["Title"] = "首页 - 我的博客"
	c.Data["IsHome"] = true
	c.Data["LoginReady"] = checkLogin(c.Ctx)
}
