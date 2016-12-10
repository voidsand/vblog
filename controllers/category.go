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
	c.Data["LoginReady"] = checkLogin(c.Ctx)
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}

	cname := c.Input().Get("cname")
	if len(cname) != 0 && checkLogin(c.Ctx) {
		err := models.AddCategory(cname)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 301)
		return
	}
}

func (c *CategoryController) Delete() {
	if !checkLogin(c.Ctx) {
		c.Redirect("/category", 301)
		return
	}
	cid := c.Ctx.Input.Param("0")
	err := models.DeleteCategory(cid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/category", 301)
	return
}
