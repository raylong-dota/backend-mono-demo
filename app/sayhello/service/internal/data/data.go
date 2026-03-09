package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	helloworldv1 "github.com/ray-dota/backend-mono/api/helloworld/service/v1"
	"github.com/ray-dota/backend-mono/app/sayhello/service/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSayhelloRepo, NewHelloworldClient)

// Data .
type Data struct{}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

// NewHelloworldClient creates a gRPC client for the helloworld service.
func NewHelloworldClient(c *conf.Data) (helloworldv1.GreeterClient, func(), error) {
	conn, err := grpc.NewClient(
		c.GetHelloworld().GetAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	return helloworldv1.NewGreeterClient(conn), func() { _ = conn.Close() }, nil
}
