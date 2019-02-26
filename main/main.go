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
		//绑定数据库服务
		app.Register(services.NewClientService())
		app.Handle(new(controllers.LoginController))
	})

	return api
}

func main() {
	app := newApp()
	_ = app.Run(iris.Addr(":8080"))
}