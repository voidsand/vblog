package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/voidsand/vblog/models"
	_ "github.com/voidsand/vblog/routers"
)

func init() {
	models.RegisterDB()
	beego.AddFuncMap("plus1", plus1)
}

func main() {
	orm.RunSyncdb("default", false, true)
	beego.Run()
}

func plus1(in int) int {
	return in + 1
}
