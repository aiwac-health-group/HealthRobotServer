package util

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"HealthRobotServer/models"
	"HealthRobotServer/manager"
	"log"
)

//0018号业务处理
//新消息通知
func  MessageNoticeHandler(notice *models.MessageNotice, WsManager manager.WSManager, Conn websocket.Connection) {
	if(notice.MessageType == "0") {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code: "0018",
			RobotResponse:models.RobotResponse{
				MessageNotice:*notice,
			},
		})
		log.Printf("MessageNotice0: %s", data)
		_ = Conn.To("robot").EmitMessage(data)
	} else {
		data, _ := json.Marshal(models.WebsocketResponse{
			Code: "0018",
			RobotResponse:models.RobotResponse{
				MessageNotice:*notice,
			},
		})
		log.Printf("MessageNotice: %s", data)
		RobotConn := WsManager.GetWSConnection(notice.UserAccount)
		_ = (*RobotConn).Write(1,data)
	}
}