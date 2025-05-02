package service

import "context"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetService(ctx context.Context, id int) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) CreateService() {

}

func (s *Service) UpdateService() {

}

func (s *Service) DeleteService() {

}
