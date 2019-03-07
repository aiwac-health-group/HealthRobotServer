package controllers

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
	"strconv"
	"strings"
	"time"
)

//service controller处理客服发出的http请求
type ServiceController struct {
	Ctx iris.Context
	Service services.ServiceService
	WsManager manager.WSManager
	CrManager manager.CRManager
}

func (c *ServiceController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandler().Serve, middleware.NewAuthToken().Serve)
	b.Handle("POST","/changeDoctor","ModifyDoctorProfile")
	b.Handle("POST","/setWebChat","AllocateDoctorForTreat")
}

//客服修改医生信息
func (c *ServiceController) ModifyDoctorProfile() {
	var request models.DoctorProfileModifyRequest
	if err := c.Ctx.ReadJSON(&request); err != nil {
		log.Println("fail to encode request")
		return
	}

	//判断医生是否已经存在
	if client := c.Service.SearchDoctorClientInfo(request.Account); client == nil {
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
		ClientType:"doctor",
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

//客服分配医生给用户
func (c *ServiceController) AllocateDoctorForTreat() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	service := claims["Account"].(string)

	var allocation models.TreatAllocation

	if err := c.Ctx.ReadJSON(&allocation); err != nil {
		fmt.Println("fail to encode request")
		return
	}

	//验证分配的医生账号是否存在
	doctor := c.Service.SearchDoctorClientInfo(allocation.Doctor)
	if doctor == nil || !strings.EqualFold(doctor.ClientType, constants.ClientType_doctor) { //账号不存在，或者分配的账号不是医生账号
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:  "2001",
			Message: "Wrong Account",
		})
	}

	//获取医生和客服的websocket连接，保证两边连接有效时，才可以发送roomID
	DoctorConn := c.WsManager.GetWSConnection(allocation.Doctor)
	PatientConn := c.WsManager.GetWSConnection(allocation.Patient)
	if DoctorConn == nil || PatientConn == nil {
		_, _ = c.Ctx.JSON(models.BaseResponse{
			Status:  "2001",
			Message: "医生或用户已经下线，请更新问诊列表",
		})
		return
	}

	//更新该treat中责任医生的账号,以及问诊单状态为正在处理
	treat := c.Service.SearchNewTreatInfo(allocation.Patient)
	treat.HandleDoctor = allocation.Doctor
	treat.Status = constants.Status_treat_onHandle
	c.Service.UpdateTreatInfoHandleDoctor(treat)

	//分配成功后，根据当前时间的纳秒值返回一个roomID给对应医生和用户
	RoomID := strconv.FormatInt(time.Now().UnixNano(),10)

	//将语音连接请求发送给医生
	dataForDoctor, _ := json.Marshal(models.WebsocketResponse{
		Code:   "2011",
		RoomID: RoomID ,
	})
	if err := (*DoctorConn).Write(1, dataForDoctor); err != nil {
		log.Println("Fail to Send Call to Doctor")
		return
	}
	//更新医生状态
	doctor.OnlineStatus = constants.Status_onbusy
	if err := c.Service.UpdateDoctorClientInfo(doctor); err != nil {
		log.Println("Fail to Update Doctor status")
		return
	}
	//把新的在线医生列表推送给客服
	CurrentConn := c.WsManager.GetWSConnection(service)
	c.SendOnlineDoctorList(*CurrentConn)

	//把语音通话房间号返回给机器人
	dataForPatient, _ := json.Marshal(models.WebsocketResponse{
		RobotResponse:models.RobotResponse{
			Account:allocation.Patient,
			UniqueID:"",
			ClientType:"robot",
		},
		Status:"2000",
		Message:"",
		RoomID:RoomID,
	})
	_ = (*PatientConn).Write(1,dataForPatient)

}

//获取在线医生列表
func (c *ServiceController) SendOnlineDoctorList(conn websocket.Connection) {
	doctors := c.Service.GetOnlineDoctor()
	var items []interface{}
	for _, value := range doctors {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2009",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("OnlineDoctorList: %s", data)
	_ = conn.To("service").EmitMessage(data)
}

