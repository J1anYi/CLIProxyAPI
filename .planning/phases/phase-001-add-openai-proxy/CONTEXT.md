# Phase Context

## Goal
在现有 Anthropic 代理基础上，新增 OpenAI 兼容代理，不影响现有功能。

## Current State
- **现有代理**: `http://localhost:8317` → `https://modelservice.jdcloud.com/coding/anthropic`
- **协议**: Claude/Anthropic API (`/v1/messages`)
- **模型**: GLM-5 (opus), Kimi-K2.5 (haiku), MiniMax-M2.5 (sonnet)
- **API Key**: `pk-be614585-f9cf-43bf-8d20-8350f3efc596`

## Requirement
新增 OpenAI 代理：
- **路径**: `http://localhost:8317/v1/chat/completions`
- **上游**: `https://modelservice.jdcloud.com/coding/openai/v1`
- **协议**: OpenAI Chat Completions API
- **模型**: GLM-5, kimi-k2.5, MiniMax-M2.5（无前缀，直接使用原始名称）

## Impact Analysis

### 对现有代理的影响分析

**结论：无影响** ✓

原因：
1. **端点隔离**: 
   - Anthropic 使用 `/v1/messages` 端点
   - OpenAI 使用 `/v1/chat/completions` 端点
   - 两个端点完全独立，互不干扰
2. **独立凭证管理**: `openai-compatibility` 配置独立于 `claude-api-key`
3. **无 prefix**: 不设置 prefix，模型名直接使用原始名称

### 路由行为

| 请求端点 | 模型名 | 路由目标 |
|---------|--------|---------|
| `/v1/messages` | `GLM-5` | Anthropic 代理 |
| `/v1/messages` | `haiku` | Anthropic 代理 |
| `/v1/chat/completions` | `GLM-5` | OpenAI 代理 |
| `/v1/chat/completions` | `kimi-k2.5` | OpenAI 代理 |

### 注意事项

1. **端点差异**:
   - Anthropic 客户端使用 `/v1/messages` 端点
   - OpenAI 客户端使用 `/v1/chat/completions` 端点
2. **模型名**: 两个代理使用相同的模型名，但通过端点区分
3. **认证方式**: 使用相同的 API Key

## Configuration Change

```yaml
# 现有配置保持不变
claude-api-key:
  - api-key: "pk-be614585-f9cf-43bf-8d20-8350f3efc596"
    base-url: "https://modelservice.jdcloud.com/coding/anthropic"
    models:
      - name: "GLM-5"
        alias: "GLM-5"
      # ... 其他模型

# 新增 OpenAI 兼容代理（无 prefix）
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

## Risk Assessment

| 风险 | 等级 | 缓解措施 |
|-----|------|---------|
| 影响现有代理 | 低 | 端点隔离确保路由独立 |
| API 格式兼容 | 中 | 验证上游 OpenAI API 兼容性 |
| 模型名冲突 | 无 | 端点不同，不会冲突 |

## Success Criteria
- [ ] 现有 Anthropic 代理正常工作 (`/v1/messages`)
- [ ] OpenAI 代理正常工作 (`/v1/chat/completions`)
- [ ] 模型名 `GLM-5`, `kimi-k2.5` 等能正确路由
- [ ] 两个代理互不干扰
