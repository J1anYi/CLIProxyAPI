# System Architecture

## High-level Design

CLIProxyAPI is a **multi-provider AI API gateway** that provides a unified, OpenAI-compatible interface for multiple AI providers (Claude, Gemini, Codex/GPT, and custom OpenAI-compatible endpoints). It acts as a reverse proxy that:

1. **Accepts requests** in various formats (OpenAI chat/completions, Claude Messages, Gemini GenerateContent)
2. **Translates** requests between API formats using a bidirectional translation layer
3. **Routes** requests to appropriate upstream providers based on model names and available credentials
4. **Manages authentication** through OAuth flows, API keys, and token refresh mechanisms
5. **Returns responses** in the client's expected format, translating upstream responses as needed

The system follows a **middleware pipeline architecture** with clear separation between:
- HTTP routing (Gin framework)
- Request authentication (Access Manager)
- Credential selection (Auth Manager / Conductor)
- Request/response translation (Translator Registry)
- Upstream execution (Provider Executors)

## Core Components

### 1. API Server (`internal/api`)
- **Server**: Main HTTP server using Gin framework
- **Handlers**: Protocol-specific handlers for OpenAI, Claude, and Gemini endpoints
- **Middleware**: CORS, authentication, request logging
- **Routes**: Unified `/v1/` endpoints that route to appropriate translators

### 2. Auth Manager (`sdk/cliproxy/auth`)
The central credential orchestration system:
- **Conductor**: Orchestrates auth state, selection, and execution lifecycle
- **Scheduler**: Manages credential cooling/quota tracking
- **Selector**: Implements round-robin/fill-first credential selection strategies
- **Auto-refresh Loop**: Background token refresh for OAuth credentials
- **Session Affinity**: Routes same-session requests to the same credential

### 3. Translator System (`sdk/translator`, `internal/translator`)
Bidirectional translation between API formats:
- **Registry**: Maps (source format, target format) → transformation functions
- **Request Transformers**: Convert request payloads between formats
- **Response Transformers**: Convert streaming/non-streaming responses
- **Built-in Formats**: OpenAI, Claude, Gemini, Gemini-CLI, Codex, Antigravity

### 4. Access Manager (`sdk/access`)
Request-level authentication:
- **Provider Registry**: Pluggable authentication providers
- **Authentication Flow**: Validates API keys from request headers
- **Metadata Extraction**: Extracts context for routing decisions

### 5. Configuration System (`internal/config`)
- **Config Loading**: YAML-based configuration with hot-reload support
- **Provider Configuration**: API keys, OAuth settings, model aliases
- **Runtime Options**: Logging, TLS, management API settings

### 6. Token Store (`sdk/auth`, `internal/store`)
Multi-backend credential persistence:
- **File Store**: Local JSON files for OAuth tokens
- **PostgreSQL Store**: Database-backed storage for cloud deployments
- **Git Store**: Git-backed storage for distributed configs
- **Object Store**: S3-compatible object storage

### 7. Management API (`internal/api/handlers/management`)
RESTful API for runtime configuration:
- Config hot-reload
- Auth file management
- Usage statistics
- OAuth flow initiation

### 8. TUI (`internal/tui`)
Terminal-based management interface using Bubble Tea framework.

## Data Flow

```
Client Request
      │
      ▼
┌─────────────────┐
│   Gin Router    │ (routes.go)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Auth Middleware │ (access.Manager)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  API Handlers   │ (openai/claude/gemini)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Auth Conductor  │ (credential selection)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Translator    │ (format conversion)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Provider Executor│ (upstream HTTP calls)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Response Stream │ (translated response)
└─────────────────┘
```

## Request Lifecycle

1. **Request Ingress**: HTTP request hits Gin router at `/v1/chat/completions`, `/v1/messages`, or `/v1beta/models/*`

2. **Authentication**: Auth middleware validates API key via Access Manager against configured `api-keys`

3. **Model Resolution**: Handler extracts model name, resolves aliases, determines target provider

4. **Credential Selection**: Conductor picks available credential:
   - Checks quota/cooldown status
   - Applies session affinity if configured
   - Uses round-robin or fill-first strategy

5. **Request Translation**: Translator Registry finds appropriate transformer:
   - OpenAI → Gemini, Claude, Codex, etc.
   - Injects provider-specific headers and payload modifications

6. **Upstream Execution**: Provider Executor:
   - Constructs HTTP request with auth credentials
   - Sends to upstream API (Anthropic, Google, OpenAI)
   - Handles streaming via Server-Sent Events or WebSocket

7. **Response Translation**: Streams back with format conversion:
   - Upstream format → Client format
   - Streaming chunks translated on-the-fly

8. **Result Recording**: Auth state updated:
   - Success/failure tracked
   - Quota consumption noted
   - Cooldown scheduling if rate-limited

## Key Design Patterns

### 1. Registry Pattern
- **Translator Registry**: Format transformations registered by (from, to) key
- **Access Provider Registry**: Authentication providers registered by identifier
- Enables extensibility without core changes

### 2. Strategy Pattern
- **Credential Selector**: Pluggable selection algorithms (round-robin, fill-first)
- **Provider Executor**: Per-provider execution logic

### 3. Middleware Pipeline
- HTTP middleware chain for cross-cutting concerns
- Clean separation of authentication, logging, routing

### 4. Observer Pattern
- **Auth Hooks**: Lifecycle callbacks for auth events
- **Config Watcher**: Hot-reload on configuration changes

### 5. Factory Pattern
- **Token Store Factory**: Creates appropriate store backend
- **Executor Factory**: Creates provider-specific executors

### 6. Adapter Pattern
- **Translators**: Adapt between different API formats
- **Providers**: Wrap different upstream APIs behind common interface

## Extension Points

### Adding a New Provider
1. Implement `ProviderExecutor` interface in `sdk/cliproxy/` or `internal/`
2. Register translator functions for supported format conversions
3. Add configuration struct in `internal/config/config.go`
4. Register access provider if using custom auth

### Adding a New API Format
1. Define `Format` constant in `sdk/translator/format.go`
2. Implement `RequestTransform` function
3. Implement `ResponseTransform` functions (stream/non-stream)
4. Register in `internal/translator/<format>/` init functions

### Adding a New Token Store Backend
1. Implement `TokenStore` interface from `sdk/auth`
2. Add backend initialization in `cmd/server/main.go`
3. Register via `sdkAuth.RegisterTokenStore()`

### Adding Management API Endpoints
1. Add handler methods in `internal/api/handlers/management/`
2. Register routes in `internal/api/server.go` `registerManagementRoutes()`
3. Add configuration fields if needed

### Custom Access Provider
1. Implement `sdk/access.Provider` interface
2. Register via `sdkaccess.RegisterProvider()`
3. Reference in configuration

## Concurrency Model

- **Server**: Single HTTP server with goroutine-per-request
- **Auth Manager**: Thread-safe with RWMutex, atomic operations
- **Translator Registry**: Read-heavy with RWMutex protection
- **Auto-refresh**: Background goroutine pool with configurable concurrency
- **Streaming**: Each stream handled in dedicated goroutine

## Configuration Hot-Reload

The system supports runtime configuration updates:
1. File watcher detects config.yaml changes
2. Config reloaded and diffed against previous state
3. Auth Manager updated with new credentials
4. Access Providers reconciled
5. Translator Registry unchanged (format definitions are static)
