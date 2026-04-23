# Technology Stack

## Runtime & Language
- **Language:** Go 1.26.0
- **Module:** github.com/router-for-me/CLIProxyAPI/v6

## Core Frameworks
- **Gin Web Framework** v1.10.1 - HTTP server and routing
- **Gorilla WebSocket** v1.5.3 - WebSocket support for real-time communication
- **golang.org/x/oauth2** v0.30.0 - OAuth2 client implementation

## Key Libraries

### Web & HTTP
- github.com/gorilla/websocket v1.5.3 - WebSocket protocol implementation
- golang.org/x/net v0.47.0 - Extended networking support

### Data Storage
- github.com/jackc/pgx/v5 v5.7.6 - PostgreSQL driver and connection pool
- github.com/minio/minio-go/v7 v7.0.66 - S3-compatible object storage client
- github.com/go-git/go-git/v6 v6.0.0-20251009132922-75a182125145 - Git repository operations

### Configuration & Data Formats
- gopkg.in/yaml.v3 v3.0.1 - YAML parsing
- github.com/tidwall/gjson v1.18.0 - Fast JSON path queries
- github.com/tidwall/sjson v1.2.5 - JSON modification

### Logging
- github.com/sirupsen/logrus v1.9.3 - Structured logging
- gopkg.in/natefinch/lumberjack.v2 v2.2.1 - Log rotation

### Authentication & Security
- golang.org/x/crypto v0.45.0 - Cryptographic primitives
- github.com/refraction-networking/utls v1.8.2 - TLS fingerprinting for OAuth flows
- github.com/google/uuid v1.6.0 - UUID generation

### UI & Terminal
- github.com/charmbracelet/bubbletea v1.3.10 - TUI framework
- github.com/charmbracelet/bubbles v1.0.0 - TUI components
- github.com/charmbracelet/lipgloss v1.1.0 - Terminal styling
- github.com/atotto/clipboard v0.1.4 - Clipboard access

### Compression & Encoding
- github.com/klauspost/compress v1.17.4 - Compression algorithms
- github.com/andybalholm/brotli v1.0.6 - Brotli compression

### Utilities
- github.com/joho/godotenv v1.5.1 - .env file loading
- github.com/fsnotify/fsnotify v1.9.0 - File system notifications
- github.com/tiktoken-go/tokenizer v0.7.0 - Token counting
- github.com/pierrec/xxHash v0.1.5 - Fast hashing
- golang.org/x/sync v0.18.0 - Concurrency primitives
- github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 - Browser URL opening

## Build Tools
- Go modules (go.mod)
- Standard Go toolchain

## Development Tools
- **Testing:** Go standard testing package with extensive test coverage (100+ test files)
- **CLI:** Cobra-like flag-based command-line interface

## Architecture Patterns
- **Modular Design:** internal/ package structure with clear separation of concerns
- **SDK Package:** External-facing SDK under sdk/ directory
- **Handler Pattern:** API handlers for OpenAI, Gemini, Claude compatibility
- **Translator Pattern:** Request/response translators between different API formats
- **Store Abstraction:** Pluggable storage backends (PostgreSQL, Git, S3/MinIO, filesystem)
