package services

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
)

type ClientService interface {
	SearchClientInfo(string) *models.ClientInfo
	SearchWebClientProfile(string) *models.WebClient
	SearchRobotProfile(string) *models.Robot
	GetOnlineDoctor() []models.DoctorItem
	UpdateClientInfo(*models.ClientInfo) error
	UpdateWebClientProfile(*models.WebClient) error
	UpdateRobotProfile(*models.Robot) error
}

func NewClientService()  ClientService {
	return &clientService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type clientService struct {
	dao *dao.Dao
}

//查询用户摘要信息
func (service *clientService) SearchClientInfo(account string) *models.ClientInfo {
	var info models.ClientInfo
	_ = service.dao.Search(&info, constants.Table_clientInfo, constants.Query_account, account)
	return &info
}

//获取web用户详细信息
func (service *clientService) SearchWebClientProfile(account string) *models.WebClient {
	var profile models.WebClient
	_ = service.dao.Search(&profile, constants.Table_webClient, constants.Query_account, account)
	return &profile
}

//获取robot用户详细信息
func (service *clientService) SearchRobotProfile(account string) *models.Robot  {
	var profile models.Robot
	_ = service.dao.Search(&profile, constants.Table_robot, constants.Query_account, account)
	return &profile
}

//查询在线的医生列表
func (service *clientService) GetOnlineDoctor() []models.DoctorItem {
	var doctors []models.DoctorItem
	service.dao.Engine.Table(constants.Table_webClient).Select("web_profile.client_name, web_profile.client_account, web_profile.department").Joins("inner join client_info on web_profile.client_account = client_info.client_account").Where("client_info.client_type = ? AND client_info.online_status = ?", "doctor", "2").Find(&doctors)
	return doctors
}

//添加或更新用户摘要信息
func (service *clientService) UpdateClientInfo(info *models.ClientInfo) error {
	var oldInfo models.ClientInfo
	_ = service.dao.Search(&oldInfo, constants.Table_clientInfo, constants.Query_account, info.ClientAccount)
	if oldInfo.ID == 0 {
		if err := service.dao.Insert(constants.Table_clientInfo, info); err != nil {
			log.Println("client_service.go Insert client info err")
			return err
		}
		return nil
	} else {
		if err := service.dao.Update(info, constants.Table_clientInfo, constants.Query_account, info.ClientAccount); err != nil {
			log.Println("client_service.go Update client info err")
			return err
		}
		return nil
	}
}

//添加或更新web用户个人信息
func (service *clientService) UpdateWebClientProfile(newProfile *models.WebClient) error {
	var oldProfile models.WebClient
	_ = service.dao.Search(&oldProfile, constants.Table_webClient, constants.Query_account, newProfile.ClientAccount)
	if oldProfile.ID == 0 { //不存在该工作人员的详细信息,则插入详细信息
		log.Println("the client's profile does not exist, try to Insert")
		if err := service.dao.Insert(constants.Table_webClient, newProfile); err != nil {
			log.Println("client_service.go Insert doctor profile err")
			return err
		}
		return nil
	} else { //否则更新医生详细信息
		if err := service.dao.Update(newProfile, constants.Table_webClient, constants.Query_account, newProfile.ClientAccount); err != nil {
			log.Println("client_service.go Update doctor profile err")
			return err
		}
		return nil
	}
}

//注册或修改机器人信息
func (service *clientService) UpdateRobotProfile(newProfile *models.Robot) error {
	var oldProfile models.Robot
	_ = service.dao.Search(&oldProfile, constants.Table_robot, constants.Query_account, newProfile.ClientAccount)
	if oldProfile.ID == 0 { //不存在该机器人的详细信息,则插入详细信息
		if err := service.dao.Insert(constants.Table_robot, newProfile); err != nil {
			log.Println("client_service.go Insert robot profile err")
			return err
		}
		return nil
	} else { //否则更新机器人详细信息
		if err := service.dao.Update(newProfile, constants.Table_robot, constants.Query_account, newProfile.ClientAccount); err != nil {
			log.Println("client_service.go Update robot profile err")
			return err
		}
		return nil
	}
}

