package services

type WebsocketService interface {
	ClientService
	TreatInfoService
	LectureService
}

func NewWebsocketService() WebsocketService {
	return &websocketService{
		NewClientService(),
		NewTreatInfoService(),
		NewLectureService(),
	}
}

type websocketService struct {
	ClientService
	TreatInfoService
	LectureService
}