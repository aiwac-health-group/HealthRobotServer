package models

type LoginRequest struct {
	ClientID string `json:"clientID"`
	Password string `json:"password"`
	IdentifyCode string `json:"identifyCode"`
}

type ServiceAddRequest struct {
	Account string `json:"account"`
	Name string `json:"name"`
	Password string `json:"pass"`
}