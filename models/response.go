package models

type BaseResponse struct {
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}

type LoginResponse struct {
	BaseResponse
	ClientType string `json:"clientType,omitempty"`
	Token string `json:"token,omitempty"`
}