package user

import (
	"context"

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
		//todo  log
		return nil, err
	}

	return &Service{
		client: userService.NewUserServiceClient(conn),
	}, nil
}

func (u *Service) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	//todo stub
	return []string{"CUSTOMER"}, nil
}

func (u *Service) Login(ctx context.Context, email, password string) (*UserWithRoles, error) {
	//todo stub
	return &UserWithRoles{
		UserID: 1,
		Roles:  []string{"CUSTOMER"},
	}, nil
}

func (u *Service) Logout(ctx context.Context) {
	//todo stub
}
