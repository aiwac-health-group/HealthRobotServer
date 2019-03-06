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

//向机器人用户返回新的token
type TokenResponse struct {
	BaseResponse
	Token string `json:"token, omitempty"`
}

type WebsocketResponse struct {
	Code string `json:"code"`
	Data List `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	//可扩展
	Account    string `json:"account,omitempty"`
	UniqueID   string `json:"uniqueID,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	Link       string `json:"link,omitempty"`
	LectureContext string `json:"lectureContext,omitempty"`
	RoomID string `json:"roomid,omitempty"`
	RobotResponse
}

type WebsocketTestResponse struct {
	Code int `json:"code"`
	Data List `json:"data,omitempty"`
}

type List struct {
	Items []interface{} `json:"items, omitempty"`
}

//定义了机器人端接口的独有字段
type RobotResponse struct {
	Account string `json:"account,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	UniqueID string `json:"uniqueID,omitempty"`
}

