package controllers

import (
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
)

type AdminController struct {
	Ctx iris.Context
	
}

func (c *AdminController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/addService","AddService")
	b.Handle("POST","/addDoctor","AddDoctor")
}

//管理员添加客服账号
func (c *AdminController) AddService() {

	var request models.ServiceAddRequest

	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}



}

//管理员添加医生账号
