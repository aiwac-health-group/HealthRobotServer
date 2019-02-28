package controllers

import (
	"HealthRobotServer/middleware"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
)

type WebsocketController struct {
	Ctx iris.Context
	Conn websocket.Connection
}

func (c *WebsocketController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.NewAuthToken().Serve)
	b.Handle("GET","/","Join")
}


func (c *WebsocketController) Join() {

}

