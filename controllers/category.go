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
	c.Data["LoginReady"] = checkSignin(c)
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
}

func (c *CategoryController) Add() {
	if !checkSignin(c) {
		c.Redirect("/category", 301)
		return
	}
	if c.Ctx.Request.Method == "POST" {
		cname := c.Input().Get("cname")
		if len(cname) != 0 && checkSignin(c) {
			err := models.AddCategory(cname)
			if err != nil {
				beego.Error(err)
			}
		}
	}
	c.Redirect("/category", 301)
}

func (c *CategoryController) Delete() {
	if !checkSignin(c) {
		c.Redirect("/category", 301)
		return
	}
	cid := c.Ctx.Input.Param("0")
	category, err := models.GetCategory(cid)
	if category.TopicCount == 0 {
		err = models.DeleteCategory(cid)
		if err != nil {
			beego.Error(err)
		}
	}
	c.Redirect("/category", 301)
}
