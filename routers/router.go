package routers

import (
	"github.com/astaxie/beego"
	"os"
	"github.com/voidsand/vblog/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/sign", &controllers.SignController{})
	beego.AutoRouter(&controllers.SignController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.AutoRouter(&controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.AutoRouter(&controllers.ReplyController{})
	beego.Router("/about", &controllers.AboutController{})

	os.Mkdir("attachment", os.ModePerm)
	beego.SetStaticPath("/attachment", "attachment")
}
