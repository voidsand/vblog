package controllers

import (
	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.TplName = "topic.html"
	c.Data["Title"] = "文章 - 我的博客"
	c.Data["IsTopic"] = true
}

func (c *TopicController) Add() {
	c.Ctx.WriteString("Add")
}
