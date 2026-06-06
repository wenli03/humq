# HU MQ 企业级消息队列平台 — 产品设计文档

> 版本：v1.0 | 日期：2026-06-06 | 状态：已确认

## 一、产品概述

**HU MQ** 是基于 Apache Kafka 构建的企业级消息队列平台。在 Kafka 引擎之上封装了监控后台、运维工具、权限管理和中文使用指引，形成开箱即用的企业级 MQ 产品。不修改 Kafka 源码，通过 Sarama Admin Client 与 Kafka 集群交互。

### 目标用户
- 运维工程师：集群管理、监控告警
- 开发工程师：消息追踪、SDK 接入、Topic 管理
- 技术管理者：全局视角、健康度评估

---

## 二、技术架构

### 2.1 架构图

```
┌─────────────────────────────────────────────────────────┐
│                    HU MQ 监控后台                        │
│  ┌──────────┐   HTTP/WS    ┌──────────────────────────┐ │
│  │ Vue3 前端 │ ◄──────────► │      Go API Server       │ │
│  │ (中文UI)  │              │  Gin + JWT + Swagger     │ │
│  └──────────┘              │                          │ │
│                            │  ┌──────────┬──────────┐ │ │
│                            │  │ Service  │Collector │ │ │
│                            │  │  Layer   │(定时采集) │ │ │
│                            │  └────┬─────┴────┬─────┘ │ │
│                            └───────┼──────────┼───────┘ │
│                                    │          │         │
│                    ┌───────────────┼──────────┼───┐     │
│                    │               ▼          ▼   │     │
│                    │  ┌─────────┐  ┌──────────┐  │     │
│                    │  │PostgreSQL│  │  Kafka   │  │     │
│                    │  │(元数据)  │  │ Cluster  │  │     │
│                    │  └─────────┘  └──────────┘  │     │
│                    │          数据层              │     │
│                    └─────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
```

### 2.2 架构原则

- **不侵入 Kafka**：通过 Sarama Admin Client 与 Kafka 交互，不修改 Kafka 源码
- **模块化单体**：单一 Go 二进制，内嵌 Vue 前端静态文件，内部按模块分层
- **数据分离**：PostgreSQL 存储 HU MQ 自身数据（用户、告警、监控历史）；Kafka 原数据不重复存储
- **单容器部署**：Docker 镜像一个进程启动，开发/测试/生产一致

---

## 三、技术选型

| 层 | 选型 | 版本 | 理由 |
|---|---|---|---|
| HTTP 框架 | Gin | v1.9+ | 高性能，Go 生态最流行 |
| ORM | GORM | v2 | 最成熟的 Go ORM |
| Kafka 客户端 | Sarama | v1.43+ | Go 原生 Kafka Admin API |
| 认证 | JWT | - | 无状态，前后端分离友好 |
| API 文档 | Swagger | - | 自动生成，前后端对齐 |
| 前端框架 | Vue 3 + Vite | 3.4+ | 现代化，生态完善 |
| UI 库 | Element Plus | 2.8+ | 企业级中文组件库 |
| 图表 | ECharts | 5.5+ | 中文文档好，功能最强 |
| 数据库 | PostgreSQL | 15 | JSON 支持，性能优 |
| 容器编排 | Docker Compose + K8s Helm | - | 覆盖开发/生产 |

---

## 四、项目目录结构

```
humq/
├── cmd/humq-server/main.go          # 启动入口
├── internal/
│   ├── api/                         # HTTP handler + 路由 + 中间件
│   ├── service/                     # 业务逻辑层
│   ├── repository/                  # GORM 数据访问层
│   ├── kafka/                       # Sarama Admin Client 封装
│   ├── collector/                   # 指标定时采集（TPS/延迟/Lag）
│   ├── alerter/                     # 告警规则引擎
│   └── auth/                        # JWT + RBAC 权限
├── web/                             # Vue 3 前端源码
│   ├── src/
│   │   ├── views/                   # 页面（集群/Topic/消费组/告警/消息/权限/运维/指引）
│   │   ├── components/              # 复用组件
│   │   ├── api/                     # 后端接口调用封装
│   │   ├── router/                  # Vue Router
│   │   └── stores/                  # Pinia 状态管理
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
├── deploy/
│   ├── docker-compose.yml           # 本地开发一键启动
│   ├── Dockerfile                   # humq-server 镜像构建
│   └── helm/                        # K8s Helm Chart
├── docs/
│   ├── README.md                    # 产品介绍
│   ├── quick-start.md               # 5分钟快速开始
│   ├── deployment.md                # Docker & K8s 部署
│   ├── user-guide.md                # 监控后台操作手册
│   ├── sdk-guide.md                 # SDK 接入示例
│   ├── best-practices.md           # 最佳实践
│   ├── faq.md                       # 常见问题
│   └── changelog.md                 # 版本记录
└── go.mod
```

