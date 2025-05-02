package user

import (
	"context"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/handmade-jewelry/auth-service/internal/config"
	userService "github.com/handmade-jewelry/user-service/pkg/api/user-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (u *Service) RefreshToken(ctx context.Context) {
	//todo stub
}

func (u *Service) Login(ctx context.Context) {
	//todo stub
}

func (u *Service) ParseToken(ctx context.Context) {
	//todo stub
}

func (u *Service) GenerateToken(ctx context.Context) {
	//todo stub
}
