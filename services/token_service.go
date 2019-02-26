package services

import "HealthRobotServer/models"

type TokenService interface {
	UpdateToken(*models.Token) error
	GetToken(string) *models.Token

}
