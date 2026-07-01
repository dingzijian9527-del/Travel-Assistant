# 检索增强问答向量库接入 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为旅行智能体接入第一阶段检索增强问答能力，先用内置旅游资料跑通链路。

**Architecture:** 新增 `pkg/rag` 包提供资料、检索和上下文格式化；配置层新增 `rag` 配置；网关智能体流式接口在调用大模型前检索资料并拼入消息。检索失败时降级为原有普通智能体回答。

**Tech Stack:** Go、Hertz（赫兹）、Viper（配置读取）、Zap（日志）、现有 `pkg/ai` 大模型客户端。

---

### Task 1: 配置层新增检索增强配置

**Files:**
- Modify: `pkg/config/config.go`
- Modify: `conf/config.yaml`
- Modify: `conf/config.yaml.tpl`
- Test: `pkg/config/config_test.go`

- [ ] **Step 1: Write the failing test**

新增测试，断言默认配置启用本地检索并包含集合名、地址、数量等默认值。

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/config`
Expected: FAIL，因为 `Config` 还没有 `RAG` 字段。

- [ ] **Step 3: Write minimal implementation**

新增 `RAGConfig`，挂到 `Config`，并设置默认值：

```go
type RAGConfig struct {
	Enabled        bool    `mapstructure:"enabled"`
	Provider       string  `mapstructure:"provider"`
	Address        string  `mapstructure:"address"`
	CollectionName string  `mapstructure:"collection_name"`
	EmbeddingDim   int     `mapstructure:"embedding_dim"`
	TopK           int     `mapstructure:"top_k"`
	MinScore       float64 `mapstructure:"min_score"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./pkg/config`
Expected: PASS。

### Task 2: 新增本地检索增强包

**Files:**
- Create: `pkg/rag/document.go`
- Create: `pkg/rag/knowledge.go`
- Create: `pkg/rag/retriever.go`
- Create: `pkg/rag/context.go`
- Test: `pkg/rag/retriever_test.go`
- Test: `pkg/rag/context_test.go`

- [ ] **Step 1: Write failing retrieval tests**

测试成都问题能命中成都资料，非旅游问题或分数不足时返回空结果。

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/rag`
Expected: FAIL，因为 `pkg/rag` 尚不存在。

- [ ] **Step 3: Implement local retriever**

实现 `Document`、`Result`、`Retriever` 接口、内置资料和基于关键词重叠的本地检索器。

- [ ] **Step 4: Add context formatting tests**

测试检索结果被格式化为带标题、城市、标签和正文的“可参考资料”文本。

- [ ] **Step 5: Run package tests**

Run: `go test ./pkg/rag`
Expected: PASS。

### Task 3: 智能体流式接口接入检索上下文

**Files:**
- Modify: `api/biz/handler/ai_stream.go`
- Test: `api/biz/handler/ai_stream_test.go`

- [ ] **Step 1: Write failing tests**

新增测试覆盖：有检索结果时，`buildAIContextMessages` 生成的消息包含“可参考资料”；检索失败时仍返回原始消息。

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./api/biz/handler`
Expected: FAIL，因为处理函数尚未接入检索上下文。

- [ ] **Step 3: Implement retrieval flow**

在 `ChatStream` 中根据配置创建本地检索器，检索用户问题，把结果格式化后传入消息构造函数。

- [ ] **Step 4: Run handler tests**

Run: `go test ./api/biz/handler`
Expected: PASS。

### Task 4: 全量验证

**Files:**
- No new files.

- [ ] **Step 1: Run all tests**

Run: `go test ./...`
Expected: PASS。

- [ ] **Step 2: Review changed files**

Run: `git diff --stat`
Expected: 只包含配置、`pkg/rag`、智能体处理和对应测试。

