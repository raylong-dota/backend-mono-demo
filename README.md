# backend-mono

基于 Go + Kratos v2 的单仓多服务（Monorepo）骨架，面向 CFD 交易平台后端。

## 目录结构

```
backend-mono/
├── api/                            # 契约层（Proto 定义 & 生成代码）
│   └── <app>/<service>/v1/        # *.proto / *.pb.go / *.swagger.json
│
├── app/                            # 应用实现层
│   └── <app>/<service>/
│       ├── cmd/server/             # 程序入口 (main.go, wire.go)
│       ├── configs/                # 服务配置文件
│       ├── internal/
│       │   ├── biz/               # 业务逻辑层（UseCase + Repo 接口）
│       │   ├── data/              # 数据访问层（Repo 实现）
│       │   ├── service/           # Transport 层（实现 proto Server 接口）
│       │   ├── server/            # HTTP / gRPC Server 注册
│       │   └── conf/              # 配置结构体（由 proto 生成）
│       └── Makefile               # 服务级构建命令（继承自 app_makefile）
│
├── bin/                            # 编译产物（git-ignored）
│   └── orbit-<svc>-svc            # 二进制命名规范
│
├── .github/
│   └── workflows/
│       └── pr-gate.yml            # PR Gate：向 main 发 PR 时自动运行 lint
│
├── configs/
│   └── golangci.yaml              # 全局 lint 规则
│
├── scripts/
│   ├── install_base.sh            # 一键安装项目本地工具链
│   └── new.sh                     # 新服务脚手架
│
├── deploy/
│   ├── build/
│   │   └── Dockerfile             # 多阶段构建，接收 APP_SVC build arg
│   └── k8s/
│       └── <svc>/                 # 每个服务一个目录
│           ├── configmap.yaml
│           ├── deployment.yaml
│           └── service.yaml
│
├── third_party/                    # proto 依赖（google/api、validate 等）
├── app_makefile                   # 服务级 Makefile 模板
├── Makefile                       # 根级构建入口
└── go.mod
```

本地工具链安装到项目目录，不依赖全局环境，已加入 `.gitignore`：

```
.go/bin/       → Go 工具二进制（protoc-gen-go, wire, golangci-lint …）
.go/pkg/mod/   → Go 模块源码缓存（IDE 可直接跳转）
.tools/        → 非 Go 工具（protoc）
```

## 快速开始

### 1. 克隆并安装工具链

```bash
git clone <repo-url>
cd backend-mono

make install   # 安装 Go 1.26.0、protoc、wire、golangci-lint 等（幂等，已装则跳过）
```

**推荐安装 [direnv](https://direnv.net/)，然后执行：**

```bash
direnv allow
```

之后 `cd` 进项目目录，终端里的 `go`、`wire` 等命令自动使用项目本地版本（`.tools/go/bin`），无需手动切换。

不使用 direnv 时，用 make 命令代替直接的 `go` 命令：

| 直接命令 | 等价的 make 命令 |
|---|---|
| `go mod tidy` | `make tidy` |
| `go get github.com/foo/bar@v1.2.3` | `make get pkg=github.com/foo/bar@v1.2.3` |

### 2. 运行模板服务（helloworld）

```bash
make run helloworld
```

服务默认监听：
- HTTP `0.0.0.0:8000`
- gRPC `0.0.0.0:9000`

### 3. 创建新服务

```bash
make new order   # 生成 app/order/service 和 api/order/service
```

生成后的工作流：

```bash
# 1. 修改 proto 定义
vim api/order/service/v1/order.proto

# 2. 生成 gRPC / HTTP / Swagger + Wire 代码
make generate order

# 3. 实现业务逻辑
#    internal/biz/     — 定义 UseCase 和 Repo 接口
#    internal/data/    — 实现 Repo（DB / cache / 外部 API）
#    internal/service/ — 实现 proto Server 接口，调用 biz

# 4. 本地运行
make run order

# 5. 编译产物
make build order   # → bin/orbit-order-svc
```

## Makefile 命令

### 根目录

| 命令 | 说明 |
|---|---|
| `make install` | 安装所有项目本地工具链（幂等，已装则跳过） |
| `make new <svc>` | 从 helloworld 模板创建新服务 |
| `make generate <svc>\|all` | 对指定服务或所有服务生成代码（api + wire + proto） |
| `make build <svc>\|all` | 编译指定服务或所有服务，产物为 `bin/orbit-<svc>-svc` |
| `make run <svc>` | 本地运行指定服务（不支持 `all`） |
| `make image <svc>` | 构建指定服务的 Docker 镜像，默认 tag `latest` |
| `make image <svc> tag=<tag>` | 构建并指定镜像 tag |
| `make clean` | 删除所有编译产物（根目录及各服务的 `bin/`） |
| `make tidy` | 运行 `go mod tidy`（使用项目本地 Go） |
| `make get pkg=<module@version>` | 添加或升级依赖（使用项目本地 Go） |
| `make lint` | 全项目运行 golangci-lint 并自动修复（本地开发） |
| `make lint-check` | 全项目运行 golangci-lint，仅检查不修复（CI 使用） |

> 服务名为位置参数，直接跟在命令后面，无需 `svc=` 前缀。

### 服务目录（`app/<app>/<service>/`）

| 命令 | 说明 |
|---|---|
| `make api` | 生成 gRPC + HTTP + Swagger + errors 代码（等同于下面四项合集） |
| `make grpc` | 生成 gRPC 代码 |
| `make http` | 生成 HTTP 代码 |
| `make errors` | 生成错误码代码 |
| `make swagger` | 生成 Swagger 文档 |
| `make proto` | 生成内部 proto struct |
| `make wire` | 生成 Wire 依赖注入代码 |
| `make generate` | 运行 `go generate ./...` |
| `make build` | 编译二进制（默认产物在服务本地 `bin/<APP_NAME>`） |
| `make run` | 本地运行服务 |
| `make test` | 运行单元测试 |
| `make ent` | 生成 ent ORM 代码 |
| `make docker` | 构建 Docker 镜像 |

## CI / PR Gate

向 `main` 分支发起 PR 时自动触发 `.github/workflows/pr-gate.yml`：

1. `make install` — 安装工具链
2. `make lint-check` — 以非 0 退出码阻断不符合规范的 PR

本地修复 lint 问题：

```bash
make lint   # 自动修复后重新提交
```

## 架构约束

**分层调用方向（严格单向）**

```
service → biz → data
```

- `service` 只能调用 `biz`，不能直接访问 `data`
- `biz` 定义 Repo 接口，`data` 负责实现
- `biz` 可以 import `api` 层的错误码（Kratos 惯例）
- 禁止跨服务 import `internal/`

**Proto 路径规范**

```
api/<app>/<service>/v1/*.proto
```

**import 顺序（由 gci 自动检查）**

```go
import (
    // 1. 标准库
    "context"

    // 2. 第三方
    "github.com/go-kratos/kratos/v2/log"

    // 3. 本项目
    "github.com/ray-dota/backend-mono/..."
)
```

## 开发工具版本

工具版本通过 `scripts/install_base.sh` 顶部变量管理：

| 工具 | 版本管理 |
|---|---|
| `go` | `GO_VERSION` 变量（当前 `1.26.0`），安装到 `.tools/go/` |
| `protoc` | `PROTOC_VERSION` 变量（当前 `33.4`），安装到 `.tools/` |
| `protoc-gen-go` / `wire` / `golangci-lint` 等 | `@latest`，安装到 `.go/bin/` |

升级工具版本只需修改 `scripts/install_base.sh` 顶部的版本变量，再重新运行 `make install`。
