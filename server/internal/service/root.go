package service

type Service struct {
	LaptopService *LaptopService
}

func NewService() *Service {
	return &Service{
		NewLaptopService(),
	}
}
