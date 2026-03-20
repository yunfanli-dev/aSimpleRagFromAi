# SimpleRAG-Go

一个从零实现的高质量、简单但完整的 RAG 系统，目标是用 Go 构建一个适合高并发场景的最小可用方案。

## 项目简写

`SRG` = `Simple Retrieval-Augmented Generation in Go`

## 项目目标

- 用尽可能少的组件实现可落地的 RAG 基线系统
- 明确区分离线索引链路与在线检索问答链路
- 优先保证并发能力、可观测性、可维护性
- 先做单机版，再保留未来水平扩展空间

## MVP 范围

- 文档导入：TXT、Markdown、PDF
- 文本清洗与切片
- Embedding 生成与向量写入
- 混合检索：向量检索 + 关键词检索
- Rerank
- Prompt 组装与大模型调用
- 引用片段返回
- REST API
- 基础监控、日志、压测与评测

## 建议技术栈

- 语言：Go 1.22+
- Web：Gin 或 Echo
- 并发控制：goroutine + worker pool + context
- 存储：PostgreSQL
- 向量：pgvector
- 全文检索：PostgreSQL FTS
- 缓存：Redis
- 队列：先使用库内任务队列，后续可替换 NATS / Kafka
- 观测：OpenTelemetry + Prometheus + Grafana

## 核心设计原则

- 简单优先，不为未来过度设计
- 在线路径短且稳定，避免重逻辑阻塞请求
- 离线处理异步化，支持重试、幂等和回放
- 数据模型清晰，便于后续扩展多租户与多数据源
- 每个阶段都有可验证的指标和测试

## 当前工程状态

- 已完成 Go 服务骨架
- 已完成 PostgreSQL / pgvector 初版 migration
- 已完成知识库 / 文档元数据的 PostgreSQL 基础接入
- 已完成文本切片与 `chunks` 入库最小链路
- 已完成基于真实 `chunks` 的关键词检索查询与日志落库
- 下一步是补文档详情 / chunk 查看接口，并继续增强检索质量

## 项目文档

完整规划与项目书见 [ProjectPlan.md](./ProjectPlan.md)。

项目进度记录见 [ProgressLog.md](./ProgressLog.md)。
