# HU MQ 快速开始

## 环境要求

- Docker 20.10+
- Docker Compose v2
- 4GB 以上可用内存

## 5 分钟部署

### 1. 克隆项目

```bash
git clone https://github.com/wenli03/humq.git
cd humq
```

### 2. 启动服务

```bash
docker compose -f deploy/docker-compose.yml up -d
```

该命令将启动三个服务：
- **humq-server**: HU MQ 监控后台 (端口 8080)
- **postgres**: PostgreSQL 数据库 (端口 5432)
- **redpanda**: Kafka 兼容消息引擎 (端口 9092)

### 3. 访问监控后台

打开浏览器访问：`http://localhost:8080`

默认管理员账号：
- 用户名：`admin`
- 密码：`admin`

### 4. 开始使用

1. 登录后在「集群管理」页面注册 Kafka 集群
2. 在「Topic 管理」创建第一个 Topic
3. 接入您的应用 SDK 开始生产/消费消息
4. 返回仪表盘查看实时监控数据

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 8080 | 服务端口 |
| DB_DSN | host=localhost... | PostgreSQL 连接串 |
| KAFKA_BOOTSTRAP | localhost:9092 | Kafka 地址 |
| JWT_SECRET | humq-secret-... | JWT 密钥 (生产环境务必修改) |

## 停止服务

```bash
docker compose -f deploy/docker-compose.yml down
```

## 下一步

- 查看[部署指南](./deployment.md)了解生产环境部署
- 查看[操作手册](./user-guide.md)了解完整功能
- 查看[SDK 指南](./sdk-guide.md)接入您的应用