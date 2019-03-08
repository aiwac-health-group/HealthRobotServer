package services

import (
	"HealthRobotServer-master/constants"
	"HealthRobotServer-master/dao"
	"HealthRobotServer-master/datasource"
	"HealthRobotServer-master/models"
	"log"
)

type ClientService interface {
	SearchRobotClientInfo(string) *models.RobotInfo
	SearchDoctorClientInfo(string) *models.DoctorInfo
	SearchServiceClientInfo(string) *models.ServiceInfo

	GetOnlineDoctor() []models.DoctorItem
	GetRobotListForHealthReport(upBoundary int64, downBoundary int64) []models.RobotItem

	UpdateRobotClientInfo(*models.RobotInfo) error
	UpdateDoctorClientInfo(*models.DoctorInfo) error
	UpdateServiceClientInfo(*models.ServiceInfo) error

	CountTotalRobotClient() int
	CountTotalDoctorClient() int

}

func NewClientService()  ClientService {
	return &clientService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type clientService struct {
	dao *dao.Dao
}

//查询机器人用户信息
func (service *clientService) SearchRobotClientInfo(account string) *models.RobotInfo {
	var info models.RobotInfo
	if err := service.dao.Engine.Table("robot_info").Where("client_account = ?", account).Find(&info).Error; err != nil {
		log.Println("client_service.go SearchRobotClientInfo() err: ", err)
	}
	return &info
}

//查询医生用户信息
func (service *clientService) SearchDoctorClientInfo(account string) *models.DoctorInfo {
	var info models.DoctorInfo
	if err := service.dao.Engine.Table("doctor_info").Where("client_account = ?", account).Find(&info).Error; err != nil {
		log.Println("client_service.go SearchDoctorClientInfo() err: ", err)
	}
	return &info
}

//查询客服用户信息
func (service *clientService) SearchServiceClientInfo(account string) *models.ServiceInfo {
	var info models.ServiceInfo
	if err := service.dao.Engine.Table("service_info").Where("client_account = ?", account).Find(&info).Error; err != nil {
		log.Println("client_service.go SearchServiceClientInfo() err: ", err)
	}
	return &info
}

//查询在线的医生列表
func (service *clientService) GetOnlineDoctor() []models.DoctorItem {
	var doctors []models.DoctorItem
	service.dao.Engine.Table("doctor_info").Select("doctor_info.client_name, doctor_info.client_account, doctor_info.department").Where("doctor_info.online_status = ?", "2").Find(&doctors)
	return doctors
}

//添加或更新机器人用户摘要信息
func (service *clientService) UpdateRobotClientInfo(info *models.RobotInfo) error {
	var oldInfo models.RobotInfo
	_ = service.dao.Search(&oldInfo, "robot_info", constants.Query_account, info.ClientAccount)
	if oldInfo.ID == 0 {
		if err := service.dao.Insert("robot_info", info); err != nil {
			log.Println("client_service.go Insert robot info err: ", err)
			return err
		}
		return nil
	} else {
		if err := service.dao.Update(info, "robot_info", constants.Query_account, info.ClientAccount); err != nil {
			log.Println("client_service.go Update robot info err: ", err)
			return err
		}
		return nil
	}
}

//添加或更新医生用户摘要信息
func (service *clientService) UpdateDoctorClientInfo(info *models.DoctorInfo) error {
	var oldInfo models.DoctorInfo
	_ = service.dao.Search(&oldInfo, "doctor_info", constants.Query_account, info.ClientAccount)
	if oldInfo.ID == 0 {
		if err := service.dao.Insert("doctor_info", info); err != nil {
			log.Println("client_service.go Insert doctor info err: ", err)
			return err
		}
		return nil
	} else {
		if err := service.dao.Update(info, "doctor_info", constants.Query_account, info.ClientAccount); err != nil {
			log.Println("client_service.go Update doctor info err: ", err)
			return err
		}
		return nil
	}
}

//添加或更新客服用户摘要信息
func (service *clientService) UpdateServiceClientInfo(info *models.ServiceInfo) error {
	var oldInfo models.ServiceInfo
	_ = service.dao.Search(&oldInfo, "service_info", constants.Query_account, info.ClientAccount)
	if oldInfo.ID == 0 {
		if err := service.dao.Insert("service_info", info); err != nil {
			log.Println("client_service.go Insert service info err: ", err)
			return err
		}
		return nil
	} else {
		if err := service.dao.Update(info, "service_info", constants.Query_account, info.ClientAccount); err != nil {
			log.Println("client_service.go Update service info err: ", err)
			return err
		}
		return nil
	}
}

//获取数据库中某一类用户的总数
func (service *clientService) CountTotalRobotClient() int {
	var count int
	err := service.dao.Engine.Table("robot_info").Count(&count).Error
	if err != nil {
		log.Println("client_service.go CountTotalRobotClient() err: ", err)
		return 0
	}
	return count
}

func (service *clientService) CountTotalDoctorClient() int {
	var count int
	err := service.dao.Engine.Table("robot_info").Count(&count).Error
	if err != nil {
		log.Println("client_service.go CountTotalDoctorClient() err: ", err)
		return 0
	}
	return count
}

//根据ID范围搜索需要处理健康报告的用户
func (service *clientService) GetRobotListForHealthReport(upBoundary int64, downBoundary int64) []models.RobotItem {
	var robots []models.RobotItem
	err := service.dao.Engine.Table("robot_info").Select("robot_info.client_account").Where("robot_info.id < ? AND robot_info.id >= ? AND health_status = ?", upBoundary, downBoundary, "1").Find(&robots).Error
	if err != nil {
		log.Println("client_service.go CountTotalDoctorClient() err: ", err)
		return nil
	}
	return robots
}