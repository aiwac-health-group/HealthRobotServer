package controllers

import (
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math"
)

type DoctorController struct {
	Ctx iris.Context
	Service services.ClientService
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/userList","GetUserList")
}

//获取待填写健康报告的用户列表
//获取用户总数和医生总数，相除，向上取整得每个医生负责的用户数N，从而获得每个医生负责的用户ID的范围即：医生ID * N ~ （医生ID + 1）* N
//注意ID是在数据库表中的序号，不是Account
//根据ID范围以及用户的HealthStatus状态来返回列表给医生
func (c *DoctorController) GetUserList() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	account := claims["Account"].(string)

	//获取该医生在数据库表中的ID
	doctor := c.Service.SearchDoctorClientInfo(account)
	doctorID := doctor.ID

	//获取注册用户和医生的数目
	RobotCount := c.Service.CountTotalRobotClient()
	DoctorCount := c.Service.CountTotalDoctorClient()

	//计算医生负载
	payLoad := int64(math.Ceil(float64(RobotCount)/float64(DoctorCount)))
	//计算该医生的负责用户ID范围
	DownBoundary := doctorID * payLoad
	UpBoundary := (doctorID + 1) * payLoad

	//根据ID范围获取需要处理的用户列表,并发送给医生
	robots := c.Service.GetRobotListForHealthReport(UpBoundary, DownBoundary)
	var items []interface{}
	for _, value := range robots {
		items = append(items, value)
	}
	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:"2000",
		Data:models.List{
			Items:items,
		},
	})
}