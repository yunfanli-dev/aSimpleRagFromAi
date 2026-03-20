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

### 2026-03-20 chunks

- 完成：文档创建时支持文本切片，并将 `chunks` 实际写入 PostgreSQL
- 当前：最小 ingest 链路已打通
- 下一步：基于 `chunks.tsv` 接入关键词检索查询
- 备注：仅支持 `txt` 和 `markdown` 文本输入

### 2026-03-20 fts-query

- 完成：基于 PostgreSQL FTS 接入真实查询、返回真实 citations，并写入 `query_logs`
- 当前：最小检索闭环已打通
- 下一步：补文档详情 / chunk 查看接口，继续增强检索质量
- 备注：当前答案仍是检索片段拼接，不含 LLM 生成

### 2026-03-20 doc-read

- 完成：补文档详情和 chunk 查看接口，便于调试 ingest 和检索
- 当前：最小读写调试接口已齐
- 下一步：增强检索排序与结果质量，再准备接 embedding
- 备注：新增 `GET /api/v1/documents/:id` 和 `GET /api/v1/documents/:id/chunks`

### 2026-03-20 retrieval-quality

- 完成：增强 FTS 排序、补充 citation 字段、整理查询结果展示，并补最小人工验证样本
- 当前：关键词检索 demo 已更适合人工校验
- 下一步：基于当前检索结构准备 embedding 接入和向量检索链路
- 备注：待用本地样本继续做接口回归验证

### 2026-03-20 embedding-reindex

- 完成：接入本地 deterministic embedding provider，打通文档 reindex 到 `chunk_vectors` 写入
- 当前：离线 embedding 生成与向量落库入口已可用
- 下一步：实现向量召回，并与现有关键词检索做混合召回
- 备注：`POST /api/v1/documents/:id/reindex` 已返回 embedding 写入结果

### 2026-03-20 hybrid-retrieval

- 完成：查询链路接入最小向量召回，并与 FTS 结果做融合排序
- 当前：最小混合检索闭环已打通
- 下一步：增强融合权重、去重质量，并继续接 rerank / LLM
- 备注：当前向量召回基于本地 deterministic embedding provider

### 2026-03-20 retrieval-rerank

- 完成：接入最小 rerank、检索来源标记、相邻 chunk 去重和文档分散控制
- 当前：混合检索结果质量已进一步收敛
- 下一步：接真实 rerank / LLM，并补查询链路人工回归
- 备注：当前 rerank 仍为本地轻量规则实现

### 2026-03-20 llm-placeholder

- 完成：新增本地 extractive LLM provider，并将 answer 生成从 query service 中拆分
- 当前：查询链路已具备独立的 LLM 层占位结构
- 下一步：替换真实 LLM / rerank provider，并补 `reindex -> query` 人工回归
- 备注：响应中已返回当前 `model`

### 2026-03-20 test-readme

- 完成：重写 `test/README.md`，固定本地启动、迁移、API 启动和手工验证流程
- 当前：测试流程文档已可直接照着执行
- 下一步：后续新增测试方法时持续更新该文件
- 备注：已明确要求新测试方法必须同步写入 `test/README.md`

### 2026-03-20 openai-llm

- 完成：接入可切换的 OpenAI Responses API provider，并保留本地 LLM fallback
- 当前：查询链路已可切到真实 OpenAI answer 生成
- 下一步：接真实 rerank provider，并补 OpenAI 路径人工回归
- 备注：`LLM_PROVIDER=openai` 时需配置 `OPENAI_API_KEY`
