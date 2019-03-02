package models

type BaseResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type LoginResponse struct {
	LoginFlag string `json:"loginFlag"`
	ClientName string `json:"userName,omitempty"`
	ClientType string `json:"userFlag,omitempty"`
	Token string `json:"token,omitempty"`
}

//定义了机器人端接口的独有字段
type RobotResponse struct {
	Accont string `json:"account,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	UniqueID string `json:"uniqueID,omitempty"`
}

type WebsocketResponse struct {
	RobotResponse
	Code string `json:"code"`
	Data List `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	//可扩展
	RoomID string `json:"roomid,omitempty"`
}

type List struct {
	Items []interface{} `json:"items, omitempty"`
}