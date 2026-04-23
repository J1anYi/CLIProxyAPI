# Coding Conventions

## Language Style

- **Language**: Go 1.26.0
- **Style**: Standard Go formatting with gofmt
- **Imports**: Grouped by stdlib, external dependencies, internal packages
- **Line length**: No hard limit, but wrapped for readability

## Naming Conventions

### Packages
- Short, lowercase names: `config`, `auth`, `executor`, `translator`
- Package names describe the domain: `claude`, `codex`, `gemini`, `kimi`
- Sub-packages use descriptive domains: `config_access`, `model_registry`

### Types
- **Structs**: PascalCase with descriptive names: `ClaudeExecutor`, `OAuthError`, `AuthenticationError`
- **Interfaces**: PascalCase, often with `-er` suffix or descriptive: `apiKeyConfigEntry`
- **Type aliases**: Used for domain-specific types: `OAuthModelAlias`, `PayloadRule`

### Functions and Methods
- **PascalCase** for exported functions: `NewClaudeAuth`, `LoadConfig`, `RefreshTokens`
- **camelCase** for unexported functions: `parseCodeAndState`, `normalizeModelPrefix`
- **Constructor pattern**: `New<Type>()` for factory functions
- **Getters**: `GetAPIKey()`, `GetBaseURL()` - avoid `Get` prefix for simple field access

### Variables
- **camelCase** for local variables: `httpClient`, `tokenData`, `cfg`
- **PascalCase** for exported constants: `DefaultPprofAddr`, `AuthURL`
- **Error variables**: `ErrInvalidState`, `ErrCodeExchangeFailed` (PascalCase with `Err` prefix)
- **Acronyms**: Preserved as uppercase: `APIKey`, `OAuthURL`, `UserID`

### Constants
- Defined at package level with PascalCase for exported
- Grouped related constants: `AuthURL`, `TokenURL`, `ClientID`, `RedirectURI`
- Magic strings avoided; constants preferred

## Code Organization

### File Structure
- One primary type per file when practical
- Test files alongside source: `api_tools.go` → `api_tools_test.go`
- Platform-specific tests: `auth_files_download_windows_test.go`

### Package Structure
```
internal/
├── api/           # HTTP handlers, middleware, server
├── auth/          # Authentication providers (claude, codex, gemini, kimi)
├── config/        # Configuration loading and validation
├── runtime/       # Request execution logic
├── translator/    # Format conversion between API styles
├── logging/       # Structured logging utilities
├── cache/         # Caching implementations
├── watcher/       # File watching and hot reload
└── util/          # Shared utilities

sdk/               # Public SDK for external consumers
├── cliproxy/      # Core proxy types
├── translator/    # Translation interfaces
└── config/        # SDK configuration
```

### Struct Layout
- Configuration structs use YAML/json tags with snake_case for serialization
- Comments document each field with its purpose
- Embedded structs use `yaml:",inline"` for composition
- Sensitive fields use `json:"-"` to prevent serialization

## Error Handling

### Error Types
- Custom error types implement `error` interface: `OAuthError`, `AuthenticationError`
- Error types include: `Code`, `Message`, `Cause` fields
- `StatusCode()` method for HTTP status extraction

### Error Creation
```go
// Constructor pattern for errors
func NewOAuthError(code, description string, statusCode int) *OAuthError
func NewAuthenticationError(baseErr *AuthenticationError, cause error) *AuthenticationError
```

### Error Wrapping
- Use `fmt.Errorf("context: %w", err)` for wrapping
- Preserve original error as `Cause` field in custom types
- `errors.As()` for type checking: `IsAuthenticationError()`, `IsOAuthError()`

### Error Messages
- Include context and actionable information
- User-friendly messages via `GetUserFriendlyMessage()`
- Log debug details, return sanitized errors to clients

## Logging

### Library
- **logrus** for structured logging: `log "github.com/sirupsen/logrus"`
- Import alias `log` for convenience

### Log Levels
```go
log.Debug()  // Detailed debugging info
log.Info()   // Normal operations
log.Warn()   // Unexpected but handled
log.Error()  // Errors requiring attention
log.Fatal()  // Unrecoverable, exit
```

### Structured Fields
```go
log.WithFields(log.Fields{
    "request_id": requestID,
    "panic":      recovered,
}).Error("recovered from panic")
```

### HTTP Request Logging
- Middleware captures: method, path, status, latency, client IP
- Request IDs for AI API paths: `/v1/chat/completions`, `/v1/messages`, etc.
- Skip logging via `SkipGinRequestLogging()` for health checks

## Comments & Documentation

### Package Documentation
- Package comment at top of file describing purpose
```go
// Package config provides configuration management for the CLI Proxy API server.
// It handles loading and parsing YAML configuration files...
package config
```

### Function Documentation
- Doc comments for all exported functions
- Include Parameters, Returns sections
```go
// NewClaudeAuth creates a new Anthropic authentication service.
// It initializes the HTTP client with a custom TLS transport...
//
// Parameters:
//   - cfg: The application configuration containing proxy settings
//
// Returns:
//   - *ClaudeAuth: A new Claude authentication service instance
func NewClaudeAuth(cfg *config.Config) *ClaudeAuth
```

### Inline Comments
- Explain "why" not "what"
- Use for non-obvious logic, workarounds, or important context
- Reference issues/PRs when applicable

### TODO Comments
```go
// NOTE: Startup legacy key migration is intentionally disabled.
// Reason: avoid mutating config.yaml during server startup.
```

## Configuration

### YAML Structure
- Snake_case keys: `api-key`, `base-url`, `proxy-url`
- Nested structures for grouping: `remote-management`, `quota-exceeded`
- Optional fields use `omitempty` tags

### Defaults
- Set defaults before unmarshaling
- Validate and sanitize after loading
- Preserve comments when saving: `SaveConfigPreserveComments()`

### Validation
- Trim whitespace from string fields
- Remove empty/invalid entries during sanitization
- Normalize and deduplicate lists

## Testing Patterns

### Test File Naming
- `*_test.go` suffix
- Platform-specific: `*_windows_test.go`, `*_linux_test.go`

### Test Function Naming
- `Test<FunctionName>_<Scenario>` pattern
- Subtests with `t.Run()` for related cases

### Test Helpers
- `t.Helper()` for helper functions
- `t.Parallel()` for independent tests
- `t.Cleanup()` for resource cleanup
