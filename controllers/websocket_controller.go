package controllers

import (
	"HealthRobotServer/middleware"
	"HealthRobotServer/models"
	"HealthRobotServer/services"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"log"
)

type WebsocketController struct {
	Ctx iris.Context
	Conn websocket.Connection
	Service services.WebsocketService
}

func (c *WebsocketController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Router().Use(middleware.JwtHandlerWS().Serve, middleware.NewAuthToken().Serve)
	b.Handle("GET","/","Join")
}

var (
	account string
	clientType string
)

func (c *WebsocketController) Join() {
	//从token中提取用户信息
	token := c.Ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	account = claims["Account"].(string)
	clientType = claims["ClientType"].(string)
	log.Println("New Websocket Connection: ",account, clientType)

	//加入对应clientType的room, 每个room存放了相应用户类型的所有websocket连接
	c.Conn.Join(clientType)
	//注册连接断开回调函数
	c.Conn.OnDisconnect(c.LoseConnection)
	//注册消息接收处理函数
	c.Conn.OnMessage(c.ReceiveRequest)

	//更新用户状态
	c.Service.UpdateClient(&models.ClientInfo{
		ClientAccount:account,
		OnlineStatus:"2",
	})

	//如果上线用户为医生
	//获取在线医生列表,并将列表推送给所有的客服
	if clientType == "doctor" {
		doctors := c.Service.GetOnlineDoctor()
		data, _ := json.Marshal(models.WebsocketResponse{
			Code: "2009",
			Data: models.List{
				Items: doctors,
			},
		})
		log.Printf("response data %s: ", data)
		_ = c.Conn.To("service").EmitMessage(data)
	}

	//开启事件监听
	c.Conn.Wait()
}

func (c *WebsocketController) LoseConnection() {
	//更新用户状态
	c.Service.UpdateClient(&models.ClientInfo{
		ClientAccount:account,
		OnlineStatus:"1",
	})
	//如果离开的用户为doctor，则更新客服的在线医生列表
	if clientType == "doctor" {
		doctors := c.Service.GetOnlineDoctor()
		data, _ := json.Marshal(models.WebsocketResponse{
			Code: "2009",
			Data: models.List{
				Items: doctors,
			},
		})
		log.Printf("response data %s: ", data)
		_ = c.Conn.To("service").EmitMessage(data)
	}
	log.Printf("%s %s lose the connection", account, clientType)
}

func (c *WebsocketController) ReceiveRequest(data []byte) {
	//在这里解析收到的请求，根据请求中的业务号跳转到指定业务处理函数中进行处理

}

