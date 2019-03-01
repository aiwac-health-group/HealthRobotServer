package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
	"time"
)

type ClientService interface {
	CreatClient(*models.ClientInfo) error
	GetClient(string) *models.ClientInfo
	UpdateClient(*models.ClientInfo)
	GetOnlineDoctor() []models.ClientInfo
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
func (service *clientService) CreatClient(info *models.ClientInfo) error {
	return service.dao.Insert(info)
}

//查询用户
func (service *clientService) GetClient(account string) *models.ClientInfo {
	return service.dao.Search(account)
}

//查询在线的医生用户
func (service *clientService) GetOnlineDoctor() []models.ClientInfo {
	return service.dao.GetAll()
}

//更新用户数据
func (service *clientService) UpdateClient(info *models.ClientInfo) {
	info.UpdatedAt = time.Now()
	if err := service.dao.Update(info); err != nil {
		log.Println("client_service.go UpdateClient() insert error: ", err)
	}
}