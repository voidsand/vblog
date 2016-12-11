package controllers

import (
	"github.com/astaxie/beego"
)

type AboutController struct {
	beego.Controller
}

func (c *AboutController) Get() {
	c.TplName = "about.html"
	c.Data["Title"] = "关于 - 我的博客"
	c.Data["IsAbout"] = true
	c.Data["LoginReady"] = checkSignin(c)
}
