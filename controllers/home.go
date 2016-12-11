package controllers

import (
	"github.com/astaxie/beego"
	"vblog/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	var err error
	c.TplName = "home.html"
	c.Data["Title"] = "首页 - 我的博客"
	c.Data["IsHome"] = true
	c.Data["LoginReady"] = checkSignin(c)
	c.Data["Topics"], err = models.GetAllTopics(true)
	if err != nil {
		beego.Error(err)
	}
}
