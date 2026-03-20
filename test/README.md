# Test Notes

本文件是当前项目的固定测试入口。

目标：

- 统一本地启动 PostgreSQL、Redis、API 的方式
- 固定一套可重复执行的最小人工验证流程
- 后续只要新增测试方法、回归步骤或排障办法，就同步更新本文件

## 维护规则

- 每次新增可重复执行的测试方法，都更新本文件
- 每次修改启动命令、环境变量、迁移方式，也更新本文件
- 优先记录“能直接照着执行”的步骤，不写大段分析

## 手工验证样本

- `test/manual_kb/launch_checklist.md`
- `test/manual_kb/ops_notes.txt`

## 本地依赖启动

在仓库根目录执行：

```bash
docker compose up -d
```

或：

```bash
./scripts/dev.sh
```

启动后默认会拉起：

- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

检查容器状态：

```bash
docker compose ps
```

如果只想看 PostgreSQL 是否就绪：

```bash
docker compose exec postgres pg_isready -U postgres -d simplerag
```

## 数据库迁移

先设置本地数据库连接：

```bash
export POSTGRES_DSN='postgres://postgres:postgres@localhost:5432/simplerag?sslmode=disable'
```

执行迁移：

```bash
./scripts/apply_migrations.sh
```

或：

```bash
make migrate-local
```

## API 启动

推荐在仓库根目录执行：

```bash
POSTGRES_DSN='postgres://postgres:postgres@localhost:5432/simplerag?sslmode=disable' go run ./cmd/api
```

如果要走真实 OpenAI LLM provider，启动前追加：

```bash
export LLM_PROVIDER=openai
export OPENAI_API_KEY=your_api_key
export LLM_MODEL=gpt-5-mini
export LLM_REASONING_EFFORT=low
```

默认监听：

- API: `http://localhost:8080`

健康检查：

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## 固定手工验证流程

### 1. 创建知识库

```bash
KB_ID=$(
  curl -s -X POST http://localhost:8080/api/v1/kbs \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "manual-test-kb",
    "description": "manual validation kb"
  }' | jq -r '.data.id'
)
```

确认 `KB_ID` 已拿到真实 UUID：

```bash
echo "$KB_ID"
```

### 2. 导入测试文档

把 `test/manual_kb/launch_checklist.md` 和 `test/manual_kb/ops_notes.txt` 的内容分别作为 `content` 导入。
不要直接把 `$(cat file)` 塞进 JSON 字符串，优先用 `jq` 生成请求体。

示例：

```bash
DOC_ID_1=$(
  jq -n \
    --arg title "launch_checklist.md" \
    --arg source_type "markdown" \
    --arg content "$(cat test/manual_kb/launch_checklist.md)" \
    '{
      title: $title,
      source_type: $source_type,
      content: $content
    }' | \
  curl -s -X POST http://localhost:8080/api/v1/kbs/$KB_ID/documents \
    -H 'Content-Type: application/json' \
    -d @- | jq -r '.data.document.id'
)
```

```bash
DOC_ID_2=$(
  jq -n \
    --arg title "ops_notes.txt" \
    --arg source_type "txt" \
    --arg content "$(cat test/manual_kb/ops_notes.txt)" \
    '{
      title: $title,
      source_type: $source_type,
      content: $content
    }' | \
  curl -s -X POST http://localhost:8080/api/v1/kbs/$KB_ID/documents \
    -H 'Content-Type: application/json' \
    -d @- | jq -r '.data.document.id'
)
```

确认文档 ID：

```bash
echo "$DOC_ID_1"
echo "$DOC_ID_2"
```

### 3. 查看切片

```bash
curl http://localhost:8080/api/v1/documents/$DOC_ID_1/chunks
```

确认：

- chunk 顺序正常
- chunk 文本和原始文档一致
- 没有空 chunk

### 4. 执行 reindex

```bash
curl -X POST http://localhost:8080/api/v1/documents/$DOC_ID_1/reindex
curl -X POST http://localhost:8080/api/v1/documents/$DOC_ID_2/reindex
```

确认返回：

- `embedded_count` 大于 `0`
- `embedding_model` 已返回

### 5. 发起查询

建议固定使用这 3 个问题：

- `startup order for the api and database`
- `what storage systems are required for the local demo`
- `what does the current query path return`

示例：

```bash
curl -X POST http://localhost:8080/api/v1/query/debug \
  -H 'Content-Type: application/json' \
  -d '{
    "knowledge_base_id": "'"$KB_ID"'",
    "question": "what storage systems are required for the local demo"
  }'
```

### 6. 核对查询结果

重点检查：

- `answer` 已生成，不是空字符串
- `model` 已返回当前 answer provider 名称
- `citations[].document_title` 与预期文档一致
- `citations[].chunk_index` 顺序合理
- `citations[].text` 是截断后的命中片段，不是整段原文堆叠
- `citations[].retrieval_source` 会标记 `keyword`、`vector` 或 `hybrid`
- 已执行 `reindex` 的文档可走向量召回 + 关键词召回融合
- 结果不会被同一文档相邻 chunk 过度占满
- `debug_info.retrieved_chunks` 能看到检索候选与得分信息

### 7. 排障检查

必要时补查：

- `GET /api/v1/documents/:id`
- `GET /api/v1/documents/:id/chunks`
- PostgreSQL 中的 `query_logs`
- API 启动日志

## 当前回归范围

- 文档导入
- 文本切片
- reindex 向量写入
- FTS + vector 混合召回
- 本地轻量 rerank
- 本地占位或真实 OpenAI LLM answer 生成

## 后续要求

- 以后只要新增测试流程、HTTP 调用方式、固定问题、预期现象或排障步骤，就继续更新本文件
- 不要把新的测试方法只留在对话里，必须落到这里
