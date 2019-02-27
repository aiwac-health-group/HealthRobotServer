package main

import (
	"HealthRobotServer/controllers"
	"HealthRobotServer/datasource"
	"HealthRobotServer/services"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func newApp() (api *iris.Application) {
	api = iris.New()

	iris.RegisterOnInterrupt(func() {
		datasource.Instance().Close()
	})

	mvc.Configure(api.Party("/login"), func(app *mvc.Application) {
		app.Register(services.NewLoginService())
		app.Handle(new(controllers.LoginController))
	})

	mvc.Configure(api.Party("/admin"), func(app *mvc.Application) {
		app.Handle(new(controllers.AdminController))
	})

	return api
}

func main() {
	app := newApp()
	_ = app.Run(iris.Addr(":8080"))
}