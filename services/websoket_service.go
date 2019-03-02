package services

type WebsocketService interface {
	ClientService
	TreatInfoService
}

func NewWebsocketService() WebsocketService {
	return &websocketService{
		NewClientService(),
		NewTreatInfoService(),
	}
}

type websocketService struct {
	ClientService
	TreatInfoService
}