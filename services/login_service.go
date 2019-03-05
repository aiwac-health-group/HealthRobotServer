package services

type LoginService interface {
	ClientService
}

func NewLoginService() LoginService {
	return &loginService{
		NewClientService(),
	}
}

type loginService struct {
	ClientService
}