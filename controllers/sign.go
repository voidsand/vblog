package controllers

import (
	"github.com/astaxie/beego"
	"reflect"
)

type SignController struct {
	beego.Controller
}

func (c *SignController) Get() {
	if checkSignin(c) {
		c.Redirect("/", 301)
		return
	}
	c.TplName = "sign.html"
	c.Data["Title"] = "登录 - 我的博客"
	c.Data["IsSignin"] = true
}

func (c *SignController) In() {
	if checkSignin(c) {
		c.Redirect("/", 301)
		return
	}
	if c.Ctx.Request.Method == "POST" {
		uname := c.Input().Get("uname")
		pwd := c.Input().Get("pwd")
		autoSignin := c.Input().Get("autoSignin") == "on"

		if beego.AppConfig.String("uname") == uname && beego.AppConfig.String("pwd") == pwd {
			if autoSignin {
				// Session自动登录如何实现?
			}
			c.SetSession("uname", string(uname))
			c.SetSession("pwd", string(pwd))
		}
	}
	c.Redirect("/", 301)
}

func (c *SignController) Out() {
	if !checkSignin(c) {
		c.Redirect("/", 301)
		return
	}
	c.DelSession("uname")
	c.DelSession("pwd")
	c.Redirect("/", 301)
}

func checkSignin(c interface{}) bool {
	m := reflect.ValueOf(c).MethodByName("GetSession")
	r := m.Call([]reflect.Value{reflect.ValueOf("uname")})
	uname := r[0].Interface()
	r = m.Call([]reflect.Value{reflect.ValueOf("pwd")})
	pwd := r[0].Interface()
	if uname == nil || pwd == nil {
		return false
	}
	return beego.AppConfig.String("uname") == uname.(string) && beego.AppConfig.String("pwd") == pwd.(string)
}
