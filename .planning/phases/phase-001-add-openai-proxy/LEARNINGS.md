# LEARNINGS: Add OpenAI Compatibility Proxy

## 执行时间
2026-04-23

## 目标
在现有 Anthropic 代理基础上，新增 OpenAI 兼容代理，实现双协议支持。

---

## ✅ 成功经验

### 1. 端点隔离是关键
**发现**: CLIProxyAPI 通过**端点**而非**前缀**来区分不同协议的代理。

| 端点 | 协议 | Provider |
|------|------|----------|
| `/v1/messages` | Anthropic/Claude | `claude-api-key` |
| `/v1/chat/completions` | OpenAI | `openai-compatibility` |

**教训**: 不需要使用 `prefix` 来隔离路由，端点本身就是天然的隔离层。

### 2. 配置独立，互不干扰
**发现**: `claude-api-key` 和 `openai-compatibility` 是完全独立的配置块，可以：
- 使用相同的 API Key
- 使用相同的模型名
- 同时运行，互不影响

**教训**: 多 provider 配置可以放心共存，系统会自动路由到正确的端点。

### 3. 无需代码修改
**发现**: 这是一个纯配置变更，无需修改任何 Go 代码。

**教训**: CLIProxyAPI 的设计足够灵活，配置即可扩展新 provider。

### 4. 配置热重载
**发现**: 修改 config.yaml 后重启服务即可生效，无需重新编译。

**教训**: 配置变更快速验证，开发效率高。

---

## ❌ 失败/问题经验

### 1. 端口占用问题
**问题**: 启动新服务时，旧进程未停止导致端口 8317 被占用。

**错误信息**:
```
listen tcp :8317: bind: Only one usage of each socket address is normally permitted
```

**解决**: 使用 `taskkill //F //IM CLIProxyAPI.exe` 停止旧进程。

**教训**: 启动前先检查并停止旧进程，或使用不同的端口。

### 2. prefix 误解
**问题**: 最初认为需要使用 `prefix: "openai"` 来隔离路由，导致模型名需要加前缀如 `openai/GLM-5`。

**纠正**: 实际上不需要 prefix，因为端点已经隔离了路由。

**教训**: 先理解系统架构再设计配置，避免不必要的复杂性。

### 3. Git 中 config.yaml 不受跟踪
**问题**: 尝试 `git checkout -- config.yaml` 恢复文件，但文件不在 Git 跟踪中。

**原因**: `config.yaml` 在 `.gitignore` 中，只有 `config.example.yaml` 被跟踪。

**教训**: 配置文件通常不提交到 Git，使用 example 文件作为模板。

---

## 📊 最终配置

```yaml
# Anthropic 代理
claude-api-key:
  - api-key: "pk-be614585-..."
    base-url: "https://modelservice.jdcloud.com/coding/anthropic"
    models:
      - name: "GLM-5"
        alias: "GLM-5"
      # ... 其他模型

# OpenAI 代理（无 prefix）
openai-compatibility:
  - name: "jdcloud-openai"
    base-url: "https://modelservice.jdcloud.com/coding/openai/v1"
    api-key-entries:
      - api-key: "pk-be614585-..."
    models:
      - name: "GLM-5"
        alias: "GLM-5"
      # ... 其他模型
```

---

## 🎯 关键决策

| 决策 | 选择 | 原因 |
|------|------|------|
| 是否使用 prefix | 否 | 端点已隔离，无需前缀 |
| 模型名是否统一 | 是 | 两个代理使用相同模型名，通过端点区分 |
| 是否需要代码修改 | 否 | 纯配置变更 |

---

## 📝 使用指南

### Anthropic API（Claude Code 等）
```bash
curl http://localhost:8317/v1/messages \
  -H "Authorization: Bearer <API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"model": "GLM-5", "messages": [...]}'
```

### OpenAI API（OpenAI SDK 等）
```bash
curl http://localhost:8317/v1/chat/completions \
  -H "Authorization: Bearer <API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"model": "GLM-5", "messages": [...]}'
```

---

## 🔄 可复用模式

### 添加新 Provider 的步骤
1. 确定协议类型（Anthropic/OpenAI/Gemini/Vertex）
2. 在 config.yaml 中添加对应的配置块
3. 配置 base-url 和 API Key
4. 配置模型映射（name/alias）
5. 重启服务验证

### 无需修改代码的场景
- 添加新的 OpenAI 兼容 provider
- 添加新的 Claude API provider
- 添加模型别名
- 修改上游 URL

---

## 📈 影响评估

| 方面 | 影响 |
|------|------|
| 现有 Anthropic 代理 | 无影响 ✓ |
| 新增 OpenAI 代理 | 正常工作 ✓ |
| 性能 | 无影响（独立路由） |
| 维护成本 | 低（纯配置） |

---

## 总结

这次任务成功验证了 CLIProxyAPI 的多 provider 架构设计。关键发现是**端点隔离**而非**前缀隔离**，这使得配置更简洁，模型名更统一。整个过程无需代码修改，体现了系统良好的扩展性。
