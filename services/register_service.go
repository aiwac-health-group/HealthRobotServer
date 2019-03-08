package services

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
)

const (
	table_registInfo = "regist_info"
)

const  (
	query_regist_id = "id = ?"
	query_service_processing= "service_account = ? AND regist_status = ?"
	query_service_account = "service_account = ?"
	query_regist_status = "regist_status = ?"
	query_user_account = "user_account = ?"
)

type RegisterService interface {
	CreateRegistInfo(*models.Registration) error
	QueryRegistrationIdle(id int64, handleService string) (bool, string)
	QueryRegistrationProcessing(a string) *models.Registration
	SaveRegistResult(in *models.RegistProcessResponse) error
	UpdateRegistStatus(id int64, handleService string)  error
	GetNonresponseRegist() []models.Registration
	SearchRegistByID(id int64) *models.Registration
	SearchRegistByUser(account string) []models.Registration
}

func NewRegisterService()  RegisterService {
	return &registerService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type registerService struct {
	dao *dao.Dao
}

// 新增挂号信息
func (service *registerService) CreateRegistInfo(info *models.Registration) error {
	err := service.dao.Insert(table_registInfo, info)
	if err != nil {
		log.Printf("regist_service.go Fails to creat regist info")
	}
	return err
}

// 查询指定id的挂号请求能否被指定客服处理
func (service *registerService) QueryRegistrationIdle(id int64, handleService string) (bool, string) {
	var info models.Registration
	// 指定id的挂号
	_ = service.dao.Search(&info, table_registInfo, query_regist_id, id)
	switch s := info.Status; s {
	case constants.Status_processing:
		return false, "该挂号请求正在处理中"
	case constants.Status_completed:
		return false, "该挂号请求已经处理完成"
	case constants.Status_noresponse:

		forSearch := models.Registration {
			//HandleSevice: handleService,
			//Status: "2",
		}

		count := 0
		err :=  service.dao.Engine.Table(table_registInfo).Where("service_account = ? AND regist_status = ?", handleService, "2").Find(&forSearch).Count(&count)
		if err != nil {
			log.Println("查询该客服处理中挂号单有误", err)
			return false ,"查询该客服处理中挂号单有误"
		}
		if count != 0 {
			return false, "该客服有未处理完的挂号请求"
		}
	}
	_ = service.UpdateRegistStatus(id, handleService)
	return true, ""
}

func (service *registerService) UpdateRegistStatus(id int64, handleService string)  error{
	info := models.Registration{
		Status: "2",
		HandleSevice: handleService,
	}

	err := service.dao.Update(info, table_registInfo, query_regist_id, id)
	return err
}

// 获取指定客服处理中的挂号信息
func (service *registerService) QueryRegistrationProcessing(account string) *models.Registration{
	var info models.Registration
	_ =  service.dao.Search(&info, table_registInfo, query_service_processing, account, "2")
	log.Println(info)
	return &info
}

//更新挂号结果
func (service *registerService) SaveRegistResult(in *models.RegistProcessResponse) error {
	service.dao.TransactionBegin()
	var newRegistInfo = &models.Registration{
		RegistDesc: in.ResponseDesc,
		Status:"3",
	}
	err := service.dao.Update(newRegistInfo, table_registInfo, query_regist_id, in.ID)
	if err != nil {
		service.dao.RollBack()
		return err
	}
	service.dao.Commit()
	return err

}

// 查询未处理的挂号信息
func (service *registerService) GetNonresponseRegist() []models.Registration {
	var regist []models.Registration
	service.dao.GetList(&regist, table_registInfo, query_regist_status, constants.Status_noresponse)
	return regist
}


// 用id查询挂号结果
func (service *registerService) SearchRegistByID(id int64) *models.Registration{
	var info models.Registration
	_ = service.dao.Search(&info, table_registInfo, query_regist_id, id)
	return &info
}

// 用用户账号查询挂号记录
func (service *registerService) SearchRegistByUser(account string) []models.Registration{
	var info []models.Registration
	_ = service.dao.Search(info, table_registInfo, query_user_account, account)
	return info
}
