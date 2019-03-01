package services

type WebsocketService interface {
	ClientService
}

func NewWebsocketService() WebsocketService {
	return &websocketService{
		NewClientService(),
	}
}

type websocketService struct {
	ClientService
}