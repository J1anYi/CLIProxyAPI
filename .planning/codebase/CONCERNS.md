# Technical Concerns

## Security Considerations

### Authentication & Secrets Management
- **API keys in config files**: `config.yaml` stores API keys in plaintext (`api-keys`, `gemini-api-key`, `claude-api-key`, `codex-api-key`, `vertex-api-key`, `openai-compatibility`)
- **Management secret-key**: Plaintext management key is hashed on startup using bcrypt, but initial value is visible in config
- **Environment variable secrets**: `MANAGEMENT_PASSWORD` env var overrides config (documented in `.env.example`)
- **Postgres DSN**: Contains credentials in connection string (`PGSTORE_DSN=postgresql://user:pass@host:port/db`)
- **Git store token**: `GITSTORE_GIT_TOKEN` for repository access
- **Object store keys**: `OBJECTSTORE_ACCESS_KEY` and `OBJECTSTORE_SECRET_KEY` for S3-compatible storage
- **OAuth tokens**: Stored in files under `~/.cli-proxy-api` or configured auth directory

### Mitigations in Place
- bcrypt hashing for management passwords (`golang.org/x/crypto/bcrypt`)
- Constant-time password comparison using `crypto/subtle` package
- IP-based rate limiting and temporary blocking for management API
- Configurable TLS support (disabled by default)
- Path traversal protection in auth file operations (`filepath.Rel` checks for `..` prefix)

### Areas Needing Attention
- TLS disabled by default (`tls.enable: false`)
- WebSocket authentication optional (`ws-auth: false`)
- `allow-remote: false` default for management API (good, but should be enforced)
- No encryption at rest for stored OAuth tokens
- Secret keys logged in debug mode should be masked

## Performance Concerns

### Concurrency Patterns
- Heavy mutex usage throughout codebase (40+ `sync.Mutex`/`sync.RWMutex` instances)
- Potential lock contention in hot paths:
  - `sdk/cliproxy/auth/conductor.go` - auth credential management
  - `sdk/cliproxy/auth/selector.go` - credential selection
  - `internal/api/handlers/management/handler.go` - management request handling
  - `sdk/access/manager.go` - access provider registry

### Retry & Backoff
- Up to 16 retries per request (`request-retry: 16`)
- Exponential backoff up to 30 minutes for quota errors
- Rate limit backoff with jitter (base: 1s, max: 30s)
- `max-retry-credentials: 0` tries all available credentials (potential thundering herd)

### Caching
- Signature cache for thinking blocks (`internal/cache/signature_cache.go`)
- Session affinity cache with TTL (default: 1h)
- User ID cache per API key
- Device profile fingerprint caching

### Memory Management
- `commercial-mode: false` default - enables higher-overhead middleware features
- In-memory usage statistics aggregation (optional)
- Failed IP attempt tracking with periodic cleanup (1h interval, 2h max idle)

## Scalability

### Storage Backend Options
- **File-based**: Default, local filesystem (`~/.cli-proxy-api`)
- **PostgreSQL**: `internal/store/postgresstore.go` with local workspace mirroring
- **Git-backed**: `internal/store/gitstore.go` for version-controlled config
- **S3-compatible**: `internal/store/objectstore.go` for object storage

### Horizontal Scaling Limitations
- Session affinity binds requests to specific credentials (can limit load distribution)
- In-memory state for rate limiting and cooldowns (not distributed)
- Model registry updates via background goroutine (per-instance)
- Auth file watcher uses fsnotify (per-instance file system events)

### Configuration for Scale
- `auth-auto-refresh-workers: 16` (configurable worker pool)
- `routing.strategy`: "round-robin" or "fill-first"
- `session-affinity: false` default (allows better distribution)

## Technical Debt

### Deprecated/Legacy Patterns
- `internal/api/modules/modules.go`: `RouteModule` interface marked as DEPRECATED, use `RouteModuleV2`
- `NewLegacy()` function for Amp module backward compatibility
- Multiple "sanitize" functions suggesting legacy data format migrations:
  - `SanitizeGeminiKeys()`, `SanitizeVertexCompatKeys()`, `SanitizeCodexKeys()`
  - `SanitizeClaudeHeaderDefaults()`, `SanitizeOAuthModelAlias()`

