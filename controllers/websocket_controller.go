package controllers

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/manager"
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
	"strconv"
	"strings"
)

const (
	BusinessRobotProfile = 6
	BusinessTreatRequest = 17 //机器人发起问诊请求
	BusinessTreatHangOut = 24 //机器人挂断问诊请求
	BusinessOnlineDoctor = 2009 //客服获取在线医生列表
	BusinessTreatList = 2010 //客服获取待问诊列表
	BusinessDoctorHangOut = 2012 //医生主动挂断问诊电话
	BusinessDoctorRejectCall = 2013 //医生拒绝接听问诊电话
)

type WebsocketController struct {
	Ctx iris.Context
	Conn websocket.Connection
	Service services.WebsocketService
	WsManager manager.WSManager
}

func (c *WebsocketController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandlerWS().Serve, middleware.NewAuthToken().Serve)
	b.Handle("GET","/","Join")
}

//从token中提取的该用户的账号和账户类别信息
var (
	ws_account string
	ws_clientType string
	ws_clientName string
)

func (c *WebsocketController) Join() {
	//从token和数据库中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	ws_account = claims["Account"].(string)
	ws_clientType = claims["ClientType"].(string)
	profile := c.Service.SearchRobotProfile(ws_account)
	ws_clientName = profile.ClientName
	log.Println("New Websocket Connection: ",ws_account, ws_clientType, ws_clientName)

	//加入对应clientType的room, 每个room存放了相应用户类型的所有websocket连接
	c.Conn.Join(ws_clientType)
	//注册连接断开回调函数
	c.Conn.OnDisconnect(c.LoseConnection)
	//注册消息接收处理函数
	c.Conn.OnMessage(c.ReceiveRequest)
	//存储该用户和对应连接的映射关系
	c.WsManager.AddMapRelationship(ws_account,&(c.Conn))

	//更新用户状态为在线
	_ = c.Service.UpdateClientInfo(&models.ClientInfo{
		ClientAccount:ws_account,
		OnlineStatus:constants.Status_online,
	})

	//如果上线用户为医生
	//获取在线医生列表,并将列表推送给所有的客服
	if ws_clientType == constants.ClientType_doctor {
		c.PushOnlineDoctorList()
	}

	//开启事件监听
	c.Conn.Wait()
}

func (c *WebsocketController) LoseConnection() {
	//删除用户与连接的映射关系
	c.WsManager.DeleteMapRelationship(ws_account)

	//更新用户状态
	_ = c.Service.UpdateClientInfo(&models.ClientInfo{
		ClientAccount:ws_account,
		OnlineStatus:"1",
	})

	//如果离开的用户为doctor，则更新客服的在线医生列表
	if ws_clientType == "doctor" {
		c.PushOnlineDoctorList()
	}
	log.Printf("%s %s lose the connection", ws_account, ws_clientType)
}

func (c *WebsocketController) ReceiveRequest(data []byte) {
	//在这里解析收到的请求，根据请求中的业务号跳转到指定业务处理函数中进行处理
	var request models.WSRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Println("Websocket request from Explore Unmarshal err: ",err)
	}
	log.Printf("websocket request, %s", request)
	//根据request中的code字段配置相应的函数处理
	businessCode, _ := strconv.Atoi(request.BusinessCode)

	switch businessCode {
	case BusinessRobotProfile: c.RobotProfileHandler(&request)
	case BusinessTreatRequest: c.TreatRequestHandler(&request)
	case BusinessTreatHangOut: c.TreatHangOutHandler(&request)
	case BusinessOnlineDoctor: c.DoctorListRequestHandler(&request)
	case BusinessTreatList: c.TreatWaitListHandler(&request)
	case BusinessDoctorHangOut: c.DoctorHangOutHandler(&request)
	case BusinessDoctorRejectCall: c.DoctorRejectCallHandler(&request)

	}
}

//6号业务处理
//处理机器人用户发起的个人信息注册及修改
func (c *WebsocketController) RobotProfileHandler(request *models.WSRequest) {
	//判断账号是否已经存在
	if client := c.Service.SearchClientInfo(request.Account); client.ID == 0 {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0006",
			Status:"2001",
			Message:"the robot doesn't exist",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	//添加详细信息
	var profile = models.Robot{
		ClientAccount:request.Account,
		ClientName:request.Name,
		ClientType:"robot",
		Sex:request.Sex,
		Birthday:request.Birthday,
		Address:request.Address,
		Wechat:request.Wechat,
	}

	if err := c.Service.UpdateRobotProfile(&profile); err != nil {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code:"0006",
			Status:"2001",
			Message:"system error",
			RobotResponse:models.RobotResponse{
				Account:ws_account,
				UniqueID:"",
			},
		})
		_ = c.Conn.Write(1,data)
		return
	}

	data, _ := json.Marshal(models.WebsocketResponse{
		Code:"0006",
		Status:"2000",
		Message:"Register or update profile successfully",
		RobotResponse:models.RobotResponse{
			Account:ws_account,
			UniqueID:"",
		},
	})
	_ = c.Conn.Write(1,data)
	return
}

