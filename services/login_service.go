package services

type LoginService interface {
	ClientService
	TokenService
}

func NewLoginService() LoginService {
	return &loginService{
		NewClientService(),
		NewTokenService(),
	}
}

type loginService struct {
	ClientService
	TokenService
}