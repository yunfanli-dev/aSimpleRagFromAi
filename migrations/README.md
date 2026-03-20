# Migrations

当前包含 3 个初始化 migration：

- `0001_init_extensions.sql`: 启用 `vector` 和 `pgcrypto`
- `0002_init_schema.sql`: 初始化核心表
- `0003_init_indexes.sql`: 初始化检索和任务相关索引
- `0004_add_document_content.sql`: 给文档增加原始文本字段

建议执行顺序：

1. `0001_init_extensions.sql`
2. `0002_init_schema.sql`
3. `0003_init_indexes.sql`
4. `0004_add_document_content.sql`

当前核心表：

- `knowledge_bases`
- `documents`
- `chunks`
- `chunk_vectors`
- `ingest_jobs`
- `query_logs`
