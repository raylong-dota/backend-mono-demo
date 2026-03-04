package service

import (
	"context"

	v1 "github.com/ray-dota/backend-mono/api/user/service/v1"
	"github.com/ray-dota/backend-mono/app/user/service/internal/biz"
)

// UserService is a user service.
type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

// NewUserService new a user service.
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// SayHello implements user.UserServer.
func (s *UserService) SayHello(ctx context.Context, in *v1.UserRequest) (*v1.UserReply, error) {
	g, err := s.uc.CreateUser(ctx, &biz.User{Hello: in.GetName()})
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{Message: "Hello " + g.Hello}, nil
}
