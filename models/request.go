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