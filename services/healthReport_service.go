package services

import (
	"HealthRobotServer/constants"
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"

)

const (
	table_reportInfo = "report_info"
)

const (
	query_reportByUser = "user_account = ?"
	query_reportByID = "id = ?"
)

type HealthReportService interface {
	CreateReportInfo(*models.HealthReport) error
	SearchReportInfo(account string) *models.HealthReport
	GetSkinInfoByAccount(account string) []models.SkinTest
	CreateSkinInfo(info *models.SkinTest) error
	GetReportByID(id int64) *models.HealthReport
	DeleteSkinInfoByAccount(account string) error
}

func NewHealthReportService() HealthReportService {
	return &healthReportService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type healthReportService struct {
	dao *dao.Dao
}

// 新增健康报告信息
func (service *healthReportService) CreateReportInfo(info *models.HealthReport) error {
	err := service.dao.Insert(table_reportInfo, info)
	if err != nil {
		log.Printf("report_service.go Fails to creat report info")
	}
	return err
}


// 查询某用户最新健康报告详细信息
func (service *healthReportService) SearchReportInfo(account string) *models.HealthReport{
	var info models.HealthReport
	_ = service.dao.SearchLast(&info, table_reportInfo, query_reportByUser, account)
	return &info
}

// 查询指定ID健康报告详细信息
func (service *healthReportService) GetReportByID(id int64) *models.HealthReport{
	var info models.HealthReport
	_ = service.dao.SearchLast(&info, table_reportInfo, query_reportByID, id)
	return &info
}

// 新增测肤数据
func (service *healthReportService) CreateSkinInfo(info *models.SkinTest) error {

	err := service.dao.Insert(constants.Table_SkinInfo, info)
	if err != nil {
		log.Printf("skin_service.go Fails to creat skin info")
	}
	return err
}

// 根据用户账号查询
func (service *healthReportService) GetSkinInfoByAccount(account string) []models.SkinTest {
	var info []models.SkinTest
	_ = service.dao.Search(&info, constants.Table_SkinInfo, constants.Query_account, account)
	return info
}

// 根据用户账号删除
func (service *healthReportService) DeleteSkinInfoByAccount(account string) error {
	if err := service.dao.Engine.Table(constants.Table_SkinInfo).Where(constants.Query_account, account).Delete(models.SkinTest{}).Error; err != nil{
		log.Println("DeleteSkinInfo error",err)
		return err
	}
	log.Println("DeleteSkinInfo successfully")
	return nil

}