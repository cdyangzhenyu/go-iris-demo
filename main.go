package main

import (
	"fpga-bms-server/router"
	"fpga-bms-server/config"
	"github.com/kataras/iris"
)

func main() {
	app := router.Router

	addr := config.Conf.Get("app.addr").(string)
	app.Run(iris.Addr(addr))
}