---

## 五、数据模型

### PostgreSQL 核心表

| 表名 | 说明 | 关键字段 |
|---|---|---|
| `users` | 用户 | id, username, password(bcrypt), role(admin/user), created_at |
| `clusters` | Kafka 集群连接 | id, name, bootstrap_servers, status, created_at |
| `topics_meta` | Topic 元数据缓存 | id, cluster_id, name, partitions, replicas, retention_ms, created_at |
| `consumer_groups` | 消费组 | id, cluster_id, group_id, topic, members, state |
| `metrics_snapshots` | 指标快照 | id, cluster_id, metric_type, value, tags(jsonb), collected_at |
| `alert_rules` | 告警规则 | id, cluster_id, metric, operator(>/=), threshold, channels(jsonb) |
| `alert_events` | 告警事件 | id, rule_id, level, message, status(triggered/resolved), triggered_at |
| `dead_messages` | 死信消息 | id, cluster_id, topic, partition, offset, key, payload, timestamp |
| `acl_rules` | 权限规则 | id, user_id, resource(topic/group), resource_name, operation(read/write/admin) |

**关系：**
- `users` 1:N `acl_rules`
- `clusters` 1:N `topics_meta`, `consumer_groups`, `metrics_snapshots`, `alert_rules`, `dead_messages`
- `alert_rules` 1:N `alert_events`

---

## 六、API 设计

### 统一格式

请求成功：`{ "code": 0, "msg": "success", "data": {...} }`
请求失败：`{ "code": 4xxx, "msg": "错误描述", "data": null }`
分页追加：`{ "total": 100, "page": 1, "page_size": 20 }`

### 认证方式
`Authorization: Bearer <JWT_TOKEN>`

### 接口列表

```
# 认证
POST   /api/v1/auth/login              登录 {username, password} → {token, refresh_token}
POST   /api/v1/auth/refresh            刷新Token

# 集群管理
GET    /api/v1/clusters                列表
POST   /api/v1/clusters                注册 {name, bootstrap_servers}
GET    /api/v1/clusters/:id/metrics    实时指标 {brokers, tps, lag_total}

# Topic 管理
GET    /api/v1/topics                  列表（支持：keyword, cluster_id）
POST   /api/v1/topics                  创建 {name, partitions, replication_factor, configs}
DELETE /api/v1/topics/:name            删除
PUT    /api/v1/topics/:name/config    修改配置
GET    /api/v1/topics/:name/partitions 分区详情

# 消费组
GET    /api/v1/consumers               列表
GET    /api/v1/consumers/:name/lag     积压详情 {per_partition_lag}

# 消息追踪
POST   /api/v1/messages/trace          查询 {keyword?, timestamp_start?, offset?}
GET    /api/v1/messages/dead           死信列表
POST   /api/v1/messages/replay         重放 {topic, from_offset, to_offset, target_topic}

# 告警管理
GET    /api/v1/alerts/rules            规则列表
POST   /api/v1/alerts/rules            创建 {metric, operator, threshold, channels}
PUT    /api/v1/alerts/rules/:id        更新
DELETE /api/v1/alerts/rules/:id        删除
GET    /api/v1/alerts/events           事件历史（分页、筛选）

# 用户与权限
GET    /api/v1/users                   用户列表
POST   /api/v1/users                   创建
PUT    /api/v1/users/:id               更新
DELETE /api/v1/users/:id               删除
POST   /api/v1/acls                    配置权限
GET    /api/v1/acls                    权限列表

# 运维操作
GET    /api/v1/ops/rebalance           分区重分配状态
POST   /api/v1/ops/rebalance           触发重分配 {topics_json}

# 实时推送
WS     /api/v1/ws/dashboard            仪表盘 WebSocket 推送
```

---

## 七、前端页面结构

```
HU MQ 监控后台（Element Plus + 深色侧边栏 + 蓝色主题）

/login            登录页 — 用户名密码登录
/dashboard        仪表盘 — 集群健康度卡片、TPS 实时曲线、节点状态、告警概览
/clusters         集群管理 — Broker 列表、节点状态、关键指标
/clusters/:id     集群详情 — 深度指标、Topic/消费组概览
/topics           Topic 管理 — 列表、搜索、新建/删除/配置
/topics/:name     Topic 详情 — 分区列表、副本分布、消息采样
/consumers        消费组监控 — Lag 排行、成员状态、消费速率
/consumers/:name  消费组详情 — 分区级 Lag、消费曲线
/messages         消息追踪 — 按 Key/Offset/时间范围查询
/messages/dead    死信队列 — 查看、重放操作
/alerts           告警中心 — 规则管理、事件列表、状态筛选
/alerts/rules     告警规则 — 创建/编辑/启停
/acl              权限管理 — 用户CRUD、角色分配、资源权限
/ops              运维工具 — 分区重分配、扩缩容引导
/guide            使用指引 — 内嵌文档浏览
```

