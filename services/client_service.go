package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
)

const (
	outline = "1"
	online = "2"
	onbusy = "3"
)

const (
	table_clientInfo = "client_info"
	table_doctor = "doctor_profile"
	table_robot = "robot_profile"
)

const  (
	query_account = "client_account = ?"
	query_online_status = "online_status = ? AND client_type = ?"
)

type ClientService interface {
	CreatClientInfo(*models.ClientInfo) error
	SearchClientInfo(string) *models.ClientInfo
	UpdateClientInfo(*models.ClientInfo) error
	GetOnlineDoctor() []models.ClientInfo
	UpdateDoctor(*models.Doctor) error
	UpdateRobot(*models.Robot) error
}

func NewClientService()  ClientService {
	return &clientService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type clientService struct {
	dao *dao.Dao
}

//创建用户摘要数据
func (service *clientService) CreatClientInfo(info *models.ClientInfo) error {
	err := service.dao.Insert(table_clientInfo, info)
	if err != nil {
		log.Printf("client_service.go Fails to creat client info")
	}
	return err
}

//查询用户摘要信息
func (service *clientService) SearchClientInfo(account string) *models.ClientInfo {
	var info models.ClientInfo
	condition := []interface{}{1:account}
	_ = service.dao.Search(&info, table_clientInfo, query_account, condition)
	return &info
}

//查询在线的医生列表
func (service *clientService) GetOnlineDoctor() []models.ClientInfo {
	var doctors []models.ClientInfo
	service.dao.GetList(doctors,table_clientInfo,query_online_status,online,"doctor")
	return doctors
}

//更新用户摘要信息
func (service *clientService) UpdateClientInfo(newInfo *models.ClientInfo) error {
	service.dao.TransactionBegin()
	err := service.dao.Update(newInfo, table_clientInfo, query_account, newInfo.ClientAccount)
	if err != nil {
		service.dao.RollBack()
		return err
	}
	var newClientInfo = &models.ClientInfo{
		ClientName:newInfo.ClientName,
	}
	err = service.dao.Update(newClientInfo, table_clientInfo, query_account, newInfo.ClientAccount)
	if err != nil {
		service.dao.RollBack()
		return err
	}
	service.dao.Commit()
	return err
}

//修改医生信息
func (service *clientService) UpdateDoctor(newProfile *models.Doctor) error {
	var oldProfile *models.Doctor
	_ = service.dao.Search(oldProfile, table_doctor, query_account, newProfile.ClientAccount)
	if oldProfile == nil { //不存在该医生的详细信息,则插入详细信息
		_ = service.dao.Insert(table_doctor, newProfile)
	} else { //否则更新医生详细信息，同时更新用户的摘要信息
		service.dao.TransactionBegin()
		err := service.dao.Update(newProfile, table_doctor, query_account, newProfile.ClientAccount)
		if err != nil {
			service.dao.RollBack()
			return err
		}
		var newClientInfo = &models.ClientInfo{
			ClientName:newProfile.ClientName,
		}
		err = service.dao.Update(newClientInfo, table_clientInfo, query_account, newProfile.ClientAccount)
		if err != nil {
			service.dao.RollBack()
			return err
		}
		service.dao.Commit()
	}
}

//注册或修改机器人信息
func (service *clientService) UpdateRobot(newProfile *models.Robot) error {
	var oldProfile *models.Doctor
	_ = service.dao.Search(oldProfile, table_robot, query_account, newProfile.ClientAccount)
	if oldProfile == nil { //不存在该机器人的详细信息,则插入详细信息
		_ = service.dao.Insert(table_robot, newProfile)
	} else { //否则更新机器人详细信息，同时更新用户的摘要信息
		service.dao.TransactionBegin()
		err := service.dao.Update(newProfile, table_robot, query_account, newProfile.ClientAccount)
		if err != nil {
			service.dao.RollBack()
			return err
		}
		var newClientInfo = &models.ClientInfo{
			ClientName:newProfile.ClientName,
		}
		err = service.dao.Update(newClientInfo, table_robot, query_account, newProfile.ClientAccount)
		if err != nil {
			service.dao.RollBack()
			return err
		}
		service.dao.Commit()
	}
}