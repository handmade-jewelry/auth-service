package user

import (
	"context"
	"github.com/handmade-jewelry/auth-service/logger"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/handmade-jewelry/auth-service/internal/config"
	userService "github.com/handmade-jewelry/user-service/pkg/api/user-service"
)

type Service struct {
	client userService.UserServiceClient
}

func NewService(opts *config.GRPCOptions) (*Service, error) {
	conn, err := grpc.NewClient(opts.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcRetry.UnaryClientInterceptor(
				grpcRetry.WithMax(opts.MaxRetry),
				grpcRetry.WithPerRetryTimeout(opts.PerRetryTimeout),
			),
		))
	if err != nil {
		logger.Error("failed to create grpc client", err)
		return nil, err
	}

	return &Service{
		client: userService.NewUserServiceClient(conn),
	}, nil
}

func (s *Service) UserRoles(ctx context.Context, userID int64) ([]string, error) {
	//todo stub
	return []string{"CUSTOMER"}, nil
}

func (s *Service) RoleMap(ctx context.Context) (map[string]string, error) {
	//todo stub
	return map[string]string{
		"CUSTOMER": "CUSTOMER",
		"ADMIN":    "ADMIN",
		"SELLER":   "SELLER",
	}, nil
}

func (s *Service) RoleList(ctx context.Context) ([]string, error) {
	//todo stub
	//todo get from cashe
	return []string{"CUSTOMER", "ADMIN", "SELLER"}, nil
	//todo set from cashe
}

func (s *Service) CheckRoles(ctx context.Context, roles []string) ([]string, error) {
	//todo stub
	return roles, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*UserWithRoles, error) {
	//todo stub
	return &UserWithRoles{
		UserID: 1,
		Roles:  []string{"CUSTOMER"},
	}, nil
}

func (s *Service) Logout(ctx context.Context) {
	//todo stub
}
