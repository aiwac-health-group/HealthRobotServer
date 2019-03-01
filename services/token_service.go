package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
	"time"
)

type TokenService interface {
	CreatToken(*models.Token)
	UpdateToken(*models.Token)
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

func (service *tokenService) CreatToken(token *models.Token) {
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.ExpressIn = time.Now().AddDate(0,0,1).Unix()
	if err := service.dao.Insert(token); err != nil {
		log.Println("token_service.go InsertToken() insert error: ", err)
	}
}

func (service *tokenService) GetToken(account string) *models.Token {
	return service.dao.Search(account)
}

func (service *tokenService) UpdateToken(token *models.Token) {
	token.ExpressIn = time.Now().AddDate(0,0,1).Unix()
	token.UpdatedAt = time.Now()
	if err := service.dao.Update(token); err != nil {
		log.Println("token_service.go InsertToken() insert error: ", err)
	}
}
