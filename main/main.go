package main

import (
	"HealthRobotServer/controllers"
	"HealthRobotServer/datasource"
	"HealthRobotServer/manager"
	"HealthRobotServer/services"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
)

func newApp() (api *iris.Application) {
	api = iris.New()

	iris.RegisterOnInterrupt(func() {
		log.Println("the server is closing")
		datasource.Instance().Close()
	})

	//注册模板
	api.RegisterView(iris.HTML("./view",".html"))
	//静态文件支持
	api.StaticWeb("/static", "./static")
	api.StaticWeb("/uploads", "./uploads")

	mvc.Configure(api.Party("/login"), func(app *mvc.Application) {
		app.Register(services.NewLoginService())
		app.Handle(new(controllers.LoginController))
	})

	mvc.Configure(api.Party("/admin"), func(app *mvc.Application) {
		app.Register(services.NewStatisticService())
		app.Register(services.NewClientService())
		app.Handle(new(controllers.AdminController))
	})

	mvc.Configure(api.Party("/service"), func(app *mvc.Application) {
		app.Register(manager.WSInstance())
		app.Register(services.NewServiceService())
		app.Handle(new(controllers.ServiceController))
	})

	mvc.Configure(api.Party("/doctor"), func(app *mvc.Application) {
		app.Register(services.NewDoctorService())
		app.Register(manager.WSInstance())
		app.Handle(new(controllers.DoctorController))
	})

	mvc.Configure(api.Party("/ws"), func(app *mvc.Application) {
		ws := websocket.New(websocket.Config{})
		app.Register(services.NewWebsocketService())
		app.Register(manager.WSInstance())
		app.Register(ws.Upgrade)
		app.Handle(new(controllers.WebsocketController))
	})

	return api
}

func main() {
	app := newApp()

	go app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
	//开启新服务器，websocket监听8082端口
	app.Run(iris.Addr(":8082"), iris.WithoutInterruptHandler)

}