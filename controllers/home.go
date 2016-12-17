package controllers

import (
	"github.com/astaxie/beego"
	"github.com/voidsand/vblog/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	var err error
	c.TplName = "home.html"
	c.Data["Title"] = "首页 - 我的博客"
	c.Data["IsHome"] = true
	c.Data["SigninReady"] = checkSignin(c)
	c.Data["Topics"], err = models.GetAllTopics(c.Input().Get("cid"), true)
	if err != nil {
		beego.Error(err)
	}
	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Categories"] = categories
}
