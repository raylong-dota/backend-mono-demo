package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	helloworldv1 "github.com/ray-dota/backend-mono/api/helloworld/service/v1"
)

// HelloworldUsecase calls the helloworld service via gRPC.
type HelloworldUsecase struct {
	client helloworldv1.GreeterClient
	log    *log.Helper
}

// NewHelloworldUsecase creates a new HelloworldUsecase.
func NewHelloworldUsecase(client helloworldv1.GreeterClient, logger log.Logger) *HelloworldUsecase {
	return &HelloworldUsecase{
		client: client,
		log:    log.NewHelper(log.With(logger, "module", "biz/helloworld")),
	}
}

// SayHello calls the helloworld service and returns the message.
func (uc *HelloworldUsecase) SayHello(ctx context.Context, name string) (string, error) {
	resp, err := uc.client.SayHello(ctx, &helloworldv1.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return resp.GetMessage(), nil
}
