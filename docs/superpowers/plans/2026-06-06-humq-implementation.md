# HU MQ Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a full enterprise MQ monitoring platform (Go + Vue3) wrapping Apache Kafka, branded as HU MQ.

**Architecture:** Modular Go monolith (Gin + GORM + Sarama + JWT) embedding Vue3 SPA (Element Plus + ECharts). Deployed via Docker Compose / K8s Helm.

**Tech Stack:** Go 1.21+, Gin, GORM, Sarama, golang-jwt, bcrypt, PostgreSQL 15, Vue 3, Vite, Element Plus, ECharts, Axios

---

## Phase 1: Project Scaffolding

### Task 1: Initialize Go module and core structure
**Files:** `go.mod`, `cmd/humq-server/main.go`, `internal/config/config.go`, `internal/database/database.go`, `internal/api/router.go`

- [ ] **Step 1:** Run `go mod init github.com/wenli03/humq`; install gin, gorm, postgres driver; run `go mod tidy`
- [ ] **Step 2:** Write `cmd/humq-server/main.go` — entry point: load config, connect DB, migrate, setup router, listen
- [ ] **Step 3:** Write `internal/config/config.go` — Config struct with Server/DB/JWT/Kafka sections, environment variable loading via `os.Getenv`
- [ ] **Step 4:** Write `internal/database/database.go` — GORM connection function + AutoMigrate for all models
- [ ] **Step 5:** Write `internal/api/router.go` — placeholder Gin router with `/api/health` endpoint
- [ ] **Step 6:** Run `go build ./...` to verify
- [ ] **Step 7:** Commit: `feat: initialize Go module and project structure`

### Task 2: Initialize Vue 3 frontend
**Files:** `web/package.json`, `web/vite.config.js`, `web/index.html`, `web/src/main.js`, `web/src/App.vue`, `web/src/router/index.js`, `web/src/api/request.js`

- [ ] **Step 1:** Create `web/package.json` — deps: vue, vue-router, pinia, element-plus, echarts, axios, @element-plus/icons-vue; devDeps: vite, @vitejs/plugin-vue
- [ ] **Step 2:** Create `web/vite.config.js` — proxy `/api` to `localhost:8080`
- [ ] **Step 3:** Create `web/index.html` — lang="zh-CN", title "HU MQ"
- [ ] **Step 4:** Create `web/src/main.js` — init Vue app with ElementPlus (zhCn locale) + icons + router
- [ ] **Step 5:** Create `web/src/App.vue` — `<router-view />`
- [ ] **Step 6:** Create `web/src/router/index.js` — hash router with 10 routes: Login, Layout(w/children: Dashboard/Clusters/Topics/Consumers/Messages/Alerts/ACL/Ops/Guide)
- [ ] **Step 7:** Create `web/src/api/request.js` — axios instance with baseURL `/api/v1`, JWT interceptor, response code/401 handling
- [ ] **Step 8:** Run `npm install` in `web/`
- [ ] **Step 9:** Commit: `feat: initialize Vue 3 frontend with Element Plus`

---

## Phase 2: Data Models & Database

### Task 3: Define all GORM models
**Files:** `internal/database/models.go`

- [ ] **Step 1:** Create models.go with 9 structs: User, Cluster, TopicMeta, ConsumerGroup, MetricSnapshot, AlertRule, AlertEvent, DeadMessage, ACLRule — each with GORM tags, json tags, custom TableName()
- [ ] **Step 2:** Run `go build ./...`
- [ ] **Step 3:** Commit: `feat: define all GORM data models`

---

## Phase 3: Auth & JWT

### Task 4: Implement authentication
**Files:** `internal/auth/jwt.go`, `internal/auth/handler.go`, `internal/auth/middleware.go`, `internal/repository/user.go`, `internal/service/user.go`

- [ ] **Step 1:** Write jwt.go — Init(), GenerateToken(), GenerateRefreshToken(), ParseToken() using golang-jwt/jwt/v5 with Claims struct (userID, username, role)
- [ ] **Step 2:** Write handler.go — Login handler: bcrypt verify, generate token+refresh_token, return user info
- [ ] **Step 3:** Write middleware.go — AuthMiddleware (Bearer token parse, set user_id/username/role in context), AdminOnly (check role=="admin")
- [ ] **Step 4:** Write repository/user.go — FindByUsername, FindByID, List(page,pageSize), Create, Update, Delete
- [ ] **Step 5:** Write service/user.go — List, Create (bcrypt hash), Update, Delete
- [ ] **Step 6:** Run `go build ./...`
- [ ] **Step 7:** Commit: `feat: implement JWT authentication with user management`

