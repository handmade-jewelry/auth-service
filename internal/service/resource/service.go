package resource

import (
	"context"
)

type Service struct {
	//todo
	//repo
}

func NewService() *Service {
	return &Service{}
}

func (r *Service) GetRouteByPath(ctx context.Context, path string) (*Resource, error) {
	//todo stub
	route := Resource{
		ID:        1,
		ServiceID: 2,
		Path:      "/items",
		Roles:     []string{"CUSTOMER_ROLE"},
		IsActive:  true,
	}

	return &route, nil
}
