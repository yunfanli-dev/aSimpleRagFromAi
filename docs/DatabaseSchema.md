# Database Schema

本文件用于说明当前项目数据库结构、表作用和字段含义。

## 规则

- 数据库结构发生变化时，必须检查本文件是否需要更新
- 如果新增表、删表、改字段、改索引、改约束，默认需要同步更新
- 目标是让后续 session 不用反复翻 migration 才能理解数据库

## 当前数据库

- 数据库名：`simplerag`
- 引擎：PostgreSQL 16
- 扩展：`pgcrypto`、`vector`
- 部署方式：Docker 容器 `simplerag-postgres`

## 表说明

### `knowledge_bases`

作用：知识库主表，表示一个独立知识集合。

字段：

- `id`: 主键，UUID
- `name`: 知识库名称
- `description`: 知识库说明
- `status`: 当前状态，默认 `active`
- `created_at`: 创建时间
- `updated_at`: 更新时间

### `documents`

作用：文档主表，记录文档元数据和原始文本内容。

字段：

- `id`: 主键，UUID
- `knowledge_base_id`: 所属知识库
- `source_type`: 文档类型，当前最小 ingest 支持 `txt`、`markdown`
- `title`: 文档标题
- `storage_path`: 原始存储路径，占位字段，后续可接文件存储
- `content_hash`: 文档内容哈希，用于幂等控制
- `content_text`: 归一化后的原始文本内容
- `status`: 文档状态，当前会经历 `pending` -> `indexed`
- `parse_version`: 解析/切片版本号
- `metadata_json`: 预留元数据
- `created_at`: 创建时间
- `updated_at`: 更新时间

约束：

- `(knowledge_base_id, content_hash)` 唯一，防止同一知识库重复导入同内容文档

### `chunks`

作用：文档切片表，保存每个文档拆分后的检索片段。

字段：

- `id`: 主键，UUID
- `document_id`: 所属文档
- `chunk_index`: 文档内切片序号
- `content`: 切片正文
- `token_count`: 粗略 token 数，当前按分词数近似
- `metadata_json`: 预留页码、段落层级等扩展信息
- `tsv`: PostgreSQL 生成列，用于全文检索
- `created_at`: 创建时间

约束：

- `(document_id, chunk_index)` 唯一，保证切片序号不重复

### `chunk_vectors`

作用：向量表，保存切片 embedding。

字段：

- `chunk_id`: 对应 `chunks.id`
- `embedding`: 向量字段，当前定义为 `vector(1024)`
- `embedding_model`: 所用 embedding 模型名
- `created_at`: 创建时间

说明：

- 当前已接入文档 reindex 时的最小 embedding 写入链路
- 现阶段使用本地 deterministic embedding provider 占位，后续可替换真实模型

### `ingest_jobs`

作用：离线导入任务表，用于后续 worker / 重试 / 异步任务处理。

字段：

- `id`: 主键，UUID
- `job_type`: 任务类型
- `target_id`: 目标对象 ID
- `status`: 任务状态
- `retry_count`: 重试次数
- `error_message`: 错误信息
- `scheduled_at`: 计划执行时间
- `created_at`: 创建时间
- `updated_at`: 更新时间

说明：

- 当前表已建好，但任务调度链路还未接入

### `query_logs`

作用：查询日志表，用于后续评测、回放、排障和统计。

字段：

- `id`: 主键，UUID
- `knowledge_base_id`: 所属知识库
- `question`: 用户问题
- `answer`: 模型回答
- `latency_ms`: 查询耗时
- `prompt_tokens`: prompt token 数
- `completion_tokens`: completion token 数
- `retrieved_chunk_ids`: 命中的切片 ID 列表
- `created_at`: 创建时间

说明：

- 当前已接入最小查询日志写入，记录问题、拼接答案、耗时和命中的 chunk id

## 当前索引

### `documents`

- `idx_documents_kb_id`: 按知识库查询文档
- `idx_documents_status`: 按状态筛选文档

### `chunks`

- `idx_chunks_document_id`: 按文档查切片
- `idx_chunks_tsv`: 全文检索索引

### `chunk_vectors`

- `idx_chunk_vectors_embedding`: HNSW 向量索引

### `ingest_jobs`

- `idx_ingest_jobs_status_scheduled`: 任务调度与扫描

### `query_logs`

- `idx_query_logs_kb_id_created_at`: 按知识库和时间回看查询

## 当前已接入的数据链路

- `knowledge_bases`: 已接入创建、列表、单条查询
- `documents`: 已接入创建、列表
- `chunks`: 已接入最小 ingest 写入
- `chunk_vectors`: 已接入文档 reindex 写入
- `query_logs`: 已接入最小查询日志写入

## 当前未接入但已预留

- `ingest_jobs`: 还未接 worker

## 对应 migration

- [0001_init_extensions.sql](../migrations/0001_init_extensions.sql)
- [0002_init_schema.sql](../migrations/0002_init_schema.sql)
- [0003_init_indexes.sql](../migrations/0003_init_indexes.sql)
- [0004_add_document_content.sql](../migrations/0004_add_document_content.sql)
