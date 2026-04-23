# Project Structure

## Directory Layout

```
CLIProxyAPI/
├── cmd/                          # Application entry points
│   ├── server/                   # Main proxy server
│   │   └── main.go              # Server entry point
│   └── fetch_antigravity_models/ # Utility command
│       └── main.go
│
├── internal/                     # Private application code
│   ├── access/                   # Access provider reconciliation
│   │   ├── config_access/       # Config-based access providers
│   │   └── reconcile.go         # Provider diff/apply logic
│   │
│   ├── api/                      # HTTP API implementation
│   │   ├── server.go            # Gin server setup
│   │   ├── handlers/
│   │   │   └── management/      # Management API handlers
│   │   ├── middleware/          # HTTP middleware
│   │   └── modules/
│   │       └── amp/             # Amp CLI integration module
│   │
│   ├── auth/                     # Provider-specific OAuth implementations
│   │   ├── antigravity/         # Antigravity OAuth
│   │   ├── claude/              # Claude/Anthropic OAuth
│   │   ├── codex/               # OpenAI Codex OAuth
│   │   ├── gemini/              # Google/Gemini OAuth
│   │   ├── kimi/                # Kimi OAuth
│   │   └── vertex/              # Vertex AI credentials
│   │
│   ├── browser/                  # Browser automation utilities
│   ├── buildinfo/               # Build version information
│   ├── cache/                   # Signature cache for thinking blocks
│   ├── cmd/                     # CLI command implementations
│   │   ├── run.go              # Server run logic
│   │   ├── login.go            # Login commands
│   │   └── auth_manager.go     # Auth management
│   │
│   ├── config/                  # Configuration loading/parsing
│   │   ├── config.go           # Main config structures
│   │   └── sdk_config.go       # SDK-exposed config
│   │
│   ├── constant/                # Application constants
│   ├── interfaces/              # Internal interfaces
│   ├── logging/                 # Logging utilities
│   ├── managementasset/         # Management UI assets
│   ├── misc/                    # Miscellaneous utilities
│   ├── registry/                # Model registry
│   ├── runtime/                 # Runtime utilities
│   ├── store/                   # Token store backends
│   │   ├── postgresstore.go    # PostgreSQL backend
│   │   ├── gitstore.go         # Git backend
│   │   └── objectstore.go      # S3-compatible backend
│   │
│   ├── thinking/                # Thinking/reasoning utilities
│   ├── translator/              # Format translators
│   │   ├── init.go             # Translator registration
│   │   ├── antigravity/        # Antigravity format translators
│   │   ├── claude/             # Claude format translators
│   │   ├── codex/              # Codex format translators
│   │   ├── gemini/             # Gemini format translators
│   │   ├── gemini-cli/         # Gemini CLI format translators
│   │   └── openai/             # OpenAI format translators
│   │
│   ├── tui/                     # Terminal UI (Bubble Tea)
│   ├── usage/                   # Usage statistics
│   ├── util/                    # General utilities
│   ├── watcher/                 # File watching for hot-reload
│   └── wsrelay/                 # WebSocket relay
│
├── sdk/                          # Public SDK for embedding
│   ├── access/                  # Request authentication SDK
│   │   ├── manager.go          # Access manager
│   │   ├── registry.go         # Provider registry
│   │   └── types.go            # Auth types
│   │
│   ├── api/                     # API handler SDK
│   │   ├── handlers/           # Base handlers
│   │   │   ├── openai/         # OpenAI handlers
│   │   │   ├── claude/         # Claude handlers
│   │   │   └── gemini/         # Gemini handlers
│   │   └── management.go       # Management SDK
│   │
│   ├── auth/                    # Auth SDK interfaces
│   │   ├── interfaces.go       # Authenticator interface
│   │   ├── filestore.go        # File token store
│   │   └── manager.go          # Auth manager
│   │
│   ├── cliproxy/               # Core proxy SDK
│   │   ├── auth/               # Auth conductor system
│   │   │   ├── conductor.go   # Main orchestrator
│   │   │   ├── scheduler.go   # Quota/cooling
│   │   │   ├── selector.go    # Credential selection
│   │   │   └── store.go       # Auth store
│   │   ├── executor/           # Execution context
│   │   ├── pipeline/           # Request pipeline
│   │   ├── service.go          # Service builder
│   │   └── types.go            # Core types
│   │
│   ├── config/                 # Config SDK
│   ├── logging/                # Logging SDK
│   ├── proxyutil/              # Proxy utilities
│   └── translator/             # Translator SDK
│       ├── registry.go         # Translation registry
│       ├── types.go            # Transformer types
│       └── format.go           # Format definitions
│
├── examples/                    # Example code
│   ├── custom-provider/        # Custom provider example
│   ├── http-request/           # HTTP request example
│   └── translator/             # Custom translator example
│
├── auths/                       # Default auth token directory
├── static/                      # Static assets
├── docs/                        # Documentation
│   ├── sdk-usage.md
│   ├── sdk-advanced.md
│   ├── sdk-access.md
│   └── sdk-watcher.md
│
├── config.yaml                  # Main configuration file
├── config.example.yaml          # Example configuration
├── go.mod                       # Go module definition
├── Dockerfile                   # Docker build
├── docker-compose.yml           # Docker Compose
└── README.md                    # Project documentation
```

