## 业务层接入姿势
no-op 阶段和真实实现阶段，业务层写法完全相同，不需要任何改动。

服务发现示例
```go
// internal/data/data.go
func NewOrderServiceClient(r registry.Registry) orderv1.OrderClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(r.Endpoint("order-svc:9000")),  // 调用方传服务名:端口，不关心底层格式
    )
	if err != nil {
		panic(err)
	}
	return orderv1.NewOrderClient(conn)
}
```


reg.Endpoint() 内部控制返回什么格式：

| 阶段 | `Endpoint("order-svc:9000")` 返回 | `grpc.WithDiscovery` 行为 |
|------|-----------------------------------|---------------------------|
| 当前 no-op | `order-svc:9000` | 直接连，CoreDNS 自动解析 |
| 切换配置中心比如 Nacos后 | `discovery:///order-svc:9000` | 接管解析，查注册中心 |
