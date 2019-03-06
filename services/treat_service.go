package services

import (
	"HealthRobotServer-master/constants"
	"HealthRobotServer-master/dao"
	"HealthRobotServer-master/datasource"
	"HealthRobotServer-master/models"
	"strings"
)

type TreatInfoService interface {
	CreatTreatInfoRequest(*models.TreatInfo)
	UpdateTreatInfoStatus(*models.TreatInfo)
	UpdateTreatInfoHandleDoctor(*models.TreatInfo)
	SearchNewTreatInfo(string) *models.TreatInfo
	SearchNotCompleteTreatInfo(string, string) *models.TreatInfo
	SearchNewTreatInfoList() []models.TreatInfo
	SearchTreatInfoOnHandleList() []models.TreatInfo
	SearchTreatInfoCompleteList() []models.TreatInfo
}

type treatInfoService struct {
	dao *dao.Dao
}

func NewTreatInfoService() TreatInfoService {
	return &treatInfoService{
		dao:dao.NewDao(datasource.Instance()),
	}
}

func (service *treatInfoService) CreatTreatInfoRequest(treatInfo *models.TreatInfo) {

}

func (service *treatInfoService) UpdateTreatInfoHandleDoctor(treatInfo *models.TreatInfo) {
	_ = service.dao.Update(treatInfo, constants.Table_TreatInfo, constants.Query_account, treatInfo.Account)
}

func (service *treatInfoService) UpdateTreatInfoStatus(treatInfo *models.TreatInfo) {

}

func (service *treatInfoService) SearchNewTreatInfo(account string) *models.TreatInfo {
	var treat models.TreatInfo
	_ = service.dao.Search(&treat, constants.Table_TreatInfo, constants.Query_treat, account, constants.Status_treat_new)
	return &treat
}

func (service *treatInfoService) SearchNotCompleteTreatInfo(column string, account string) *models.TreatInfo {
	var treat models.TreatInfo
	if strings.EqualFold(column, "patient") {
		_ = service.dao.Search(&treat, constants.Table_TreatInfo, "client_account <> ?", account, constants.Status_treat_complete)
	} else {
		_ = service.dao.Search(&treat, constants.Table_TreatInfo, "handle_doctor <> ?", account, constants.Status_treat_complete)
	}

	return &treat
}

func (service *treatInfoService) SearchNewTreatInfoList() []models.TreatInfo {
	var treats []models.TreatInfo
	service.dao.GetList(&treats,constants.Table_TreatInfo,constants.Query_treat_status,constants.Status_treat_new)
	return treats
}

func (service *treatInfoService) SearchTreatInfoOnHandleList() []models.TreatInfo {
	return nil
}

func (service *treatInfoService) SearchTreatInfoCompleteList() []models.TreatInfo {
	return nil
}
