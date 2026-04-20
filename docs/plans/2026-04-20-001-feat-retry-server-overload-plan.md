---
title: feat: Add retry support for server overload errors
type: feat
status: active
date: 2026-04-20
---

# feat: Add retry support for server overload errors

## Overview

CLIProxyAPI 目前无法自动重试京东云 API 返回的 "Decode server is overloaded" 错误（HTTP 400）。这导致 Claude Code 任务被中断。需要扩展重试检测逻辑，将服务器过载错误识别为可重试错误。

## Problem Frame

用户报告京东云 API 返回以下错误时，CLIProxyAPI 没有自动重试：

```json
{
  "error": {
    "cause": "{\"detail\":\"Decode server is overloaded\"}",
    "code": 400,
    "message": "模型服务调用失败",
    "status": "FAILED_RESPONSE"
  },
  "requestId": "8b53d485fd30ba9f00bc7bb5ba5eb7bf-eUNmU",
  "result": null
}
```

当前行为：
- 错误直接返回给客户端
- Claude Code 任务中断

期望行为：
- CLIProxyAPI 自动重试
- Claude Code 保持等待状态
- 任务不中断

## Requirements Trace

- R1. 检测京东云 "Decode server is overloaded" 错误并触发重试
- R2. 保持与现有重试机制的兼容性
- R3. 可通过配置控制重试行为

## Scope Boundaries

- 仅修改错误检测逻辑，不改变重试流程本身
- 不添加新的配置项，复用现有 `request-retry` 配置

## Context & Research

### Relevant Code and Patterns

**核心文件：**
- `sdk/cliproxy/auth/conductor.go` - 包含 `isRateLimitDisguisedAs400` 函数
- `sdk/cliproxy/auth/rate_limit_400_test.go` - 测试文件

**现有检测模式：**
```go
// isRateLimitDisguisedAs400 当前检测：
// - ModelArts.81101 (华为云)
// - "type":"TooManyRequests"
// - "too many requests"
// - "rate limit exceeded" / "rate_limit exceeded"
```

**重试配置：**
```yaml
request-retry: 3           # 重试次数
max-retry-credentials: 0   # 尝试不同凭证
max-retry-interval: 30     # 最大等待时间(秒)
```

### Institutional Learnings

- 京东云错误格式：`{"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}","code":400,...}}`
- 错误嵌套在 `cause` 字段中，需要递归解析或字符串匹配

## Key Technical Decisions

- **决策**: 在 `isRateLimitDisguisedAs400` 函数中添加 "overload" / "overloaded" 检测
- **理由**: 
  - 服务器过载是临时状态，应该重试
  - 与现有检测模式一致
  - 最小化代码改动

## Open Questions

### Resolved During Planning

- Q: 是否需要添加新的配置项？
- A: 不需要，复用现有 `request-retry` 配置

### Deferred to Implementation

- 无

## Implementation Units

- [ ] **Unit 1: 扩展 isRateLimitDisguisedAs400 函数**

**Goal:** 添加对服务器过载错误的检测

**Requirements:** R1, R2

**Dependencies:** None

**Files:**
- Modify: `sdk/cliproxy/auth/conductor.go`
- Test: `sdk/cliproxy/auth/rate_limit_400_test.go`

**Approach:**
在 `isRateLimitDisguisedAs400` 函数中添加以下检测：
- 检测 "overload" 或 "overloaded" 关键词
- 检测京东云特定错误格式 `FAILED_RESPONSE` 状态

**Patterns to follow:**
- 现有 `isRateLimitDisguisedAs400` 函数的字符串匹配模式

**Test scenarios:**
- Happy path: "Decode server is overloaded" 被识别为可重试
- Happy path: "server overload" 被识别为可重试
- Edge case: 空错误返回 false
- Edge case: 不相关的 400 错误返回 false
- Integration: 京东云完整错误格式被正确识别

**Verification:**
- 单元测试通过
- 现有测试不受影响

- [ ] **Unit 2: 添加测试用例**

**Goal:** 验证新增的过载检测逻辑

**Requirements:** R1

**Dependencies:** Unit 1

**Files:**
- Modify: `sdk/cliproxy/auth/rate_limit_400_test.go`

**Approach:**
在 `TestIsRateLimitDisguisedAs400` 中添加京东云过载错误的测试用例

**Patterns to follow:**
- 现有测试用例格式

**Test scenarios:**
- Happy path: 京东云 "Decode server is overloaded" 错误被检测
- Happy path: 通用 "overload" 关键词被检测
- Error path: 不包含 overload 的错误不被检测

**Verification:**
- `go test ./sdk/cliproxy/auth/... -run TestIsRateLimitDisguisedAs400 -v` 通过

## System-Wide Impact

- **Interaction graph:** 无变化，仅影响错误检测逻辑
- **Error propagation:** 过载错误现在会触发重试，而不是直接返回
- **State lifecycle risks:** 无
- **API surface parity:** 无影响
- **Integration coverage:** 需要验证京东云错误格式

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| 过度重试导致服务器压力更大 | 使用现有 `max-retry-interval` 限制重试间隔 |
| 误判其他 400 错误为过载 | 仅检测特定关键词，保持精确匹配 |

## Documentation / Operational Notes

- 用户可通过 `request-retry` 配置控制重试次数
- 日志中会显示重试信息

## Sources & References

- 错误日志: `API Error: 400 {"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}"...}}`
- 相关代码: `sdk/cliproxy/auth/conductor.go` (`isRateLimitDisguisedAs400` 函数)
