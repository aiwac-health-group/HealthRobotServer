package models

import "time"

type BaseResponse struct {
	Status string `json:"status"`
	Message string `json:"message, omitempty"`
	Data List `json:"data,omitempty"`
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
	RobotResponse
	WebResponse
}

type List struct {
	Items []interface{} `json:"items, omitempty"`
}

//定义了Web端接口返回字段
type WebResponse struct {
	TreatResponse
}

//定义了机器人端接口的返回字段
type RobotResponse struct {
	Account string `json:"account,omitempty"`
	ClientType string `json:"clientType,omitempty"`
	UniqueID string `json:"uniqueID,omitempty"`
	LectureResponse
	MessageNotice
	HealthReportResponse
	RegistResponse
	TreatResponse
	ExamineResponse
}

// 消息通知
type MessageNotice struct {
	MessageType string `json:"messageType,omitempty"`
	MessageID	int64  `json:"messageID,omitempty"`
	ExamCover	string `json:"cover,omitempty"`
	UserAccount string `json:"userAccount,omitempty"`
}

// 语音问诊
type TreatResponse struct {
	RoomID string `json:"roomid,omitempty"`
}

//挂号记录
type RegistResponse struct {
	RegistInfo  	Regist
	Description		string		`json:"description,omitempty"`
	RegisterStatus  string		`json:"registerStatus,omitempty"`
	UpdateTime		time.Time	`json:"updateTime,omitempty"`
	CreateTime		time.Time	`json:"createTime,omitempty"`
}

//健康报告
type HealthReportResponse struct {
	Report 			string `json:"resultContext,omitempty"`
}

//体检反馈
type ExamineResponse struct {
	Examine 		string `json:"examContext,omitempty"`
	ExaminePackage 	string `json:"link,omitempty"` // 体检套餐文件链接
}

//健康讲座
type LectureResponse struct {
	Link string `json:"link,omitempty"`
	LectureContext string `json:"lectureContext,omitempty"`
}