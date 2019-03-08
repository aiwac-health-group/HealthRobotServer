package services

type DoctorService interface {
	ClientService
	HealthReportService
}

func NewDoctorService() DoctorService {
	return &doctorService{
		NewClientService(),
		NewHealthReportService(),
	}
}

type doctorService struct {
	ClientService
	HealthReportService
}