package controllers

import (
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type ServiceController struct {
	Ctx iris.Context
	WsManager manager.WSManager
	CrManager manager.CRManager
}

func (c *ServiceController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/setWebChat","AllocateDoctorForTreat")
}

//客服分配医生给用户
func (c *ServiceController) AllocateDoctorForTreat() {
	var allocation models.TreatAllocation

	if err := c.Ctx.ReadJSON(&allocation); err != nil {
		fmt.Println("fail to encode request")
		return
	}

	//验证账号是否存在

	_, _ = c.Ctx.JSON(models.BaseResponse{
		Status:  "2000",
		Message: "Successful",
	})

	//分配成功后，从空闲的房间号中选择一个返回给对应医生和用户
	idleRoom := strconv.Itoa(c.CrManager.GetIdleRoom())

	DoctorConn := c.WsManager.GetWSConnection(allocation.Doctor)
	dataForDoctor, _ := json.Marshal(models.WebsocketResponse{
		Code:   "2011",
		RoomID: idleRoom ,
	})
	_ = (*DoctorConn).Write(1, dataForDoctor)

	PatientConn := c.WsManager.GetWSConnection(allocation.Patient)
	dataForPatient, _ := json.Marshal(models.WebsocketResponse{
		RobotResponse:models.RobotResponse{
			Accont:allocation.Patient,
			UniqueID:"",
			ClientType:"robot",
		},
		Status:"2000",
		Message:"",
		RoomID:idleRoom,
	})
	_ = (*PatientConn).Write(1,dataForPatient)

}