package websocket

type WSManager struct {
	//内存中保存在线用户列表
	clients map[*WSClient]bool
	//连接请求消息写入注册管道
	register chan *WSClient
	//连接断开消息写入解绑管道
	unregister chan *WSClient
}

func NewWSManager() *WSManager {
	return &WSManager{
		clients:make(map[*WSClient]bool),
		register:make(chan *WSClient),
		unregister:make(chan *WSClient),
	}
}

func (manager *WSManager) start() {
	for true {
		select {
		//新的用户连接
		case client := <- manager.register:
			manager.clients[client] = true
			//医生上线，把在线医生列表推送给客服
			if client.ClientType == "doctor" {
				manager.pushOnlineDoctorList()
			}
		//用户断开连接
		case leave := <- manager.unregister:
			manager.clients[leave] = false
			//医生下线，重新更新在线医生列表至客服
			if leave.ClientType == "doctor" {
				manager.pushOnlineDoctorList()
			}

		}
	}
}

//推送在线医生列表
func (manager *WSManager) pushOnlineDoctorList() {

}