---

## 八、功能模块详情

### 8.1 集群概览
- 仪表盘首页：总 Broker 数、在线/离线、总 Topic、总 Lag
- TPS 实时曲线（通过 Sarama 采集 Broker JMX 指标）
- 节点列表：每台 Broker 的 CPU、内存、磁盘、网络
- WebSocket 推送实现实时刷新

### 8.2 Topic 管理
- 列表展示：名称、分区数、副本、保留时间、消息数
- 创建向导：名称、分区数、副本因子、保留策略（删除/压缩）
- 配置修改：可修改 retention.ms、segment.bytes 等
- 分区详情：每个分区的 Leader、Replicas、ISR、Log End Offset

### 8.3 消息追踪
- 按 Key/Timestamp/Offset 组合条件查询
- 支持 deserialize 预览（string/json/avro 基础支持）
- 死信队列：查看滞留消息、批量/单选重放到指定 Topic

### 8.4 消费组监控
- 消费组列表：组名、订阅 Topic、消费模式、成员数
- Lag 排行：按积压量排序，红色高亮警告
- 分区级 Lag：每个分区的 current offset、log end offset、lag

### 8.5 性能监控
- 定时采集（可配置间隔，默认 30s）
- 采集指标：TPS(生产/消费)、延迟 P99、磁盘使用率、网络 IO
- 历史数据存储到 PostgreSQL metrics_snapshots 表
- 前端 ECharts 时序图表展示

### 8.6 告警管理
- 规则引擎：metric + operator + threshold 组合
- 支持指标：lag、tps_drop、broker_down、disk_usage
- 通知渠道：Webhook、邮件（MVP），预留钉钉/飞书
- 告警生命周期：triggered → acknowledged → resolved

### 8.7 权限管理
- 用户管理：增删改查，密码 bcrypt 加密
- 角色：admin（全局）、user（受限）
- ACL：用户 → 资源(Topic/消费组) → 操作(read/write/admin)
- JWT Token 过期 + Refresh Token 续期

### 8.8 消息重放与死信
- 死信消息自动记录（消费失败超过重试次数）
- 手动重放：选择死信消息 → 指定目标 Topic → 批量/单条重放
- 时间回溯消费：指定时间范围，将历史消息重放到新 Topic

### 8.9 集群运维
- 分区重分配：选择 Topic → 生成/执行 reassignment plan
- 配置热更新：Broker/Topic 级别动态配置修改
- 操作审计日志

### 8.10 使用指引
- 快速入门：Docker Compose 3 步启动
- 部署指南：Docker 详细配置、K8s Helm 安装
- 操作手册：后台每个功能的使用说明（截图+文字）
- SDK 指南：Go/Java/Python SDK 代码示例
- 最佳实践：Topic 命名规范、分区规划、告警阈值建议
- FAQ：常见问题排查

---

## 九、部署方案

### Docker Compose（开发/测试）

```yaml
services:
  humq-server:
    build: .
    ports: ["8080:8080"]
    environment:
      - DB_DSN=postgres://humq:password@postgres:5432/humq
      - KAFKA_BOOTSTRAP=redpanda:9092
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: humq
      POSTGRES_USER: humq
      POSTGRES_PASSWORD: password
  redpanda:
    image: redpandadata/redpanda:latest
    command: redpanda start --smp 1
```

> 开发阶段使用 Redpanda 替代 Kafka（API 兼容、资源更省）

### K8s（生产）
- Helm Chart：humq-server Deployment + PostgreSQL StatefulSet
- Kafka 集群独立部署（不纳入 Helm Chart）
- ConfigMap 管理环境配置
- Ingress 暴露前端

---

## 十、验收标准

| 功能 | 验收条件 |
|---|---|
| 集群概览 | 仪表盘正确显示 Broker 状态、TPS 曲线实时更新 |
| Topic 管理 | 可创建/删除/配置 Topic，分区列表与 Kafka 一致 |
| 消息追踪 | 可按条件查到消息内容，死信可重放 |
| 消费组监控 | Lag 数值正确，支持到分区级别 |
| 性能监控 | 指标曲线与 Kafka 实际指标一致，历史可查询 |
| 告警管理 | 规则创建后，触发条件满足时产生告警事件 |
| 权限管理 | 非 admin 用户不能执行越权操作 |
| 运维工具 | 分区重分配可提交并观察进度 |
| 使用指引 | 文档可独立阅读，新人按文档能完成部署和基础操作 |
| 部署 | Docker Compose 一键启动，所有服务正常 |
