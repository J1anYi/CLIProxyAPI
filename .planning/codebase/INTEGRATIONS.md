# External Integrations

## APIs & Services

### AI Provider APIs
- **Google Gemini / Cloud Code** (`cloudcode-pa.googleapis.com`, `googleapis.com`)
  - Gemini CLI authentication and API proxy
  - OAuth2 with Google identity
  - Model discovery and inference

- **Anthropic Claude** (`claude.ai`, `api.anthropic.com`)
  - Claude OAuth authentication flow
  - Messages API proxy
  - Console: https://console.anthropic.com/

- **OpenAI Codex** (`auth.openai.com`)
  - OAuth2 device flow and web flow
  - PKCE-based authentication
  - Chat completions API proxy

- **Kimi (Moonshot AI)** (`auth.kimi.com`, `api.kimi.com`)
  - Device authorization grant flow (RFC 8628)
  - Coding API integration

- **Google Antigravity** (`cloudcode-pa.googleapis.com`)
  - Google Cloud Code Assist integration
  - Project onboarding flow
  - Multi-tier user support

### Vertex AI
- Google Vertex AI service account import
- Vertex-compatible API key management
- Model registry integration

## Authentication Providers

### OAuth2 Providers
- **Google OAuth2** (`accounts.google.com`, `oauth2.googleapis.com`)
  - Scopes: cloud-platform, userinfo.email, userinfo.profile
  - Used for Gemini, Antigravity authentication

- **Anthropic OAuth** (`claude.ai/oauth/authorize`, `api.anthropic.com/v1/oauth/token`)
  - Custom OAuth flow for Claude

- **OpenAI Auth** (`auth.openai.com`)
  - PKCE-enhanced authorization code flow
  - Device code flow support

- **Kimi OAuth** (`auth.kimi.com`)
  - Device authorization grant (RFC 8628)
  - Device ID tracking

## Database/Storage

### PostgreSQL
- Connection via pgx driver (v5)
- Schema-managed tables for config and auth storage
- Environment variable configuration (PGSTORE_DSN, PGSTORE_SCHEMA)
- Local workspace mirroring

### Object Storage (S3-Compatible)
- MinIO client for S3-compatible storage
- Bucket-based config/auth persistence
- Environment configuration (OBJECTSTORE_ENDPOINT, OBJECTSTORE_BUCKET, etc.)
- SSL/TLS support with path-style addressing

### Git-Based Storage
- Remote Git repository for config/auth sync
- Credentials via environment (GITSTORE_GIT_URL, GITSTORE_GIT_TOKEN)
- Branch management support

### Local Filesystem
- Default file-based token storage
- Auth directory with JSON files
- Hot-reload via fsnotify

## Third-party Services

### GitHub API
- Release checking for management panel updates
- URL: `https://api.github.com/repos/router-for-me/Cli-Proxy-API-Management-Center/releases/latest`
- Fallback CDN: `https://cpamc.router-for.me/`

### Google Cloud Services
- Cloud Resource Manager API (`cloudresourcemanager.googleapis.com`)
- Service Usage API (`serviceusage.googleapis.com`)
- Compute metadata service

## Webhooks & Callbacks

### OAuth Callback Endpoints
- `/anthropic/callback` - Claude OAuth callback
- `/codex/callback` - OpenAI Codex OAuth callback
- `/google/callback` - Google/Gemini OAuth callback
- `/antigravity/callback` - Antigravity OAuth callback

### WebSocket Relay
- `/v1/ws` - WebSocket endpoint for streaming responses
- Session management with heartbeat
- Message routing by request ID

## API Compatibility Layers

### OpenAI-Compatible Endpoints
- `POST /v1/chat/completions` - Chat completions
- `POST /v1/completions` - Legacy completions
- `POST /v1/messages` - Claude-style messages
- `GET /v1/models` - Model listing
- `POST /v1/responses` - OpenAI Responses API

### Gemini-Compatible Endpoints
- `GET/POST /v1beta/models/*` - Gemini API proxy
- `POST /v1internal:method` - Gemini CLI internal API

### Management API
- RESTful management endpoints under `/v0/management/`
- Configuration hot-reload
- Auth file management
- Usage statistics

## Proxy Configuration
- HTTP/HTTPS proxy support via PROXY_URL environment
- TLS fingerprint customization (utls transport)
- Custom certificate support
