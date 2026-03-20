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
7. 如需排障，再检查 `query_logs` 是否写入命中的 chunk id。

## 当前目标

- 先验证关键词检索排序和 citation 展示是否稳定
- 后续在这套样本上继续做 embedding 接入前后的回归
