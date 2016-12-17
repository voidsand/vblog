package controllers

import (
	"github.com/astaxie/beego"
	"github.com/voidsand/vblog/models"
	"strconv"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.TplName = "category.html"
	c.Data["Title"] = "分类 - 我的博客"
	c.Data["IsCategory"] = true
	c.Data["SigninReady"] = checkSignin(c)
	cates, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	for i := range cates {
		models.TotalViewsChange(strconv.FormatInt(cates[i].Id, 10))
	}
	c.Data["Categories"] = cates
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
