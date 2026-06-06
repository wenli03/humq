# HU MQ 部署指南

## Docker 部署

### 开发环境

```bash
docker compose -f deploy/docker-compose.yml up -d
```

### 自定义配置

创建 `.env` 文件覆盖默认配置：

```env
SERVER_PORT=8080
DB_DSN=host=postgres user=humq password=humq dbname=humq port=5432 sslmode=disable
KAFKA_BOOTSTRAP=kafka1:9092,kafka2:9092,kafka3:9092
JWT_SECRET=your-production-secret-key
```

## Kubernetes 部署

### 前提条件

- Kubernetes 1.24+
- Helm 3+
- 已部署 Kafka 集群
- 已部署 PostgreSQL

### Helm 安装

```bash
helm install humq deploy/helm/humq \
  --set db.dsn="host=postgres.default user=humq password=humq dbname=humq" \
  --set kafka.bootstrapServers="kafka:9092" \
  --set jwt.secret="production-secret"
```

## 生产环境建议

1. **数据库**：使用 PostgreSQL 主从或云数据库
2. **JWT 密钥**：使用强随机密钥，定期轮换
3. **日志**：配置日志收集 (ELK/Loki)
4. **监控**：配置 Prometheus + Grafana 监控 HU MQ 自身
5. **备份**：定期备份 PostgreSQL 数据库
6. **HTTPS**：生产环境配置 TLS 证书