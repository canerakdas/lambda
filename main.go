package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"lambda/cache"
	"lambda/conf"
	"lambda/crawlers"
	_ "lambda/routers"
)

func main() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", static.OrmConnectionString)

	c := cron.New()
	c.AddFunc("@every 1h0m0s", cache.GameViewCollect)
	c.AddFunc("@every 12h0m0s", crawlers.UpdateSteamGameList)
	c.AddFunc("@every 0h10m0s", crawlers.SyncSteamGameList)
	//go crawlers.SyncSteamGameList()
	c.Start()

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	beego.SetStaticPath("/static", "static")

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
