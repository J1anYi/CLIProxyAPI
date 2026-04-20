---
title: feat: Support API keys from environment variables
type: feat
status: active
date: 2026-04-20
---

# feat: Support API keys from environment variables

## Overview

CLIProxyAPI 的 `config.yaml` 中 `api-keys` 字段包含硬编码的占位符值（如 `your-api-key-1`），这会导致：
1. Git 提交时可能意外包含真实 API key
2. 配置文件不能安全地提交到版本控制
3. 不同环境需要手动修改配置文件

需要支持通过环境变量配置 API keys，使配置文件可以安全提交到 Git。

## Problem Frame

**当前状态：**
```yaml
# config.yaml
api-keys:
  - "your-api-key-1"  # 硬编码，需要手动修改
  - "your-api-key-2"
  - "your-api-key-3"
```

**期望状态：**
```bash
# 通过环境变量传入
export CLIPROXYAPI_API_KEYS="key1,key2,key3"
```

```yaml
# config.yaml 可以保持为空或移除该字段
api-keys: []  # 安全，可以从环境变量读取
```

## Requirements Trace

- R1. 支持通过环境变量 `CLIPROXYAPI_API_KEYS` 配置 API keys（逗号分隔）
- R2. 环境变量优先级高于配置文件
- R3. 保持向后兼容：如果未设置环境变量，仍从配置文件读取
- R4. 更新 `.env.example` 文档说明

## Scope Boundaries

- 仅修改 API keys 的加载逻辑
- 不修改其他配置项的行为
- 不修改现有的 Management API 行为

## Context & Research

### Relevant Code and Patterns

**配置加载流程：**
1. `cmd/server/main.go` - 启动入口，已使用 godotenv 加载 `.env`
2. `internal/config/config.go` - 配置结构定义，`APIKeys []string`
3. `internal/config/sdk_config.go` - SDK 配置，`APIKeys []string`

**已有环境变量支持：**
```go
// internal/api/server.go:239
envAdminPassword, envAdminPasswordSet := os.LookupEnv("MANAGEMENT_PASSWORD")
```

**模式参考：** `MANAGEMENT_PASSWORD` 环境变量的处理方式

### Institutional Learnings

- 项目已使用 godotenv 加载 `.env` 文件
- 已有环境变量覆盖配置的模式（MANAGEMENT_PASSWORD）
- 配置文件使用 YAML 格式

## Key Technical Decisions

- **决策**: 使用单一环境变量 `CLIPROXYAPI_API_KEYS`，逗号分隔多个 key
- **理由**:
  - 简单直观，易于配置
  - 与 Docker/K8s 部署方式兼容
  - 避免需要预定义多个环境变量名

- **决策**: 环境变量优先级高于配置文件
- **理由**:
  - 允许在不修改配置文件的情况下覆盖配置
  - 符合 12-factor app 原则

## Open Questions

### Resolved During Planning

- Q: 使用单一环境变量还是多个环境变量？
- A: 使用单一环境变量 `CLIPROXYAPI_API_KEYS`，逗号分隔

- Q: 环境变量与配置文件的优先级？
- A: 环境变量优先，如果设置了环境变量则忽略配置文件中的值

### Deferred to Implementation

- 无

## Implementation Units

- [ ] **Unit 1: 添加环境变量读取逻辑**

**Goal:** 在配置加载时读取 `CLIPROXYAPI_API_KEYS` 环境变量

**Requirements:** R1, R2, R3

**Dependencies:** None

**Files:**
- Modify: `internal/config/config.go`
- Test: `internal/config/config_test.go`

**Approach:**
1. 在配置加载完成后检查 `CLIPROXYAPI_API_KEYS` 环境变量
2. 如果环境变量存在，解析逗号分隔的值并覆盖 `APIKeys`
3. 支持空格处理：`"key1, key2, key3"` 和 `"key1,key2,key3"` 都有效
4. **输入验证**：
   - 对每个 key 进行 `strings.TrimSpace()` 处理
   - 跳过空字符串和纯空白字符串
   - 最小 key 长度要求：8 字符
   - 无效 key 时记录警告日志（不包含 key 值）

**Patterns to follow:**
- 现有 `MANAGEMENT_PASSWORD` 的处理模式

**Test scenarios:**
- Happy path: 环境变量设置单个 key
- Happy path: 环境变量设置多个 keys（逗号分隔）
- Happy path: 环境变量设置多个 keys（带空格）
- Edge case: 环境变量为空字符串
- Edge case: 环境变量未设置，使用配置文件值
- Edge case: 环境变量包含空值 `",,key1,,"` - 跳过空值
- Edge case: 环境变量包含短于 8 字符的 key - 跳过并警告
- Integration: 环境变量覆盖配置文件值
- Security: key 值不出现在日志输出中

**Verification:**
- 单元测试通过
- 现有测试不受影响

- [ ] **Unit 2: 更新 .env.example 文档**

**Goal:** 添加 `CLIPROXYAPI_API_KEYS` 的说明

**Requirements:** R4

**Dependencies:** Unit 1

**Files:**
- Modify: `.env.example`

**Approach:**
在 `.env.example` 中添加：
```
# ------------------------------------------------------------------------------
# CLIProxyAPI API Keys
# ------------------------------------------------------------------------------
# If set, overrides the api-keys field in config.yaml.
# Multiple keys should be comma-separated.
# WARNING: Never commit the .env file to version control!
# CLIPROXYAPI_API_KEYS=key1,key2,key3
```

同时确认 `.env` 已在 `.gitignore` 中。

**Test scenarios:**
- Test expectation: none -- 文档更新

**Verification:**
- 文件包含正确的说明

- [ ] **Unit 3: 清理 config.yaml 中的硬编码值**

**Goal:** 移除配置文件中的占位符 API keys

**Dependencies:** Unit 1

**Files:**
- Modify: `config.yaml`
- Modify: `config.example.yaml`

**Approach:**
将 `api-keys` 改为空数组或移除该字段：
```yaml
# API keys for authentication
# Can also be set via CLIPROXYAPI_API_KEYS environment variable
api-keys: []
```

**Test scenarios:**
- Test expectation: none -- 配置文件更新

**Verification:**
- 配置文件不包含硬编码的 API key 值

## System-Wide Impact

- **Interaction graph:** 无变化，仅影响配置加载
- **Error propagation:** 无影响
- **State lifecycle risks:** 无
- **API surface parity:** 无影响
- **Integration coverage:** 需要验证环境变量加载

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| 环境变量泄露到日志 | 使用日志脱敏：不直接打印 config 结构体；在调试日志中使用 `[REDACTED]` 替代 key 值 |
| 向后兼容性破坏 | 保持配置文件加载作为 fallback |
| 空值或无效 key 被接受 | 输入验证：trim、跳过空值、最小长度检查 |
| .env 文件被提交到 Git | 确认 .env 在 .gitignore 中；在 .env.example 中添加警告注释 |

## Documentation / Operational Notes

- 用户可以通过环境变量配置 API keys
- Docker 部署时使用 `-e CLIPROXYAPI_API_KEYS=...`
- systemd 部署时在 service 文件中设置 `Environment=`
- **安全建议**：
  - 确保 `.env` 文件在 `.gitignore` 中
  - 设置 `.env` 文件权限：`chmod 600 .env`
  - 注意：环境变量在 Linux 上可通过 `/proc/<pid>/environ` 查看

## Sources & References

- 相关代码: `internal/config/config.go` (配置加载)
- 相关代码: `internal/api/server.go:239` (MANAGEMENT_PASSWORD 模式)
- 相关代码: `.env.example`
