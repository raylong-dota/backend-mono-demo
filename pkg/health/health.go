package health

import (
	"net/http"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Register mounts /healthz and /readyz on the given Kratos HTTP server.
//
// /healthz — liveness:  always returns 200, indicates the process is alive.
// /readyz  — readiness: always returns 200, indicates the process is ready to serve traffic.
//
// Both endpoints can be extended by passing optional check functions that return an error
// when the service is not ready (e.g. database connectivity).
func Register(srv *kratoshttp.Server, checks ...func() error) {
	r := srv.Route("/")
	r.GET("/healthz", func(ctx kratoshttp.Context) error {
		return ctx.Result(http.StatusOK, map[string]string{"status": "ok"})
	})
	r.GET("/readyz", func(ctx kratoshttp.Context) error {
		for _, check := range checks {
			if err := check(); err != nil {
				return ctx.Result(http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "error": err.Error()})
			}
		}
		return ctx.Result(http.StatusOK, map[string]string{"status": "ok"})
	})
}
