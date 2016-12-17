package controllers

import (
	"github.com/astaxie/beego"
	"github.com/voidsand/vblog/models"
	"strconv"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	var tId string
	if c.Ctx.Request.Method == "POST" {
		tId = c.Input().Get("tid")
		nickname := c.Input().Get("nickname")
		content := c.Input().Get("content")
		err := models.AddReply(tId, nickname, content)
		if err != nil {
			beego.Error(err)
		}
		models.TopicViewsChange(tId, false)
	}
	c.Redirect("/topic/view/"+tId, 301)
}

func (c *ReplyController) Delete() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	rId := c.Ctx.Input.Param("0")
	topic, err := models.GetTopicByReply(rId)
	if err != nil {
		beego.Error(err)
	}
	err = models.DeleteReply(rId)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+strconv.FormatInt(topic.Id, 10), 301)
}
