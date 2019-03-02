package services

import (
	"HealthRobotServer/dao"
	"HealthRobotServer/datasource"
	"HealthRobotServer/models"
)

type TreatInfoService interface {
	CreatTreatInfoRequest(*models.TreatInfo)
	UpdateTreatInfoStatus(*models.TreatInfo)
	GetNewTreatInfoList() []models.TreatInfo
	GetTreatInfoOnHandleList() []models.TreatInfo
	GetTreatInfoCompleteList() []models.TreatInfo
}

type treatInfoService struct {
	dao *dao.TreatInfoDao
}

func NewTreatInfoService() TreatInfoService {
	return &treatInfoService{
		dao:dao.NewTreatInfoDao(datasource.Instance()),
	}
}

func (service *treatInfoService) CreatTreatInfoRequest(TreatInfo *models.TreatInfo) {

}

func (service *treatInfoService) UpdateTreatInfoStatus(TreatInfo *models.TreatInfo) {

}

func (service *treatInfoService) GetNewTreatInfoList() []models.TreatInfo {
	return nil
}

func (service *treatInfoService) GetTreatInfoOnHandleList() []models.TreatInfo {
	return nil
}

func (service *treatInfoService) GetTreatInfoCompleteList() []models.TreatInfo {
	return nil
}
