# GoMall 微服务电商平台

## 项目简介

GoMall是一个基于微服务架构的现代电商平台，采用Go语言开发，集成了多种先进的技术组件和功能特性。
- [API文档](API.md)

### 核心特性

- **微服务架构**：采用分布式微服务设计，服务间通过RPC通信
- **服务治理**：
  - ETCD服务注册与发现
  - Consul服务治理
  - Jaeger链路追踪
- **可观测性**：
  - Grafana + Promtail + Loki 日志收集与可视化
  - Prometheus 指标监控与告警
- **消息队列**：集成NATS实现异步通信
- **容器化部署**：支持本地部署、Docker、Docker Compose和Kubernetes多种部署方式
- **AI赋能**：集成Doubao Function Call AI大模型实现智能查询和预下单

## 项目架构

### 服务模块

- **前端服务**：Vue.js构建的整个服务前端界面
- **核心服务**：
  - 用户服务(User Service)
    - 用户注册登入登出，用户角色管理
  - 商品服务(Product Service)
    - 商品的增删改查
  - 购物车服务(Cart Service)
    - 查看/添加商品到购物车
  - 订单服务(Order Service)
    - 订单取消/定时取消，订单二次支付
  - 支付服务(Payment Service)
    - 支付订单/取消支付
  - 结算服务(Checkout Service)
  - 邮件服务(Email Service)
    - 发送邮件，针对支付和预下单
  - 前后端中间层服务(Frontend Service)
    - 用户认证和角色鉴权: jwt + rbac，令牌续费
  - AI服务(AIEino Service)
    - Ai查询订单，Ai预下单

## 快速开始

### 环境要求

- Linux操作系统(推荐Ubuntu 18.04)
- Go语言环境
- Docker
- Kubernetes
- Kind工具
- Kubectl工具

### 安装步骤

1. **克隆项目**
```bash
git clone <project-url>
cd gomall
```

2. **安装依赖**
- 安装后端服务依赖
```bash
cd app/<service-name>
go mod tidy
```

- 安装前端依赖
```bash
cd frontend-vue
npm install
```

3. **配置**
- 按需检查`exec.sh`文件中的服务端口配置和一些开发时使用过的命令
- 修改`components/prometheus/config/prometheus.yaml`中的黑盒监控地址为实际的Ubuntu主机地址
- 项目的app目录下各个服务的配置采用的是k8s的configmap，可根据需要去到deploy/k8s/gomall-dev-app.yaml修改

### 部署流程

1. **创建Kubernetes集群**
项目根目录下:
```bash
kind create cluster --config=./deploy/k8s/gomall-dev-cluster.yaml
```

2. **准备基础镜像**

确保以下镜像已在系统中准备就绪：
- mysql:5.7
- etcd:v3.5.5
- redis:7.2.5
- jaegertracing/all-in-one:1.62.0
- nats:2.10.19
- prometheus:v2.34.0
- consul:1.8.8
- grafana/promtail:2.9.2
- grafana/loki:2.9.2
- grafana:8.3.3
- alpine:3.14
- nginx:stable-alpine
- node:18-alpine

3. **构建服务镜像**
可去到deploy/mk目录下，修改仓库地址为自己的仓库地址，然后回到根目录执行以下命令：
```bash
# 构建应用服务镜像
make release-dev

# 构建监控组件镜像
make release-prometheus-exports

# 构建前端服务镜像
cd frontend-vue
make frontend-vue-dev
```

4. **加载镜像到集群**
```bash
kind load docker-image --name gomall-dev mysql:5.7 etcd:v3.5.5 redis:7.2.5 jaegertracing/all-in-one:1.62.0 nats:2.10.19 prometheus:v2.34.0 consul:1.8.8 grafana/promtail:2.9.2 grafana/loki:2.9.2 grafana:8.3.3 go-mall-frontend-dev:latest go-mall-product-dev:latest go-mall-email-dev:latest go-mall-checkout-dev:latest go-mall-order-dev:latest go-mall-cart-dev:latest go-mall-payment-dev:latest go-mall-user-dev:latest go-mall-aieino-dev:latest blackbox:latest alertmanager:latest go-mall-frontend-vue-dev:latest
```

5. **部署服务**
```bash
# 部署基础设施和组件
kubectl apply -f deploy/k8s/gomall-dev-base.yaml

# 部署应用服务
kubectl apply -f deploy/k8s/gomall-dev-app.yaml
```

> 注意：如果基础设施部署失败，可能需要增加kind节点容器中`gomall/components`目录的权限

## 开发指南

### 目录结构

```
├── app/                # 微服务应用目录
├── common/            # 公共代码
├── components/        # 基础组件配置
├── deploy/            # 部署相关配置
├── frontend-vue/     # 前端应用
├── idl/              # 接口定义文件
└── rpc_gen/          # RPC生成代码
```

### 开发流程

1. 在`idl`目录下定义服务接口
2. 生成RPC代码
3. 实现服务业务逻辑
4. 在`idl\frontend`目录下按需定义中间层frontend接口
5. 生成HTTP代码
6. 实现对RPC的调用逻辑和结果对前端的返回
7. 前端方面相关调用HTTP路由接口
8. 本地测试
9. 构建镜像并部署

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 发起Pull Request
