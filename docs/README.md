# HU MQ - 企业级消息队列平台

HU MQ 是基于 Apache Kafka 构建的企业级消息队列平台，提供中文监控后台、运维工具、权限管理和完整的中文使用指引。

## 核心功能

- **集群概览** — 多集群管理，Broker 状态、TPS 实时监控
- **Topic 管理** — 可视化创建/删除/配置 Topic，分区管理
- **消息追踪** — 按 Key/Offset/时间查询消息，支持死信队列和消息重放
- **消费组监控** — 消费进度、Lag 排行、成员状态实时展示
- **性能监控** — TPS、延迟、磁盘/CPU 指标实时图表
- **告警管理** — 积压告警、节点宕机告警，支持多通知渠道
- **权限管理** — 用户管理、角色分配、Topic/消费组级别 ACL
- **运维工具** — 分区重分配、配置热更新
- **使用指引** — 新手快速入门、SDK 示例、最佳实践

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go + Gin + GORM |
| 前端 | Vue 3 + Element Plus + ECharts |
| 引擎 | Apache Kafka (兼容 Redpanda) |
| 数据库 | PostgreSQL 15 |
| 部署 | Docker Compose / Kubernetes Helm |

## 快速开始

参考 [快速入门指南](./quick-start.md)

## 文档

- [快速入门](./quick-start.md)
- [部署指南](./deployment.md)
- [操作手册](./user-guide.md)
- [SDK 指南](./sdk-guide.md)
- [最佳实践](./best-practices.md)
- [常见问题](./faq.md)