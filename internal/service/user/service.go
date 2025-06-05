package user

import (
	"context"
	"github.com/handmade-jewelry/auth-service/logger"
	"google.golang.org/protobuf/types/known/emptypb"

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
	req := &userService.GetUserRolesRequest{
		UserId: userID,
	}

	res, err := s.client.GetUserRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(res.GetRoles()))
	for _, role := range res.GetRoles() {
		roles = append(roles, role.GetName())
	}

	return roles, nil

	//todo stub
	//return []string{"CUSTOMER"}, nil
}

func (s *Service) RoleMap(ctx context.Context) (map[string]struct{}, error) {
	res, err := s.client.ListRoles(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	roles := make(map[string]struct{}, len(res.GetRoles()))
	for _, role := range res.GetRoles() {
		name := role.GetName()
		roles[name] = struct{}{}
	}

	return roles, nil

	//todo stub
	//return map[string]struct{}{
	//	"CUSTOMER": {},
	//	"ADMIN":    {},
	//	"SELLER":   {},
	//}, nil
}

func (s *Service) RoleList(ctx context.Context) ([]string, error) {
	res, err := s.client.ListRoles(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(res.GetRoles()))
	for _, role := range res.GetRoles() {
		roles = append(roles, role.GetName())
	}

	return roles, nil

	//todo stub
	//return []string{"CUSTOMER", "ADMIN", "SELLER"}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*UserWithRoles, error) {
	res, err := s.client.Login(ctx, &userService.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(res.GetRoles()))
	for _, role := range res.GetRoles() {
		roles = append(roles, role.GetName())
	}

	return &UserWithRoles{
		UserID: res.GetUserId(),
		Roles:  roles,
	}, nil

	//todo stub
	//return &UserWithRoles{
	//	UserID: 1,
	//	Roles:  []string{"CUSTOMER"},
	//}, nil
}
