package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"vblog/models"
	_ "vblog/routers"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.RunSyncdb("default", false, true)
	beego.Run()
}
