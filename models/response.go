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

type WebsocketResponse struct {
	Code int `json:"code"`
	Data List `json:"data"`
}

type List struct {
	Items []ClientInfo `json:"items"`
}