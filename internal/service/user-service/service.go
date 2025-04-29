package user_service

import "context"

type UserService struct {
	//todo client
}

func NewService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserRoles(ctx context.Context, token string) (*UserRoles, error) {
	//todo stub
	return &UserRoles{
		UserID: 1,
		Roles:  []string{},
	}, nil
}