---

## Phase 4: Kafka Integration & Repositories

### Task 5: Implement Kafka admin client wrapper
**Files:** `internal/kafka/admin.go`

- [ ] **Step 1:** Write admin.go — AdminClient struct wrapping sarama.Client+ClusterAdmin; methods: NewAdminClient, Close, ListTopics, DescribeTopic, CreateTopic, DeleteTopic, ListConsumerGroups, DescribeConsumerGroup, GetClusterInfo, FetchMessages, GetTopicOffsets, AlterTopicConfig, HealthCheck, ListBrokers
- [ ] **Step 2:** Run `go build ./...`
- [ ] **Step 3:** Commit: `feat: implement Kafka admin client wrapper with Sarama`

### Task 6: Implement all repositories
**Files:** `internal/repository/cluster.go`, `topic.go`, `consumer.go`, `metric.go`, `alert.go`, `deadmsg.go`, `acl.go`

- [ ] **Step 1:** Write cluster.go — List, Create, FindByID, Update, Delete
- [ ] **Step 2:** Write topic.go — ListByCluster(keyword filter), FindByName, Create, Delete, Update
- [ ] **Step 3:** Write consumer.go — ListByCluster, FindByGroupID, Upsert
- [ ] **Step 4:** Write metric.go — Insert, Query(time range), CleanupOlderThan
- [ ] **Step 5:** Write alert.go — ListRules, CreateRule, UpdateRule, DeleteRule, FindRuleByID, ListEnabledRules, CreateEvent, ListEvents(pagination)
- [ ] **Step 6:** Write deadmsg.go — List(pagination, topic filter), Create, Delete
- [ ] **Step 7:** Write acl.go — List, Create, Delete, FindByUserAndResource
- [ ] **Step 8:** Run `go build ./...`
- [ ] **Step 9:** Commit: `feat: implement all data repositories`

---

## Phase 5: Service Layer

### Task 7: Implement service layer
**Files:** `internal/service/cluster.go`, `topic.go`, `consumer.go`, `message.go`, `alert.go`, `acl.go`

- [ ] **Step 1:** Write cluster.go — ClusterService with in-memory AdminClient cache (map+sync.RWMutex), List/Create/GetClient/GetInfo/Delete
- [ ] **Step 2:** Write topic.go — List/Create/Delete/Describe/AlterConfig, each using clusterSvc.GetClient
- [ ] **Step 3:** Write consumer.go — List (fetch groups from Kafka + calculate lag + upsert to DB), Describe
- [ ] **Step 4:** Write message.go — Trace (fetch N messages from offset), ListDead (pagination), Replay (fetch dead msg, produce to target topic)
- [ ] **Step 5:** Write alert.go — ListRules/CreateRule/UpdateRule/DeleteRule/ListEvents
- [ ] **Step 6:** Write acl.go — List/Create/Delete
- [ ] **Step 7:** Run `go build ./...`
- [ ] **Step 8:** Commit: `feat: implement all service layers`

---

## Phase 6: API Handlers & Router

### Task 8: Implement API handlers and wire router
**Files:** `internal/api/handler.go`, `internal/api/router.go`

- [ ] **Step 1:** Write handler.go — response helpers (OK, OKPage, Fail) + all handler functions: ListClusters, CreateCluster, DeleteCluster, GetClusterInfo, ListTopics, CreateTopic, DeleteTopic, DescribeTopic, AlterTopicConfig, ListConsumers, DescribeConsumer, TraceMessages, ListDeadMessages, ReplayMessages, ListAlertRules, CreateAlertRule, UpdateAlertRule, DeleteAlertRule, ListAlertEvents, ListUsers, CreateUser, UpdateUser, DeleteUser, ListACLs, CreateACL, DeleteACL, GetCurrentUser. Each handler is a closure receiving the service as parameter.
- [ ] **Step 2:** Rewrite router.go — SetupRouter: init auth, create services, register all routes under `/api/v1` with auth middleware on protected routes, AdminOnly on sensitive endpoints
- [ ] **Step 3:** Run `go build ./...`
- [ ] **Step 4:** Commit: `feat: implement all API handlers and router`

---

## Phase 7: Frontend Pages

