package main

import (
	"time"

	"WatchBitcoinAddress/controllers"

	"github.com/astaxie/beego"
)

func main() {
	mainController := controllers.GetMainController()
	controllers.StoreAllBitSite()
	go mainController.Timer(1 * time.Second)

	beego.Router("/", mainController, "get:Index")

	beego.Run()
}
