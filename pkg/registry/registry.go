package registry

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
)

// Registry 当前为 no-op 实现。
// 对外行为与真实注册中心实现完全一致，业务层无需感知。
type Registry struct{}

func New() *Registry {
	return &Registry{}
}

// Endpoint 返回服务的拨号地址。
// 调用方始终使用 reg.Endpoint("service-name:port")，不感知底层格式。
//
// 当前（no-op）：返回服务名:端口，K8s 集群内 CoreDNS 自动解析。
// 切换注册中心后：返回 discovery:///service-name:port，grpc.WithDiscovery 接管解析。
// 两种情况下调用方代码不需要任何改动。
func (r *Registry) Endpoint(endpoint string) string {
	return endpoint
}

func (r *Registry) Register(_ context.Context, _ *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) Deregister(_ context.Context, _ *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) GetService(_ context.Context, _ string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (r *Registry) Watch(_ context.Context, _ string) (registry.Watcher, error) {
	return &noopWatcher{}, nil
}

type noopWatcher struct{}

func (w *noopWatcher) Next() ([]*registry.ServiceInstance, error) { return nil, nil }
func (w *noopWatcher) Stop() error                                { return nil }
