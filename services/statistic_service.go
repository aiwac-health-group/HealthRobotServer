package services
import (
	"HealthRobotServer-master/dao"
	"HealthRobotServer-master/datasource"
	"HealthRobotServer-master/models"
	"HealthRobotServer-master/constants"
	//"log"
)
type StatisticService interface{ 
	 SelectLecturework(*models.StatisticRequest) []models.ClientLecture
	 SelectReportwork(*models.StatisticRequest)  []models.ClientReport
	 SelectRegistwork(*models.StatisticRequest)  []models.ClientRegist
	 WorkerList() [] models.ServiceInfo
}


func NewStatisticService()  StatisticService {
	return &statisticService{
		dao:dao.NewDao(datasource.Instance()),
	}
}



type statisticService struct {
	dao *dao.Dao
}



//查询健康讲座工作量
func (service *statisticService) SelectLecturework(request *models.StatisticRequest) []models.ClientLecture  {
	workers := service.WorkerList()
	var clientLectures []models.ClientLecture
	for _, worker := range workers {
		var cnt int
		service.dao.Engine.Table(constants.Table_Lecture).Where("lecture.handle_service = ? AND lecture.updated_at between ? AND ?", worker.ClientAccount,request.StartTime,request.EndTime).Count(&cnt)
		var clientLecture = models.ClientLecture{
			ClientAccount:worker.ClientAccount,
			ClientName:worker.ClientName,
			CountLecture:cnt,
		}
		clientLectures = append(clientLectures, clientLecture)
	}
	return clientLectures
}

//查询健康报告工作量
func (service *statisticService) SelectReportwork(request *models.StatisticRequest) []models.ClientReport  {
	workers := service.WorkerList()
	var clientReports []models.ClientReport
	for _, worker := range workers {
		var cnt int
		service.dao.Engine.Table("examine_info").Where("examine_info.service_account = ? AND examine_info.updated_at between ? AND ?", worker.ClientAccount,request.StartTime,request.EndTime).Count(&cnt)
		var clientReport = models.ClientReport{
			ClientAccount:worker.ClientAccount,
			ClientName:worker.ClientName,
			CountReport:cnt,
		}
		clientReports = append(clientReports, clientReport)
	}

	return clientReports
}

//查询挂号工作量
func (service *statisticService) SelectRegistwork(request *models.StatisticRequest) []models.ClientRegist  {
	workers  := service.WorkerList()
	var clientRegists []models.ClientRegist
	for _, worker := range workers {
		var cnt int
		service.dao.Engine.Table("regist_info").Where("regist_info.service_account = ? AND regist_info.regist_status = ? AND lecture.updated_at between ? AND ?", worker.ClientAccount,"3",request.StartTime,request.EndTime).Count(&cnt)
		var clientRegist = models.ClientRegist{
			ClientAccount:worker.ClientAccount,
			ClientName:worker.ClientName,
			CountRegist:cnt,
		}
		clientRegists = append(clientRegists, clientRegist)
	}
	return clientRegists
}

//查询客服列表
func (service *statisticService)  WorkerList() [] models.ServiceInfo {
	var info []models.ServiceInfo
	service.dao.Engine.Table("service_info").Select("service_info.client_account,service_info.client_name").Where("service_info.client_type = ? ", "service").Find(&info)
    return info
}
