package service_service

import "context"

type ServiceService struct {
}

func NewService() *ServiceService {
	return &ServiceService{}
}

func (s *ServiceService) GetService(ctx context.Context, id int) (*Service, error) {
	return &Service{}, nil
}

func (s *ServiceService) CreateService() {

}

func (s *ServiceService) UpdateService() {

}

func (s *ServiceService) DeleteService() {

}
