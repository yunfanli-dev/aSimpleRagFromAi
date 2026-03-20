# Test Notes

当前先维护最小人工验证闭环，便于在接 embedding 前稳定回归关键词检索。

## 手工验证样本

- `test/manual_kb/launch_checklist.md`
- `test/manual_kb/ops_notes.txt`

## 最短验证流程

1. 启动 PostgreSQL、Redis 和 API 服务。
2. 创建一个知识库。
3. 导入 `test/manual_kb/` 下的两个文档。
4. 调用 `GET /api/v1/documents/:id/chunks`，确认 chunk 顺序和内容正常。
5. 发起查询并检查 `citations`：
   - `startup order for the api and database`
   - `what storage systems are required for the local demo`
   - `what does the current query path return`
6. 核对返回结果：
   - `document_title` 与预期文档一致
   - `chunk_index` 顺序合理
   - `text` 为截断后的命中片段，不是整段原文堆叠
   - 已执行 `reindex` 的文档在查询时可走向量召回 + 关键词召回融合
   - `retrieval_source` 会标记 `keyword`、`vector` 或 `hybrid`
   - 返回结果不会被同一文档相邻 chunk 过度占满
   - `model` 当前会返回本地占位 LLM 名称
7. 如需排障，再检查 `query_logs` 是否写入命中的 chunk id。
8. 调用 `POST /api/v1/documents/:id/reindex`，确认返回 `embedded_count` 和 `embedding_model`。

## 当前目标

- 先验证关键词检索排序和 citation 展示是否稳定
- 在这套样本上继续做 embedding、向量召回和混合检索回归
- 后续补真实 rerank 与 LLM 接入后的回归基线
