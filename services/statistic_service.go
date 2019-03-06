package services
import (
	"HealthRobotServer-master/dao"
	"HealthRobotServer-master/datasource"
	"HealthRobotServer-master/models"
	"HealthRobotServer-master/constants"
	//"log"
	//"HealthRobotServer-master/constants"
)
type StatisticService interface{ 
	 SelectLecturework(*models.StatisticRequest) []models.ClientLecture
	 SelectReportwork(*models.StatisticRequest)  []models.ClientReport
	 SelectRegistwork(*models.StatisticRequest)  []models.ClientRegist
	 WorkerList() [] models.Client
}


func NewStatisticService()  LectureService {
	return &lectureService{
		dao:dao.NewDao(datasource.Instance()),
	}
}



type statisticService struct {
	dao *dao.Dao
}



//查询健康讲座工作量
func (service *statisticService) SelectLecturework(request *models.StatisticRequest) []models.ClientLecture  {
	worker := service.WorkerList()
	var clientLecture []models.ClientLecture
	for i:= 0;i<len(worker);i++{
		account:=worker[i].ClientAccount
		var num int
		service.dao.Engine.Table(constants.Table_Lecture).Where("lecture.Service = ? AND lecture.updated_at between ? AND ?", account,request.StartTime,request.EndTime).Count(&num)
		clientLecture[i].ClientAccount = worker[i].ClientAccount
		clientLecture[i].ClientName = worker[i].ClientName
		clientLecture[i].CountLecture = num
	}
	return clientLecture
}

//查询健康报告工作量
func (service *statisticService) SelectReportwork(request *models.StatisticRequest) []models.ClientReport  {
	worker := service.WorkerList()
	var clientReport []models.ClientReport
	for i:= 0;i<len(worker);i++{
		account:=worker[i].ClientAccount
		var num int
		service.dao.Engine.Table("examine_info").Where("examine_info.service_account = ? AND examine_info.updated_at between ? AND ?", account,request.StartTime,request.EndTime).Count(&num)
		clientReport[i].ClientAccount = worker[i].ClientAccount
		clientReport[i].ClientName = worker[i].ClientName
		clientReport[i].CountReport= num
	}
	return clientReport
}

//查询挂号工作量
func (service *statisticService) SelectRegistwork(request *models.StatisticRequest) []models.ClientRegist  {
	worker := service.WorkerList()
	var clientRegist []models.ClientRegist
	for i:= 0;i<len(worker);i++{
		account:=worker[i].ClientAccount
		var num int
		service.dao.Engine.Table("regist_info").Where("regist_info.service_account = ? AND regist_info.regist_status = ? AND lecture.updated_at between ? AND ?", account,"3",request.StartTime,request.EndTime).Count(&num)
		clientRegist[i].ClientAccount = worker[i].ClientAccount
		clientRegist[i].ClientName = worker[i].ClientName
		clientRegist[i].CountRegist = num
	}
	return clientRegist
}

//查询客服列表
func (service *statisticService)  WorkerList() [] models.Client {
	var info []models.Client
	service.dao.Engine.Table(constants.Table_webClient).Select("web_profile.client_account,web_profile.client_name,").Where("web_profile.client_type = ? ", "service").Find(&info)
    return info
}
