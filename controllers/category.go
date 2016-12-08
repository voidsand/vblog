package controllers

import (
	"github.com/astaxie/beego"
	"vblog/models"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	var err error
	c.TplName = "category.html"
	c.Data["Title"] = "分类 - 我的博客"
	c.Data["IsCategory"] = true
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	op := c.Input().Get("op")
	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 301)
		return
	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err := models.DelCategory(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 301)
		return
	}
}
