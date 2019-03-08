package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
	"log"
)

const (
	table_examineInfo = "examine_info"
)

const (
	query_examine = "where id = ?"
)

type ExamineService interface {
	CreatePhysicalExamine(*models.PhysicalExamine) error
	GetAllExamine() []models.PhysicalExamine
	SearchExamineInfo(id int64) *models.PhysicalExamine
	Get3NewExamine() []models.PhysicalExamine
	GetTheNewExamine() *models.PhysicalExamine
}

func NewExamineService()  ExamineService {
	return &examineService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

type examineService struct {
	dao *dao.Dao
}

// 新增体检推荐信息
func (service *examineService) CreatePhysicalExamine(info *models.PhysicalExamine) error {
	err := service.dao.Insert(table_examineInfo, info)
	if err != nil {
		log.Printf("examine_service.go Fails to creat examine info")
	}
	return err
}

// 返回体检推荐列表
func (service *examineService) GetAllExamine() []models.PhysicalExamine {
	var list []models.PhysicalExamine
	service.dao.GetAllList(&list, table_examineInfo)
	return list
}

// 返回最新3条体检推荐列表
func (service *examineService) Get3NewExamine() []models.PhysicalExamine {
	var list []models.PhysicalExamine
	service.dao.GetAllList(&list, table_examineInfo)
	len := len(list)
	var resultList []models.PhysicalExamine
	for i := len - 1; i >= len - 3; i-- {
		resultList = append(resultList, list[i])
	}
	return resultList
}

// 查询体检推荐详细信息
func (service *examineService) SearchExamineInfo(id int64) *models.PhysicalExamine{
	var info models.PhysicalExamine
	_ = service.dao.SearchLast(&info, table_examineInfo, query_examine, id)
	return &info
}

//查询最新的体检推荐
func (service *examineService) GetTheNewExamine() *models.PhysicalExamine {
	var info models.PhysicalExamine
	if err := service.dao.Engine.Table(table_examineInfo).Last(&info).Error; err != nil {
		log.Println("GetTheNewExamine erro",err)
	}
	return &info
}