### Task 9: Create core pages (Layout, Login, Dashboard, Clusters)
**Files:** `web/src/views/Layout.vue`, `Login.vue`, `Dashboard.vue`, `Clusters.vue`

- [ ] **Step 1:** Write Layout.vue — el-container with el-aside sidebar (logo + el-menu with router links to all 8 sections using Element Plus icons), el-header (username + logout), el-main (router-view)
- [ ] **Step 2:** Write Login.vue — centered card with form (username+password), call POST /auth/login, store token+username, navigate to /dashboard
- [ ] **Step 3:** Write Dashboard.vue — 4 stat cards (clusters/topics/messages/alerts count), 2 ECharts (TPS line chart, Lag bar chart)
- [ ] **Step 4:** Write Clusters.vue — table listing all clusters, register cluster dialog, delete with confirmation
- [ ] **Step 5:** Commit: `feat: create layout, login, dashboard, and clusters pages`

### Task 10: Create data pages (Topics, Consumers, Messages)
**Files:** `web/src/views/Topics.vue`, `Consumers.vue`, `Messages.vue`

- [ ] **Step 1:** Write Topics.vue — cluster selector, keyword search, topic list table, create topic dialog (cluster/name/partitions/replicas), delete
- [ ] **Step 2:** Write Consumers.vue — cluster selector, consumer group list with topic tags, member count, state, lag (red if >1000)
- [ ] **Step 3:** Write Messages.vue — cluster selector, topic/partition/offset inputs, trace results table
- [ ] **Step 4:** Commit: `feat: create topics, consumers, and messages pages`

### Task 11: Create management pages (Alerts, ACL, Ops, Guide)
**Files:** `web/src/views/Alerts.vue`, `ACL.vue`, `Ops.vue`, `Guide.vue`

- [ ] **Step 1:** Write Alerts.vue — tabs for rules/events; rules: list + create dialog (cluster/metric/operator/threshold) + toggle/delete; events: list with level/status/time
- [ ] **Step 2:** Write ACL.vue — tabs for users/acls; users: list + create/edit/delete with role; acls: list + create with user/resource/operation
- [ ] **Step 3:** Write Ops.vue — rebalance trigger interface, cluster health status
- [ ] **Step 4:** Write Guide.vue — embedded documentation viewer with quick-start, sdk-guide, best-practices sections
- [ ] **Step 5:** Commit: `feat: create alerts, acl, ops, and guide pages`

---

## Phase 8: Docker & Deployment

### Task 12: Docker Compose and Dockerfile
**Files:** `deploy/Dockerfile`, `deploy/docker-compose.yml`

- [ ] **Step 1:** Write multi-stage Dockerfile — build Go binary, copy to distroless, embed Vue build artifacts in embedded filesystem
- [ ] **Step 2:** Write docker-compose.yml — 3 services: humq-server, postgres (15), redpanda (kafka-compatible, single-node)
- [ ] **Step 3:** Add `embed` directives in main.go to serve Vue dist/
- [ ] **Step 4:** Commit: `feat: add Dockerfile and docker-compose for development`

---

## Phase 9: Documentation

### Task 13: Write Chinese usage documentation
**Files:** `docs/README.md`, `quick-start.md`, `deployment.md`, `user-guide.md`, `sdk-guide.md`, `best-practices.md`, `faq.md`

- [ ] **Step 1:** Write README.md — product intro, features, tech stack, quick links
- [ ] **Step 2:** Write quick-start.md — 5-minute Docker Compose setup guide
- [ ] **Step 3:** Write deployment.md — detailed Docker and K8s Helm deployment
- [ ] **Step 4:** Write user-guide.md — operation manual covering all 10 modules
- [ ] **Step 5:** Write sdk-guide.md — Go/Java/Python SDK examples
- [ ] **Step 6:** Write best-practices.md — topic naming, partition planning, alert thresholds
- [ ] **Step 7:** Write faq.md — common issues and solutions
- [ ] **Step 8:** Commit: `docs: add Chinese usage documentation`

---

## Phase 10: Final Integration & Push

### Task 14: Build, verify, and push
- [ ] **Step 1:** Run `go build ./...` to verify Go code compiles
- [ ] **Step 2:** Run `npm run build` in web/ to verify frontend builds
- [ ] **Step 3:** Push all commits to GitHub
- [ ] **Step 4:** Create v1.0.0 tag and push
