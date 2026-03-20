# ProgressLog

用于记录当前项目推进到哪一步，尽量简短。

## 规则

- 只写结果，不写长解释
- 每次更新 4 项：时间、完成、当前、下一步
- 一条记录控制在几行内

## 记录

### 2026-03-20

- 完成：补充 `ProjectPlan.md`，重写 `README.md`
- 当前：项目规划阶段完成
- 下一步：初始化 Go 项目骨架与基础目录
- 备注：规划文档已提交并推送到 `main`

### 2026-03-20 skeleton

- 完成：初始化 Go 目录结构、API 骨架、`docker-compose.yml`
- 当前：基础工程搭建中
- 下一步：补充真实存储、迁移脚本、索引链路
- 备注：本机未安装 Go，暂未本地编译验证

### 2026-03-20 migration

- 完成：补充 PostgreSQL/pgvector 初版 migration、索引与执行脚本
- 当前：数据层初始化完成
- 下一步：接入 pgx repository，替换内存存储
- 备注：未执行数据库验证，当前仅完成 schema 落盘

### 2026-03-20 local-env

- 完成：安装 Go、Docker、Compose、psql，启动本地 Postgres/Redis，执行 migration
- 当前：本地开发环境可用
- 下一步：接入 pgx repository，替换内存存储
- 备注：`go mod tidy` 已完成，修复骨架编译小问题后继续验证

### 2026-03-20 pgx-repo

- 完成：接入 pgx 连接池，替换内存仓储，打通知识库/文档元数据真实读写
- 当前：元数据链路已接数据库
- 下一步：实现文本切片与 chunks 入库
- 备注：`go test ./...` 通过，知识库/文档接口已完成本地数据库验证
