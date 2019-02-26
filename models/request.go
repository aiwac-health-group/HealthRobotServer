package models

type LoginRequest struct {
	ClientID string `json:"clientID"`
	Password string `json:"password"`
	IdentifyCode string `json:"identifyCode"`
}