package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/ray-dota/backend-mono/api/sayhello/service/v1"
	"github.com/ray-dota/backend-mono/app/sayhello/service/internal/biz"
)

// SayhelloService is a sayhello service.
type SayhelloService struct {
	v1.UnimplementedSayhelloServer

	uc  *biz.SayhelloUsecase
	hw  *biz.HelloworldUsecase
	log *log.Helper
}

// NewSayhelloService new a sayhello service.
func NewSayhelloService(uc *biz.SayhelloUsecase, hw *biz.HelloworldUsecase, logger log.Logger) *SayhelloService {
	return &SayhelloService{
		uc:  uc,
		hw:  hw,
		log: log.NewHelper(log.With(logger, "module", "service/sayhello")),
	}
}

// SayHello implements sayhello.SayhelloServer.
func (s *SayhelloService) SayHello(ctx context.Context, in *v1.SayhelloRequest) (*v1.SayhelloReply, error) {
	g, err := s.uc.CreateSayhello(ctx, &biz.Sayhello{Hello: in.GetName()})
	if err != nil {
		return nil, err
	}
	return &v1.SayhelloReply{Message: "Hello " + g.Hello}, nil
}

// HelloProxy proxies the request to the helloworld service.
func (s *SayhelloService) HelloProxy(ctx context.Context, in *v1.HelloProxyRequest) (*v1.HelloProxyReply, error) {
	s.log.Infof("received request: name=%s", in.GetName())
	msg, err := s.hw.SayHello(ctx, in.GetName())
	if err != nil {
		return nil, err
	}
	s.log.Infof("helloworld replied: %s", msg)
	return &v1.HelloProxyReply{Message: msg}, nil
}
