package controllers

import (
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
	"unsafe"
)

const (
	BusinessTreatRequest = 17 //机器人发起问诊请求
	BusinessOnlineDoctor = 2009 //客服获取在线医生列表
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
	account string
	clientType string
	clientName string
)

func (c *WebsocketController) Join() {
	//从token中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	account = claims["Account"].(string)
	clientType = claims["ClientType"].(string)
	//clientName = claims["ClientName"].(string)
	log.Println("New Websocket Connection: ",account, clientType, clientName)

	//加入对应clientType的room, 每个room存放了相应用户类型的所有websocket连接
	c.Conn.Join(clientType)
	//注册连接断开回调函数
	c.Conn.OnDisconnect(c.LoseConnection)
	//注册消息接收处理函数
	c.Conn.OnMessage(c.ReceiveRequest)
	//存储该用户和对应连接的映射关系
	c.WsManager.AddMapRelationship(account,&(c.Conn))

	//更新用户状态为在线
	c.Service.UpdateClient(&models.ClientInfo{
		ClientAccount:account,
		OnlineStatus:"2",
	})

	//如果上线用户为医生
	//获取在线医生列表,并将列表推送给所有的客服
	if clientType == "doctor" {
		c.SendOnlineDoctorList()
	}

	//开启事件监听
	c.Conn.Wait()
}

func (c *WebsocketController) LoseConnection() {
	//删除用户与连接的映射关系
	c.WsManager.DeleteMapRelationship(account)

	//更新用户状态
	c.Service.UpdateClient(&models.ClientInfo{
		ClientAccount:account,
		OnlineStatus:"1",
	})

	//如果离开的用户为doctor，则更新客服的在线医生列表
	if clientType == "doctor" {
		c.SendOnlineDoctorList()
	}
	log.Printf("%s %s lose the connection", account, clientType)
}

func (c *WebsocketController) ReceiveRequest(data []byte) {
	//在这里解析收到的请求，根据请求中的业务号跳转到指定业务处理函数中进行处理
	var request models.WSRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Println("Websocket request from Explore Unmarshal err: ",err)
	}
	log.Println("Websocket request from Explore: ",request)

	if unsafe.Sizeof(request) == 0 { //空请求
		return
	}

	//根据request中的code字段配置相应的函数处理
	businessCode, _ := strconv.Atoi(request.BusinessCode)

	switch businessCode {
	case BusinessOnlineDoctor: c.DoctorListRequestHandler(&request)
	case BusinessTreatRequest: c.TreatRequestHandler(&request)

	}
}

//2009号业务处理
//获取在线医生账号
func (c *WebsocketController) SendOnlineDoctorList() {
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
	log.Printf("doctorList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}

func (c *WebsocketController) DoctorListRequestHandler(request *models.WSRequest) {
	c.SendOnlineDoctorList()
}


//17号业务处理
//处理用户发起的问诊请求
//把问诊请求存进数据库，并把等待列表推送给在线客服
func (c *WebsocketController) TreatRequestHandler(request *models.WSRequest) {
	c.Service.CreatTreatInfoRequest(&models.TreatInfo{
		Account:account,
		ClientName:clientName,
	})
	treats := c.Service.GetNewTreatInfoList()
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
	log.Printf("doctorList: %s", data)
	_ = c.Conn.To("service").EmitMessage(data)
}

