package models

//http请求
type LoginRequest struct {
	Account string `json:"account"`
	Password string `json:"password"`
	IdentifyCode string `json:"identifyCode"`
}

type AccountAddRequest struct {
	Account string `json:"account"`
	Name string `json:"name"`
	Password string `json:"pass"`
}

type AccountModifyRequest struct {
	Account string `json:"account"`
	OperationType string `json:"type"`
	Value string `json:"value"`
}

//websocket请求
type WSRequest struct {
	BusinessCode string `json:"code"`
	Message string `json:"message,,omitempty"`
	//可在下方扩展字段，比如uuid,time,要查询的字段，允许为空
	WSRobotRequest
}

type WSRobotRequest struct {
	//公共字字段
	Account    string `json:"account,omitempty"`
	UniqueID   string `json:"uniqueID,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	Time       string `json:"time,omitempty"`
	WSRobotRequest
}

type WSRobotProfile struct {
	//个人信息
	Sex string 	`json:"sex,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Address string `json:"address,omitempty"`
	Wechat string `json:"wechat,omitempty"`
}