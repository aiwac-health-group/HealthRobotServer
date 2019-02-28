package controllers

import (
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"strings"
)

type AdminController struct {
	Ctx iris.Context
	Service services.ClientService
}

func (c *AdminController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/addService","AddService")
	b.Handle("POST","/addDoctor","AddDoctor")
	b.Handle("POST","/changeAccount","ModifyClientAccount")
}

//管理员添加客服账号
func (c *AdminController) AddService() {

	var request models.AccountAddRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}
	err := c.Service.CreatClient(&models.ClientInfo{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientPassword:request.Password,
		ClientType:"service",
	})

	if err != nil {
		_, _ =c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"failed to add service",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"successfully",
	})

	return
}

//管理员添加医生账号
func (c *AdminController) AddDoctor() {

	var request models.AccountAddRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	err := c.Service.CreatClient(&models.ClientInfo{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientPassword:request.Password,
		ClientType:"doctor",
	})

	if err != nil {
		_, _ =c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"failed to add doctor",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"successfully",
	})

	return
}

//管理员修改其他账号信息
func (c *AdminController) ModifyClientAccount() {

	var request models.AccountModifyRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//判断账号是否存在
	if client := c.Service.GetClient(request.Account); client == nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account doest not exist",
		})
		return
	} else {
		if strings.EqualFold(request.OperationType, "changPass") { //修改密码
			c.Service.UpdateClient(&models.ClientInfo{
				ClientAccount:request.Account,
				ClientPassword:request.Value,
			})
		} else { //修改账户姓名
			c.Service.UpdateClient(&models.ClientInfo{
				ClientAccount:request.Account,
				ClientName:request.Value,
			})
		}
	}

}