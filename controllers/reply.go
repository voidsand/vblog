package controllers

import (
	"github.com/astaxie/beego"
	"vblog/models"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	var tid string
	if c.Ctx.Request.Method == "POST" {
		tid = c.Input().Get("tid")
		nickname := c.Input().Get("nickname")
		content := c.Input().Get("content")
		err := models.AddReply(tid, nickname, content)
		if err != nil {
			beego.Error(err)
		}
	}
	c.Redirect("/topic/view/"+tid, 301)
}

func (c *ReplyController) Delete() {
	if !checkSignin(c) {
		c.Redirect("/topic", 301)
		return
	}
	tid := c.Ctx.Input.Param("0")
	rid := c.Ctx.Input.Param("1")
	err := models.DeleteReply(tid, rid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 301)
}
