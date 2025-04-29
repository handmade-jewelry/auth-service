package resource_service

import (
	"context"
)

type ResourceService struct {
	//todo
	//repo
}

func NewService() *ResourceService {
	return &ResourceService{}
}

func (r *ResourceService) GetRouteByPath(ctx context.Context, path string) (*Resource, error) {
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
