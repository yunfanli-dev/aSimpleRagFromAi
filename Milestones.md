# Milestones

用于记录项目当前处在哪个里程碑、每个里程碑的大致状态，以及下一步重点。

## 规则

- 只保留阶段性里程碑，不记录细碎操作
- 每次项目改动后，都检查本文件是否需要更新
- 如果阶段推进、完成定义或下一步重点发生变化，必须同步更新本文件

## 当前里程碑

### Milestone 0: 项目初始化与本地环境

- 状态：已完成
- 范围：工程骨架、Docker Compose、本地 Postgres/Redis、migration 执行
- 结果：项目可本地启动，数据库基础结构可落盘

### Milestone 1: 最小 ingest 与关键词检索闭环

- 状态：已完成
- 范围：知识库/文档元数据入库、文本切片、`chunks` 写入、FTS 检索、查询日志、文档调试接口
- 结果：已具备最小可运行检索 demo

### Milestone 2: 检索质量增强与人工校验

- 状态：已完成
- 范围：FTS 排序优化、citation 增强、人工验证样本、固定测试流程、调试文档
- 结果：当前检索链路已更适合人工回归和排障

### Milestone 3: 向量链路与混合检索

- 状态：已完成
- 范围：本地 deterministic embedding、`reindex` 向量写入、向量召回、混合召回、轻量 rerank、去重与文档分散
- 结果：已打通最小混合检索链路

### Milestone 4: 真实 LLM 接入

- 状态：进行中
- 范围：从本地占位 answer 生成切换到真实大模型 provider
- 已完成：MiniMax provider 已接入，支持 `LLM_PROVIDER=minimax`
- 未完成：真实 rerank provider、MiniMax 路径人工回归、生成质量评估

## 当前重点

- 优先补真实 rerank provider
- 补 MiniMax 路径的固定人工回归
- 再考虑 PDF、评测与可观测性扩展
