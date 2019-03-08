package services

type LoginService interface {
	ClientService
	SMSService
}

func NewLoginService() LoginService {
	return &loginService{
		NewClientService(),
		NewSMSService(),
	}
}

type loginService struct {
	ClientService
	SMSService
}