package services

type ServiceService interface {
	ClientService
	TreatInfoService
	RegisterService
	ExamineService
	LectureService
}

func NewServiceService() ServiceService {
	return &serviceService{
		NewClientService(),
		NewTreatInfoService(),
		NewRegisterService(),
		NewExamineService(),
		NewLectureService(),
	}
}

type serviceService struct {
	ClientService
	TreatInfoService
	RegisterService
	ExamineService
	LectureService
}