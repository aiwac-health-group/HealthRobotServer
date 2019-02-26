package services

import "HealthRobot/models"

type TokenService interface {
	UpdateToken(*models.Token) error
	GetToken(string) *models.Token

}
