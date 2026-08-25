# 任务队列服务（task-queue）

一个基于 Go 标准库实现的轻量级任务队列后端服务，支持队列管理、任务投递/消费/确认/失败重试、死信队列、消费者心跳与订阅关系，以及全局统计。

## 技术栈

- 纯 Go 标准库（`net/http` + 标准库），零第三方依赖
- 内存存储（`sync.RWMutex` 保证并发安全）
- 标准分层架构：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`

## 运行

```bash
go run ./cmd/server
# 或
go build -o taskqueue-server ./cmd/server && ./taskqueue-server
```

默认监听 `:8080`，可通过环境变量配置：

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| `PORT` | `8080` | 监听端口 |
| `ADDR` | `:8080` | 完整监听地址（优先级高于 PORT） |
| `MAX_PAGE_SIZE` | `100` | 分页最大页大小 |
| `DEFAULT_MAX_RETRY` | `3` | 队列默认最大重试次数 |
| `BACKOFF_BASE_MS` | `1000` | 重试退避基数（毫秒） |
| `LOG_LEVEL` | `info` | 日志级别：debug/info/warn/error |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误码：400 校验失败、404 不存在、409 冲突、500 服务器错误。

### 队列 Queue

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/queues` | 创建队列（name/topic/max_retry/visibility_timeout） |
| GET | `/api/queues` | 列表（status/keyword + 分页） |
| GET | `/api/queues/{id}` | 详情 |
| PUT | `/api/queues/{id}` | 更新 |
| DELETE | `/api/queues/{id}` | 删除 |
| POST | `/api/queues/{id}/status` | 变更状态（active/paused） |

### 任务 Task

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/queues/{id}/tasks` | 投递任务（payload/priority） |
| POST | `/api/queues/{id}/dequeue` | 取出最高优先级待执行任务（置为 running） |
| GET | `/api/tasks` | 列表（queue_id/status + 分页） |
| GET | `/api/tasks/{id}` | 详情 |
| DELETE | `/api/tasks/{id}` | 删除 |
| POST | `/api/tasks/{id}/ack` | 确认成功（running → succeeded） |
| POST | `/api/tasks/{id}/fail` | 标记失败（reason，自动重试或进死信） |
| POST | `/api/tasks/{id}/requeue` | 死信重新投递（dead → pending） |

### 消费者 Consumer

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/consumers` | 创建消费者（name/concurrency） |
| GET | `/api/consumers` | 列表（status/keyword + 分页） |
| GET | `/api/consumers/{id}` | 详情 |
| PUT | `/api/consumers/{id}` | 更新 |
| DELETE | `/api/consumers/{id}` | 删除 |
| POST | `/api/consumers/{id}/heartbeat` | 心跳（更新 LastHeartbeat） |

### 订阅 Subscription

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/subscriptions` | 创建订阅（consumer_id/queue_id） |
| GET | `/api/subscriptions` | 列表（consumer_id/queue_id/status + 分页） |
| GET | `/api/subscriptions/{id}` | 详情 |
| PUT | `/api/subscriptions/{id}` | 更新 |
| DELETE | `/api/subscriptions/{id}` | 删除 |

### 重试 / 死信

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/retry-records` | 重试记录列表（task_id + 分页） |
| GET | `/api/retry-records/{id}` | 重试记录详情 |
| GET | `/api/dead-letters` | 死信列表（queue_id/keyword + 分页） |
| GET | `/api/dead-letters/{id}` | 死信详情 |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 全局统计 |

## 核心实体

1. **Queue**：任务队列（name/topic/max_retry），状态 `active/paused`。
2. **Task**：任务，状态机 `pending → running → succeeded/failed`，失败可 `pending`（重试）或 `dead`（死信），`dead → pending`（重新投递）。
3. **Consumer**：消费者实例，含心跳与并发度。
4. **Subscription**：消费者与队列的订阅关系。
5. **RetryRecord**：任务重试记录。
6. **DeadLetter**：死信（耗尽重试次数的任务）。

## 任务生命周期

1. `Enqueue` 投递：任务进入 `pending`，`MaxAttempts = MaxRetry + 1`。
2. `Dequeue` 取任务：按优先级（高优先）与入队时间取 `pending` 任务，置为 `running`。
3. `Ack` 确认：`running → succeeded`。
4. `Fail` 失败：`Attempts++`；若 `Attempts >= MaxAttempts` 则 `running → dead` 并写入死信，否则回 `pending` 并写重试记录（指数退避）。
5. `Requeue`：`dead → pending`，清空尝试次数并移除死信。

## 测试

```bash
go test ./...
```
