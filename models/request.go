package models

//http请求
type LoginRequest struct {
	Account string `json:"account"`
	Password string `json:"password"`
	IdentifyCode string `json:"identifyCode"`
}

type TokenGetRequest struct {
	Account string `json:"account"`
	OldToken string `json:"oldToken"`
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

type DoctorProfileModifyRequest struct {
	Account string `json:"account"`
	Name string `json:"name"`
	Department string `json:"class"`
	Brief string `json:"blief"`
}

//websocket请求
type WSRequest struct {
	BusinessCode string `json:"code"`
	//可在下方扩展字段，比如uuid,time,要查询的字段，允许为空
	WSRobotRequest
	WSWebRequest
}

//Web工作人员的websocket请求字段
type WSWebRequest struct {
	Message string `json:"message,,omitempty"`
}

//机器人端的websocket请求字段
type WSRobotRequest struct {
	//公共字字段
	Account    string `json:"account,omitempty"`
	UniqueID   string `json:"uniqueID,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	Time       string `json:"time,omitempty"`
	WSRobotProfile
}

type WSRobotProfile struct {
	//个人信息
	Name string `json:"name"`
	Sex string 	`json:"sex,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Address string `json:"address,omitempty"`
	Wechat string `json:"wechat,omitempty"`
}