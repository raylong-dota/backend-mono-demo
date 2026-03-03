# backend-mono 

Kratos-based Monorepo Skeleton for Microservices
单仓多服务（Monorepo）微服务基础骨架与标准模板


📌 项目目标

backup-mono 是一个基于 Go + Kratos 的单仓多服务基础结构，提供：

✅ 统一的微服务目录结构

✅ helloworld 作为标准模板（Golden Template）

✅ 统一的 proto 生成策略

✅ 统一的 Makefile 构建入口

✅ 自动化 make new SERVICE=xxx 服务生成能力

✅ 便于未来接入 K8s / Gateway / Config Center / Event Bus

本项目用于规范公司内部 10–20 个微服务的统一开发模式。

🏗️ 目录结构

```
backup-mono/
├── api/                        # 契约层（Proto 定义）
│   ├── helloworld/        # 模板服务 proto
│   └── ...
│
├── app/                        # 应用实现层
│   ├── helloworld/        # 标准模板服务
│   └── <new-service>/          # 业务服务
│
├── pkg/                        # 公共库（middleware / errors / utils 等）
├── scripts/                    # 自动化脚本（new-service 等）
├── deployments/                # k8s / helm 等部署文件
├── Makefile                    # 统一构建入口
├── go.work                     # 多 module workspace
└── README.md
```

🧱 架构约束

1️⃣ 分层原则（强制）

每个服务必须遵循：

```
internal/
  ├── service   # transport 层（实现 proto）
  ├── biz       # 业务逻辑层（UseCase）
  ├── data      # 数据访问层（repo）
  └── server    # http/grpc server 注册
```
禁止：
	•	❌ service 直接访问 data
	•	❌ 跨服务 import internal
	•	❌ proto 生成物散落各目录


2️⃣ Proto 规范
	•	proto 统一放在：

```
api/<service>/service/v1/*.proto
```

🚀 快速开始

1️⃣ 初始化仓库
```
git clone <repo-url>
cd backup-mono
```
默认会读取
```
make proto SERVICE=helloworld
make run SERVICE=helloworld
```

3️⃣ 创建新服务
```
make new SERVICE=user-svc
```
生成
```
api/user-svc/
app/user-svc/
```
然后
```
make proto SERVICE=user-svc
make run SERVICE=user-svc
```

🛠️ Makefile 统一命令

命令
作用
make new SERVICE=xxx
创建新服务
make proto SERVICE=xxx
生成 proto
make build SERVICE=xxx
构建服务
make run SERVICE=xxx
本地运行
make test
运行全部单测

服务创建流程
1.	执行：
```
make new SERVICE=order-svc
```
2. 修改
```
api/order-svc/service/v1/*.proto
```
3. 生成代码
```
make proto SERVICE=order-svc
```
4. 实现业务逻辑
	•	internal/service
	•	internal/biz
	•	internal/data
5.	本地启动：
```
make run SERVICE=order-svc
```