## Key Files

### Entry Points
- `cmd/server/main.go` - Main server entry point, config loading, service startup
- `cmd/fetch_antigravity_models/main.go` - Utility to fetch Antigravity model catalog

### Core Server
- `internal/api/server.go` - Gin server, route setup, middleware configuration
- `internal/config/config.go` - Configuration structures and YAML loading
- `internal/cmd/run.go` - Server lifecycle management

### Auth System
- `sdk/cliproxy/auth/conductor.go` - Central credential orchestrator
- `sdk/cliproxy/auth/selector.go` - Credential selection strategies
- `sdk/cliproxy/auth/scheduler.go` - Quota tracking and cooldown
- `sdk/access/manager.go` - Request authentication manager

### Translation Layer
- `sdk/translator/registry.go` - Format transformation registry
- `internal/translator/init.go` - Built-in translator registration
- `internal/translator/<source>/<target>/` - Format-specific transformers

### Handlers
- `sdk/api/handlers/handlers.go` - Base handler with common logic
- `sdk/api/handlers/openai/openai_handlers.go` - OpenAI-compatible endpoints
- `sdk/api/handlers/claude/code_handlers.go` - Claude Messages endpoints
- `sdk/api/handlers/gemini/gemini_handlers.go` - Gemini GenerateContent endpoints

### Management
- `internal/api/handlers/management/handler.go` - Management API handler
- `internal/api/handlers/management/config_*.go` - Config manipulation handlers
- `internal/api/handlers/management/auth_files.go` - Auth file management

## Module Organization

### `internal/` - Private Packages
Contains implementation details not exposed to external consumers:
- Provider-specific OAuth flows
- Internal translators
- Storage backends
- HTTP handlers

### `sdk/` - Public SDK
Stable APIs for embedding CLIProxyAPI in other applications:
- `sdk/cliproxy/` - Core proxy functionality
- `sdk/auth/` - Authentication interfaces
- `sdk/translator/` - Format translation
- `sdk/access/` - Request authentication

### `cmd/` - Executables
Each subdirectory is a separate executable:
- `server/` - Main proxy server
- `fetch_antigravity_models/` - Utility command

### `examples/` - Example Code
Demonstrates SDK usage patterns:
- Custom providers
- Custom translators
- HTTP request customization

## Entry Points

### Server Entry Point (`cmd/server/main.go`)
```
main()
  ├── Parse CLI flags (-login, -codex-login, -claude-login, etc.)
  ├── Load environment variables (.env)
  ├── Initialize token store (file/postgres/git/object)
  ├── Load configuration (config.yaml)
  ├── Register access providers
  ├── Handle login commands if specified
  └── Start server (cmd.StartService or TUI mode)
```

### API Routes
| Route | Handler | Description |
|-------|---------|-------------|
| `GET /healthz` | inline | Health check |
| `GET /v1/models` | unifiedModelsHandler | List available models |
| `POST /v1/chat/completions` | openaiHandlers | OpenAI chat API |
| `POST /v1/messages` | claudeCodeHandlers | Claude Messages API |
| `POST /v1/responses` | openaiResponsesHandlers | OpenAI Responses API |
| `POST /v1beta/models/*` | geminiHandlers | Gemini GenerateContent API |
| `GET /v0/management/*` | mgmt handlers | Management API |

## Configuration Files

### `config.yaml` - Main Configuration
Contains:
- Server settings (port, host, TLS)
- API keys for authentication
- Provider configurations (Gemini, Claude, Codex, OpenAI-compatible)
- OAuth settings and model aliases
- Routing strategy and retry settings

### `.env` - Environment Variables
Optional overrides for:
- `PGSTORE_DSN` - PostgreSQL connection
- `GITSTORE_GIT_URL` - Git token store
- `OBJECTSTORE_*` - Object store settings
- `MANAGEMENT_PASSWORD` - Management API secret

### `auths/` - Auth Token Directory
Contains OAuth token files:
- `gemini-*.json` - Gemini CLI tokens
- `claude-*.json` - Claude tokens
- `codex-*.json` - Codex tokens
- `antigravity-*.json` - Antigravity tokens

### `config.example.yaml` - Example Configuration
Template with all available options documented.

## Import Graph (Simplified)

```
cmd/server/main.go
    ├── internal/config
    ├── internal/cmd
    │       └── internal/api
    │               ├── sdk/api/handlers
    │               │       ├── sdk/translator
    │               │       └── sdk/cliproxy/auth
    │               ├── sdk/access
    │               └── internal/api/handlers/management
    ├── sdk/auth (TokenStore)
    └── internal/access
            └── sdk/access
```

## Test Organization

Tests are co-located with implementation:
- `*_test.go` files next to source
- `internal/translator/*/` - Extensive translator tests
- `sdk/cliproxy/auth/*_test.go` - Auth system tests
- `test/` - Integration test utilities
