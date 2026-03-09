package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	v1 "github.com/ray-dota/backend-mono/api/sayhello/service/v1"
	"github.com/ray-dota/backend-mono/app/sayhello/service/internal/conf"
	"github.com/ray-dota/backend-mono/app/sayhello/service/internal/service"
	"github.com/ray-dota/backend-mono/pkg/health"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, sayhello *service.SayhelloService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.GetHttp().GetNetwork() != "" {
		opts = append(opts, http.Network(c.GetHttp().GetNetwork()))
	}
	if c.GetHttp().GetAddr() != "" {
		opts = append(opts, http.Address(c.GetHttp().GetAddr()))
	}
	if c.GetHttp().GetTimeout() != nil {
		opts = append(opts, http.Timeout(c.GetHttp().GetTimeout().AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterSayhelloHTTPServer(srv, sayhello)
	health.Register(srv)
	return srv
}
