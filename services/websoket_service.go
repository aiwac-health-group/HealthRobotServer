package services

type WebsocketService interface {
	ClientService
	TreatInfoService
	HealthReportService
	RegisterService
	ExamineService
	LectureService
}

func NewWebsocketService() WebsocketService {
	return &websocketService{
		NewClientService(),
		NewTreatInfoService(),
		NewRegisterService(),
		NewExamineService(),
		NewHealthReportService(),
		NewLectureService(),
	}
}

type websocketService struct {
	ClientService
	TreatInfoService
	RegisterService
	ExamineService
	HealthReportService
	LectureService
}