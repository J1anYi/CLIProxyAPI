# PLAN: Add OpenAI Compatibility Proxy

## Phase Goal
在现有 Anthropic 代理基础上，新增 OpenAI 兼容代理，实现双协议支持，互不影响。

## Tasks

### Task 1: 更新配置文件
**Goal**: 添加 openai-compatibility 配置块

**File**: `config.yaml`

**Changes**:
```yaml
# 新增 OpenAI 兼容代理（无 prefix，通过端点区分）
openai-compatibility:
  - name: "jdcloud-openai"
    base-url: "https://modelservice.jdcloud.com/coding/openai/v1"
    api-key-entries:
      - api-key: "pk-be614585-f9cf-43bf-8d20-8350f3efc596"
    models:
      - name: "GLM-5"
        alias: "GLM-5"
      - name: "kimi-k2.5"
        alias: "kimi-k2.5"
      - name: "MiniMax-M2.5"
        alias: "MiniMax-M2.5"
```

**Verification**: 配置语法正确，服务能启动

---

### Task 2: 重启服务验证
**Goal**: 验证配置生效，两个代理都正常工作

**Steps**:
1. 停止现有服务
2. 启动服务（无需重新编译）
3. 检查日志确认 OpenAI-compat 加载

**Verification**: 日志显示 `X OpenAI-compat` 客户端加载

---

### Task 3: 测试现有代理
**Goal**: 确认现有 Anthropic 代理不受影响

**Test**: 使用现有 Claude Code 配置测试请求

**Verification**: 请求正常返回

---

### Task 4: 测试新 OpenAI 代理
**Goal**: 验证 OpenAI 代理正常工作

**Test**:
```bash
curl http://localhost:8317/v1/chat/completions \
  -H "Authorization: Bearer pk-be614585-f9cf-43bf-8d20-8350f3efc596" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-5",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

**Verification**: 请求正常返回，路由到京东云 OpenAI 端点

---

## Dependencies
- 无（配置级修改，无需代码变更）

## Rollback Plan
删除 `openai-compatibility` 配置块即可回滚

## Estimated Time
2 分钟

## Notes
- 这是一个纯配置变更，不涉及代码修改
- 通过端点区分：`/v1/messages` → Anthropic，`/v1/chat/completions` → OpenAI
- 模型名直接使用原始名称，无需前缀
