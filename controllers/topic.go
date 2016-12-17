package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/voidsand/vblog/models"
	"os"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	var err error
	c.TplName = "topic.html"
	c.Data["Title"] = "文章 - 我的博客"
	c.Data["IsTopic"] = true
	c.Data["SigninReady"] = checkSignin(c)
	c.Data["Topics"], err = models.GetAllTopics("", false)
	if err != nil {
		beego.Error(err)
	}
}

func (c *TopicController) Add() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	switch c.Ctx.Request.Method {
	case "GET":
		var err error
		c.TplName = "topic_add.html"
		c.Data["Title"] = "添加文章 - 我的博客"
		c.Data["Categories"], err = models.GetAllCategories()
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
		}
	case "POST":
		tTitle := c.Input().Get("title")
		cId := c.Input().Get("category")
		tContent := c.Input().Get("content")

		var tAttachment string
		_, fh, err := c.GetFile("attachment")
		if err != nil {
			beego.Error(err)
		}
		if fh != nil {
			tAttachment = fh.Filename
		}
		tId, err := models.AddTopic(tTitle, cId, tContent, tAttachment)
		if err != nil {
			beego.Error(err)
		}
		if fh != nil {
			os.Mkdir(path.Join("attachment", tId), os.ModePerm)
			err = c.SaveToFile("attachment", path.Join("attachment", tId, tAttachment))
			if err != nil {
				beego.Error(err)
			}
		}
		c.Redirect("/topic", 301)
	}
}

func (c *TopicController) Modify() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	switch c.Ctx.Request.Method {
	case "GET":
		var err error
		var topic *models.Topic
		var category *models.Category

		c.TplName = "topic_modify.html"
		c.Data["Title"] = "文章修改 - 我的博客"
		c.Data["Categories"], err = models.GetAllCategories()
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
			return
		}
		tId := c.Ctx.Input.Param("0")
		category, err = models.GetCategoryByTopic(tId)
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
			return
		}
		c.Data["Cid"] = category.Id

		topic, err = models.GetTopic(tId)
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
			return
		}
		c.Data["Topic"] = topic
	case "POST":
		tId := c.Input().Get("tid")
		tTitle := c.Input().Get("title")
		cId := c.Input().Get("category")
		TContent := c.Input().Get("content")

		var tAttachment string
		_, fh, err := c.GetFile("attachment")
		fmt.Println(fh)
		if err != nil {
			beego.Error(err)
		}
		if fh != nil {
			tAttachment = fh.Filename
		}
		err = models.ModifyTopic(tId, tTitle, cId, TContent, tAttachment)
		if err != nil {
			beego.Error(err)
		}
		if fh != nil {
			os.Mkdir(path.Join("attachment", tId), os.ModePerm)
			err := c.SaveToFile("attachment", path.Join("attachment", tId, tAttachment))
			if err != nil {
				beego.Error(err)
			}
		}
		c.Redirect("/topic", 301)
	}
}

func (c *TopicController) Delete() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	tId := c.Ctx.Input.Param("0")
	err := models.DeleteTopic(tId)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 301)
}

func (c *TopicController) View() {
	c.TplName = "topic_view.html"
	c.Data["Title"] = "文章浏览 - 我的博客"
	c.Data["SigninReady"] = checkSignin(c)

	tId := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tId)
	if err != nil {
		beego.Error(err)
		c.Redirect("/topic", 301)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tId

	replies, err := models.GetAllReplies(tId)
	if err != nil {
		beego.Error(err)
		c.Redirect("/topic", 301)
		return
	}
	c.Data["Replies"] = replies
	models.TopicViewsChange(tId, true)
}
