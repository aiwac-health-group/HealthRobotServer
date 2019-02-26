package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
)

type ClientService interface {
	Creat(info *models.ClientInfo) error
	GetByName(string) *models.ClientInfo
}

func NewClientService()  ClientService {
	return &clientService{
		dao:dao.NewClientDao(datasource.Instance()),
	}
}

type clientService struct {
	dao *dao.ClientDao
}

//创建新的用户数据
func (service *clientService) Creat(info *models.ClientInfo) error {
	return service.dao.Insert(info)
}

//查询用户
func (service *clientService) GetByName(name string) *models.ClientInfo {
	return service.dao.Search(name)
}