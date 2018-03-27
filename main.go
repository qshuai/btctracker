package main

import (
	"time"

	"github.com/astaxie/beego"
	"WatchBitcoinAddress/controllers"
)

func main() {
	mainController := controllers.GetMainController()
	controllers.StoreAllBitSite()
	go mainController.Timer(1 * time.Second)

	beego.Router("/", mainController, "get:Index")

	beego.Run()
}