### Workarounds
- OAuth tool rename map (`oauthToolRenameMap`) to bypass Anthropic third-party client detection
- Cloaking configuration for non-Claude-Code clients
- Device profile stabilization for fingerprint consistency
- Model alias and prefix systems for routing flexibility

### Code Complexity
- `internal/runtime/executor/claude_executor.go`: Handles multiple encoding formats (gzip, brotli, zstd)
- `sdk/cliproxy/auth/conductor.go`: Large file managing auth lifecycle
- `internal/api/handlers/management/auth_files.go`: 2500+ lines handling auth file operations
- `cmd/server/main.go`: 580+ lines in main function with multiple code paths

## Dependencies at Risk

### Outdated Dependencies (Go 1.26.0)
| Package | Current | Latest | Risk Level |
|---------|---------|--------|------------|
| `github.com/gin-gonic/gin` | v1.10.1 | v1.12.0 | Medium |
| `github.com/cloudflare/circl` | v1.6.1 | v1.6.3 | High (crypto) |
| `github.com/ProtonMail/go-crypto` | v1.3.0 | v1.4.1 | High (crypto) |
| `golang.org/x/crypto` | v0.45.0 | - | Current |
| `github.com/cloudwego/base64x` | v0.1.4 | v0.1.6 | Low |
| `github.com/andybalholm/brotli` | v1.0.6 | v1.2.1 | Low |
| `github.com/bytedance/sonic` | v1.11.6 | v1.15.0 | Low |
| `cloud.google.com/go/compute/metadata` | v0.3.0 | v0.9.0 | Medium |

### Direct Dependencies (No Known Issues)
- `golang.org/x/oauth2` v0.30.0
- `github.com/gorilla/websocket` v1.5.3
- `github.com/jackc/pgx/v5` v5.7.6
- `github.com/minio/minio-go/v7` v7.0.66

## Error-prone Areas

### Panic/Recover Patterns
- `internal/registry/model_registry.go`: Panics in hook execution are recovered
- `internal/api/modules/amp/routes.go`: Panic recovery in request handling
- `internal/logging/gin_logger.go`: Panics with `http.ErrAbortHandler`
- `examples/` directory: Multiple panics in example code (expected)

### Context Usage
- `context.Background()` used in 20+ files - some may miss cancellation propagation
- Short timeouts (30s) for database bootstrap operations
- Background goroutines for watchers, updaters, cleanup tasks

### File System Operations
- `os.RemoveAll` used in `postgresstore.go` for auth directory reset
- Path traversal validation relies on `filepath.Rel` prefix checking
- Windows path handling via `filepath.FromSlash`/`filepath.ToSlash`

### Complex Logic
- Credential selection with multiple blocking conditions (cooldown, disabled, model cooldown)
- Multi-provider request routing with alias/prefix resolution
- Thinking block signature validation with protobuf structure checks

## Missing Features

### Security Features
- No audit logging for management API operations
- No encryption at rest for stored credentials
- No key rotation automation
- No MFA support for management access
- No rate limiting per API key (only per IP for management)

### Operational Features
- No health check endpoint for load balancers
- No metrics export (Prometheus/OpenMetrics format)
- No distributed tracing support
- No graceful degradation mode when all credentials exhausted

### Protocol Support
- Thinking parameters support varies by provider
- Streaming error recovery varies by provider

## Documentation Gaps

### Missing Package Documentation
- No README for `internal/` packages
- Limited godoc comments for exported functions
- Complex algorithms (signature validation, credential selection) lack inline documentation

### Configuration Documentation
- `config.example.yaml` is comprehensive but some options lack explanations:
  - `antigravity-signature-cache-enabled`
  - `antigravity-signature-bypass-strict`
  - `experimental-cch-signing`
- SDK documentation in `docs/` but limited architecture overview

### Operational Documentation
- No deployment guide for production environments
- No troubleshooting guide for common errors
- No performance tuning documentation
- No migration guide between storage backends
