---
title: refactor: Configure CLIProxyAPI as transparent proxy
type: refactor
status: active
date: 2026-04-20
---

# refactor: Configure CLIProxyAPI as transparent proxy

## Overview

当前 CLIProxyAPI 的 `config.yaml` 中 `api-keys` 包含硬编码占位符值，需要修改配置使其作为透明代理工作：
- CLIProxyAPI 不验证客户端 API key
- 直接转发请求到上游 API（京东云）
- 真实 API key 由 Claude Code 的 settings.json 直接传递

## Problem Frame

**当前架构：**
```
Claude Code → CLIProxyAPI (验证 api-keys) → 京东云 API
                    ↑
              需要配置 api-keys（硬编码问题）
```

**目标架构：**
```
Claude Code (携带 ANTHROPIC_AUTH_TOKEN) → CLIProxyAPI (透明转发) → 京东云 API
```

**Claude Code 配置（已完成）：**
```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "pk-be614585-f9cf-43bf-8d20-8350f3efc596",
    "ANTHROPIC_BASE_URL": "http://localhost:8317",
    ...
  }
}
```

## Requirements Trace

- R1. CLIProxyAPI 作为透明代理，不验证客户端 API key
- R2. 清理 config.yaml 中的硬编码占位符值
- R3. 保持向上游 API 正确转发 Authorization header

## Scope Boundaries

- 仅修改 CLIProxyAPI 配置
- 不修改 CLIProxyAPI 代码逻辑
- 不修改 Claude Code 配置（已完成）

## Context & Research

### Relevant Code and Patterns

**CLIProxyAPI 认证逻辑：**
- `api-keys` 配置：当为空或未设置时，跳过客户端认证
- 请求转发时：直接传递客户端的 Authorization header 到上游

**当前 config.yaml 问题：**
```yaml
api-keys:
  - "your-api-key-1"  # 硬编码占位符
  - "your-api-key-2"
  - "your-api-key-3"
```

## Key Technical Decisions

- **决策**: 将 `api-keys` 设为空数组，禁用客户端认证
- **理由**:
  - CLIProxyAPI 作为本地代理，无需验证客户端
  - 真实认证由上游 API（京东云）处理
  - 简化配置，避免硬编码问题

## Open Questions

### Resolved During Planning

- Q: 是否需要修改 CLIProxyAPI 代码？
- A: 不需要，只需修改配置文件

### Deferred to Implementation

- 无

## Implementation Units

- [ ] **Unit 1: 清理 config.yaml 中的 api-keys**

**Goal:** 移除硬编码占位符，设为空数组

**Requirements:** R1, R2

**Dependencies:** None

**Files:**
- Modify: `config.yaml`

**Approach:**
将 `api-keys` 改为空数组：
```yaml
# API keys for authentication
# Set to empty to disable client authentication (transparent proxy mode)
# The real API key is passed through in the Authorization header
api-keys: []
```

**Verification:**
- CLIProxyAPI 启动正常
- 请求可以正常转发

- [ ] **Unit 2: 验证透明代理工作正常**

**Goal:** 确认请求正确转发

**Requirements:** R3

**Dependencies:** Unit 1

**Files:**
- None (手动验证)

**Approach:**
1. 重启 CLIProxyAPI
2. 发送测试请求，验证 Authorization header 正确传递

**Verification:**
- `curl http://localhost:8317/v1/models -H "Authorization: Bearer pk-be614585-..."` 返回正确结果

## System-Wide Impact

- **Interaction graph:** 无变化
- **Error propagation:** 无影响
- **State lifecycle risks:** 无
- **API surface parity:** 无影响
- **Integration coverage:** 透明转发模式

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| 本地网络暴露 | CLIProxyAPI 仅监听 localhost，外部无法访问 |
| 向后兼容性 | 空数组是合法配置，表示禁用认证 |

## Documentation / Operational Notes

- CLIProxyAPI 作为透明代理运行
- 真实 API key 由 Claude Code 直接传递
- 无需在 CLIProxyAPI 中存储任何凭证

## Sources & References

- Claude Code settings.json 配置（已完成）
- CLIProxyAPI config.yaml
