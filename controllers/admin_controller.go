package controllers

import (
	"HealthRobotServer/constants"
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
	b.Handle("POST","/changeDoctor","ModifyDoctorProfile")
}

//管理员添加客服账号
func (c *AdminController) AddService() {

	var request models.AccountAddRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//查询该客服的账户信息是否已经存在
	if service := c.Service.SearchServiceClientInfo(request.Account); service.ID != 0 {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account have already exist",
		})
		return
	}

	//添加客服信息
	var info = models.ServiceInfo{
		ClientAccount:request.Account,
		ClientPassword:request.Password,
		ClientName:request.Name,
		ClientType:constants.ClientType_service,
	}
	
	if err := c.Service.UpdateServiceClientInfo(&info); err != nil {
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

	//查询账户信息是否已经存在
	if client := c.Service.SearchDoctorClientInfo(request.Account); client.ID != 0 {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account have already exist",
		})
		return
	}

	//添加医生信息
	var info = models.DoctorInfo{
		ClientAccount:request.Account,
		ClientPassword:request.Password,
		ClientName:request.Name,
		ClientType:constants.ClientType_doctor,
	}

	if err := c.Service.UpdateDoctorClientInfo(&info); err != nil {
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

}

//管理员修改账号信息
func (c *AdminController) ModifyClientAccount() {

	var request models.AccountModifyRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//判断账号是否存在
	service := c.Service.SearchServiceClientInfo(request.Account)
	doctor := c.Service.SearchDoctorClientInfo(request.Account)
	if service.ID == 0 && doctor.ID == 0 { //账户不存在
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account doest not exist",
		})
		return
	}

	//账户存在
	if service.ID != 0 { //为客服账号
		if strings.EqualFold(request.OperationType, "changPass") { //修改密码
			var info = models.ServiceInfo{
				ClientAccount:request.Account,
				ClientPassword:request.Value,
			}
			if err := c.Service.UpdateServiceClientInfo(&info); err != nil {
				_, _ = c.Ctx.JSON(models.BaseResponse{
					Status:"2001",
					Message:"Fail to update service's password",
				})
				return
			}
		} else { //修改账户姓名,则更新详细信息表
			var info = models.ServiceInfo{
				ClientAccount:request.Account,
				ClientName:request.Value,
			}
			if err := c.Service.UpdateServiceClientInfo(&info); err != nil {
				_, _ = c.Ctx.JSON(models.BaseResponse{
					Status:"2001",
					Message:"Fail to update service's name",
				})
				return
			}
		}
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2000",
			Message:"Successful",
		})
		return
	} else { //为医生账户
		if strings.EqualFold(request.OperationType, "changPass") { //修改密码
			var info = models.DoctorInfo{
				ClientAccount:request.Account,
				ClientPassword:request.Value,
			}
			if err := c.Service.UpdateDoctorClientInfo(&info); err != nil {
				_, _ = c.Ctx.JSON(models.BaseResponse{
					Status:"2001",
					Message:"Fail to update doctor's password",
				})
				return
			}
		} else { //修改账户姓名,则更新详细信息表
			var profile = models.DoctorInfo{
				ClientAccount:request.Account,
				ClientName:request.Value,
			}
			if err := c.Service.UpdateDoctorClientInfo(&profile); err != nil {
				_, _ = c.Ctx.JSON(models.BaseResponse{
					Status:"2001",
					Message:"Fail to update doctor's name",
				})
				return
			}
		}
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2000",
			Message:"Successful",
		})
		return
	}


}

//管理员修改医生信息
func (c *AdminController) ModifyDoctorProfile() {
	var request models.DoctorProfileModifyRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//判断医生是否已经存在
	if doctor := c.Service.SearchDoctorClientInfo(request.Account); doctor.ID == 0 {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"Account doesn't exist",
		})
		return
	}

	//添加详细信息
	var profile = models.DoctorInfo{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientType:constants.ClientType_doctor,
		Department:request.Department,
		Brief:request.Brief,
	}

	if err := c.Service.UpdateDoctorClientInfo(&profile); err != nil {
		_, _ =c.Ctx.JSON(models.BaseResponse{
			Status:"2001",
			Message:"failed to modify doctor profile",
		})
		return
	}

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Message:"successfully",
	})

}