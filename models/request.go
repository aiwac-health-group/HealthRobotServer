package models

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
	Message string `json:"message"`
	//可在下方扩展字段，比如uuid,time,要查询的字段，允许为空
}