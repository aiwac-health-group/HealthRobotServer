package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"time"
)

type TokenService interface {
	UpdateToken(*models.Token) error
	GetToken(string) *models.Token
}

func NewTokenService() TokenService {
	return &tokenService{
		dao:dao.NewTokenDao(datasource.Instance()),
	}
}

type tokenService struct {
	dao *dao.TokenDao
}

func (service *tokenService) GetToken(raw string) *models.Token {
	return service.dao.Search(raw)
}

func (service *tokenService) UpdateToken(token *models.Token) error {
	token.ExpressIn = time.Now().AddDate(0,0,1).Unix()
	token.UpdatedAt = time.Now()
	return service.dao.Update(token)
}
