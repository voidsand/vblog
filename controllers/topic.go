package controllers

import (
	"github.com/astaxie/beego"
	"vblog/models"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	var err error
	c.TplName = "topic.html"
	c.Data["Title"] = "文章 - 我的博客"
	c.Data["IsTopic"] = true
	c.Data["LoginReady"] = checkLogin(c.Ctx)
	c.Data["Topics"], err = models.GetAllTopics(false)
	if err != nil {
		beego.Error(err)
	}
}

func (c *TopicController) Post() {
	var err error
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	tid := c.Input().Get("tid")
	if len(tid) == 0 {
		err = models.AddTopic(title, content)
	} else {
		err = models.ModifyTopic(tid, title, content)
	}
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 301)
	return
}

func (c *TopicController) Add() {
	if !checkLogin(c.Ctx) {
		c.Redirect("/topic", 301)
		return
	}
	c.TplName = "topic_add.html"
	c.Data["Title"] = "添加文章 - 我的博客"
}

func (c *TopicController) Modify() {
	if !checkLogin(c.Ctx) {
		c.Redirect("/topic", 301)
		return
	}
	c.TplName = "topic_modify.html"
	c.Data["Title"] = "文章修改 - 我的博客"

	tid := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/topic", 301)
		return
	}
	c.Data["Topic"] = topic
}

func (c *TopicController) Delete() {
	if !checkLogin(c.Ctx) {
		c.Redirect("/topic", 301)
		return
	}
	tid := c.Ctx.Input.Param("0")
	err := models.DeleteTopic(tid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 301)
	return
}

func (c *TopicController) View() {
	c.TplName = "topic_view.html"
	c.Data["Title"] = "文章浏览 - 我的博客"
	c.Data["LoginReady"] = checkLogin(c.Ctx)

	tid := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 301)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}
