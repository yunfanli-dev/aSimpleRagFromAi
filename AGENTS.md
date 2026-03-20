# AGENTS.md

本文件约束后续 agent 在本仓库中的工作方式。

## 启动规则

每次开启新的 session，必须先阅读以下文件：

1. [README.md](/home/ubuntu/aiCodingTest/aSimpleRagFromAi/aSimpleRagFromAi/README.md)
2. [ProgressLog.md](/home/ubuntu/aiCodingTest/aSimpleRagFromAi/aSimpleRagFromAi/ProgressLog.md)

目标：

- 先理解项目简略说明
- 明确当前项目做到哪一步了
- 从当前进度继续，不重复做已经完成的内容

## 阅读升级规则

如果仅通过 `README.md` 仍无法看明白项目目标、范围、架构或下一步方向，必须继续阅读：

3. [ProjectPlan.md](/home/ubuntu/aiCodingTest/aSimpleRagFromAi/aSimpleRagFromAi/ProjectPlan.md)

原则：

- 先读简略文档
- 不清楚再读完整项目书
- 不要在没读清楚上下文前直接改代码

## 执行规则

- 默认从 `ProgressLog.md` 记录的“当前 / 下一步”继续推进
- 每次只推进一个清晰的小阶段，避免一次改动过散
- 完成一部分后，必须更新 [ProgressLog.md](/home/ubuntu/aiCodingTest/aSimpleRagFromAi/aSimpleRagFromAi/ProgressLog.md)
- 更新进度时保持极简，只写：时间、完成、当前、下一步、备注

## 进度更新规则

每次完成阶段性工作后，必须同步更新进度，包括：

- 本次完成了什么
- 当前处于什么阶段
- 接下来准备做什么
- 是否有阻塞或验证限制

要求：

- 尽量简短
- 不写长篇分析
- 不重复抄项目书

## 工作优先级

1. 先看 `README.md`
2. 再看 `ProgressLog.md`
3. 看不明白再看 `ProjectPlan.md`
4. 按当前进度继续
5. 完成阶段性工作后更新 `ProgressLog.md`

## 禁止事项

- 不读取进度就直接开始大改
- 忽略 `ProgressLog.md` 里的当前状态
- 做完工作不更新进度
- 在 `README.md` 已足够清楚时无意义重复通读所有文档
