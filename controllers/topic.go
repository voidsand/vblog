package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
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
	c.Data["LoginReady"] = checkSignin(c)
	c.Data["Topics"], err = models.GetAllTopics(false)
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
		var err error
		var category *models.Category
		title := c.Input().Get("title")
		cid := c.Input().Get("category")
		content := c.Input().Get("content")

		category, err = models.GetCategory(cid)
		if err != nil {
			beego.Error(err)
		}
		err = models.AddTopic(title, category.Title, content)
		if err != nil {
			beego.Error(err)
		}
		err = models.TopicCountUp(cid, true)
		if err != nil {
			beego.Error(err)
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
		tid := c.Ctx.Input.Param("0")
		category, err = models.GetCategoryByTopicId(tid)
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
			return
		}
		c.Data["Cid"] = category.Id

		topic, err = models.GetTopic(tid)
		if err != nil {
			beego.Error(err)
			c.Redirect("/topic", 301)
			return
		}
		c.Data["Topic"] = topic
	case "POST":
		var err error
		var category *models.Category
		var mcategory *models.Category
		tid := c.Input().Get("tid")
		title := c.Input().Get("title")
		cid := c.Input().Get("category")
		content := c.Input().Get("content")

		category, err = models.GetCategory(cid)
		if err != nil {
			beego.Error(err)
		}
		mcategory, err = models.GetCategoryByTopicId(tid)
		if err != nil {
			beego.Error(err)
		}
		err = models.TopicCountUp(strconv.FormatInt(mcategory.Id, 10), false)
		if err != nil {
			beego.Error(err)
		}
		err = models.ModifyTopic(tid, title, category.Title, content)
		if err != nil {
			beego.Error(err)
		}
		err = models.TopicCountUp(cid, true)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/topic", 301)
	}
}

func (c *TopicController) Delete() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	var category *models.Category
	tid := c.Ctx.Input.Param("0")
	category, err := models.GetCategoryByTopicId(tid)
	if err != nil {
		beego.Error(err)
	}
	err = models.TopicCountUp(strconv.FormatInt(category.Id, 10), false)
	if err != nil {
		beego.Error(err)
	}
	err = models.DeleteTopic(tid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 301)
}

func (c *TopicController) View() {
	c.TplName = "topic_view.html"
	c.Data["Title"] = "文章浏览 - 我的博客"
	c.Data["LoginReady"] = checkSignin(c)

	tid := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 301)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
	var replies []*models.Comment
	replies, err = models.GetAllReplies(tid)
	fmt.Println(replies)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 301)
		return
	}
	c.Data["Replies"] = replies
}