//17号业务处理
//处理用户发起的问诊请求
//把问诊请求存进数据库，并把等待列表推送给在线客服
func (c *WebsocketController) TreatRequestHandler(request *models.WSRequest) {
	//把问诊请求存进数据库
	c.Service.CreatTreatInfoRequest(&models.TreatInfo{
		Account:ws_account,
		ClientName:ws_clientName,
		Others:"",
	})

	//把正在等待的问诊请求列表推送给客服
	c.PushTreatWaitList()
}

//24号业务处理
//机器人端挂断问诊请求
func (c *WebsocketController) TreatHangOutHandler(request *models.WSRequest) {
	//获取该用户问诊请求信息,更新其状态
	treat := c.Service.SearchNotCompleteTreatInfo("patient", ws_account)
	treat.Status = constants.Status_treat_complete
	c.Service.UpdateTreatInfoStatus(treat)
	//推送新的问诊列表给客服
	c.PushTreatWaitList()

	//同时更新对应医生的状态为空闲
	doctor := c.Service.SearchClientInfo(treat.HandleDoctor)
	if doctor.ID != 0 {
		doctor.OnlineStatus = constants.Status_online
		_ = c.Service.UpdateClientInfo(doctor)
	}
	//把空闲状态的医生列表推送给客服
	c.PushOnlineDoctorList()
}


//2009号业务处理
//推送列表至客服
func (c *WebsocketController) PushOnlineDoctorList() {
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
	_ = c.Conn.To("service").EmitMessage(data)
}

//获取在线医生列表
func (c *WebsocketController) GetOnlineDoctorList() {
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
	_ = c.Conn.Write(1,data)
}

func (c *WebsocketController) DoctorListRequestHandler(request *models.WSRequest) {
	if strings.EqualFold(request.Message, " getDoctorList ") {
		println("getDoctorList")
		c.GetOnlineDoctorList()
	}
}

//2010号业务处理
//推送等候问诊列表至所有客服
func (c *WebsocketController) PushTreatWaitList() {
	treats := c.Service.SearchNewTreatInfoList()
	var items []interface{}
	for _, value := range treats {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2010",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("TreatWaitList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}
//客服主动获取等待问诊列表
func (c *WebsocketController) GetTreatWaitList() {
	treats := c.Service.SearchNewTreatInfoList()
	var items []interface{}
	for _, value := range treats {
		items = append(items, value)
	}
	data, _ := json.Marshal(models.WebsocketResponse{
		Code: "2010",
		Data: models.List{
			Items: items,
		},
	})
	log.Printf("TreatWaitList: %s", data)
	_ = c.Conn.Write(1,data)
}

func (c *WebsocketController) TreatWaitListHandler(request *models.WSRequest)  {
	if !strings.EqualFold(request.Message, " getWaitList ") {
		return
	}
	c.GetTreatWaitList()
}

//2012号业务处理
//医生主动挂断语音
func (c *WebsocketController) DoctorHangOutHandler(request *models.WSRequest) {
	//更新医生状态
	doctor := c.Service.SearchClientInfo(ws_account)
	doctor.OnlineStatus = constants.Status_online
	_ = c.Service.UpdateClientInfo(doctor)
	//推送新的列表到客服
	c.PushOnlineDoctorList()
}

//2013号业务处理
//医生拒绝接听电话
func (c *WebsocketController) DoctorRejectCallHandler(request *models.WSRequest)  {
	//根据医生获取对应的未完成的问诊单,删除掉责任医生,把问诊状态重置为未处理
	treat := c.Service.SearchNotCompleteTreatInfo("doctor", ws_account)
	treat.HandleDoctor = "-"
	treat.Status = constants.Status_treat_new
	c.Service.UpdateTreatInfoStatus(treat)
	//更新客服的问诊列表
	c.PushTreatWaitList()

	//更新医生状态为空闲
	doctor := c.Service.SearchClientInfo(ws_account)
	doctor.OnlineStatus = constants.Status_online
	_ = c.Service.UpdateClientInfo(doctor)
	//推送新的列表到客服
	c.PushOnlineDoctorList()
}
