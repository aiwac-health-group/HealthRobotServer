package services

type ServiceService interface {
	ClientService
	TreatInfoService
}

func NewServiceService() ServiceService {
	return &serviceService{
		NewClientService(),
		NewTreatInfoService(),
	}
}

type serviceService struct {
	ClientService
	TreatInfoService
}