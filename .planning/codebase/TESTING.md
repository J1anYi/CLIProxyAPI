# Testing Strategy

## Test Frameworks

- **testing** - Go standard library test framework
- **testify** - Not used; tests use standard library patterns
- **httptest** - HTTP server mocking via `net/http/httptest`
- **gin** - Web framework test mode: `gin.SetMode(gin.TestMode)`

## Test Organization

### File Structure
- Tests colocated with source code
- Naming: `<source>_test.go`
- Platform-specific: `<source>_windows_test.go`, `<source>_linux_test.go`

### Package Organization
```
internal/
├── api/handlers/management/
│   ├── api_tools.go
│   ├── api_tools_test.go
│   ├── auth_files_download_test.go
│   └── auth_files_download_windows_test.go
├── runtime/executor/
│   ├── claude_executor_test.go
│   ├── codex_executor_retry_test.go
│   └── ...
```

### Test Package
- Tests use same package name (not `_test` suffix)
- Access to unexported functions and types

## Coverage Areas

### Core Components Tested
1. **API Handlers** (`internal/api/handlers/management/`)
   - HTTP request/response handling
   - Authentication/authorization
   - File operations (download, delete, patch)
   - Configuration management

2. **Executors** (`internal/runtime/executor/`)
   - Request building and signing
   - Response handling and decoding
   - Retry logic with backoff
   - Streaming support
   - Token counting

3. **Translators** (`internal/translator/`)
   - Format conversion: Claude ↔ OpenAI ↔ Gemini
   - Request/response transformation
   - Streaming event translation

4. **Authentication** (`internal/auth/`)
   - OAuth flows (Claude, Codex, Gemini)
   - Token refresh and storage
   - Error handling

5. **Configuration** (`internal/config/`)
   - YAML loading and validation
   - Default value application
   - Sanitization and normalization

6. **Caching** (`internal/cache/`)
   - Signature cache validation
   - TTL enforcement

7. **Watchers** (`internal/watcher/`)
   - File change detection
   - Configuration diff
   - Hot reload

## Test Types

### Unit Tests
- Most common test type
- Single function/method focus
- Mock external dependencies
- Examples:
  - `TestAPICallTransportDirectBypassesGlobalProxy`
  - `TestParseCodexRetryAfter`
  - `TestApplyClaudeToolPrefix`

### Integration Tests
- Test component interactions
- Use `httptest.Server` for HTTP mocking
- Examples:
  - `TestClaudeExecutor_ReusesUserIDAcrossModelsWhenCacheEnabled`
  - `TestClaudeExecutor_ExecuteStream_SetsIdentityAcceptEncoding`

### Table-Driven Tests
- Used for testing multiple scenarios
- Slice of test cases with expected outcomes
```go
cases := []struct {
    name      string
    auth      *coreauth.Auth
    wantProxy string
}{
    {name: "gemini", auth: ..., wantProxy: "..."},
    {name: "claude", auth: ..., wantProxy: "..."},
}
for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        // test logic
    })
}
```

### Concurrency Tests
- Race condition detection
- `sync.WaitGroup` and channels for coordination
- Example: `TestResolveClaudeDeviceProfile_RechecksCacheBeforeStoringCandidate`

### Platform-Specific Tests
- Build constraints: `//go:build windows`
- Separate files for OS-specific behavior
- Example: `auth_files_download_windows_test.go`

## Test Patterns

### Test Helpers
```go
func newClaudeHeaderTestRequest(t *testing.T, incoming http.Header) *http.Request {
    t.Helper()  // Marks as helper for better error reporting
    // setup code
}

func assertClaudeFingerprint(t *testing.T, headers http.Header, ...) {
    t.Helper()
    // assertion code with t.Fatalf on failure
}
```

### Parallel Execution
```go
func TestSomething(t *testing.T) {
    t.Parallel()  // Run concurrently with other parallel tests
}
```

### Cleanup
```go
func TestWithCleanup(t *testing.T) {
    server := httptest.NewServer(handler)
    t.Cleanup(func() { server.Close() })
    // test code
}
```

### Mock HTTP Servers
```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // inspect request, write response
}))
defer server.Close()
```

### Test Data
- Inline JSON/byte slices for test cases
- `[]byte(`{"key":"value"}`)` patterns
- `gjson.GetBytes()` and `sjson.SetBytes()` for manipulation

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Specific Package
```bash
go test ./internal/runtime/executor/...
```

### Run Specific Test
```bash
go test -run TestFunctionName ./path/to/package
```

### Run with Verbose Output
```bash
go test -v ./...
```

### Run with Race Detection
```bash
go test -race ./...
```

### Run with Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Parallel Tests
```bash
go test -parallel 4 ./...
```

## Test Assertions

### Standard Library Style
- Use `t.Fatalf()` for fatal failures with formatted message
- Use `t.Errorf()` for non-fatal failures
- Use `t.Skip()` for skipping tests

### Common Patterns
```go
if got != want {
    t.Fatalf("value = %q, want %q", got, want)
}

if err == nil {
    t.Fatal("expected error, got nil")
}

if !strings.Contains(err.Error(), "expected substring") {
    t.Fatalf("error = %v, want containing %q", err, "expected substring")
}
```

### Type Assertions in Tests
```go
httpTransport, ok := transport.(*http.Transport)
if !ok {
    t.Fatalf("transport type = %T, want *http.Transport", transport)
}
```

## Coverage Metrics

### Running Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Coverage Scope
- Unit tests cover individual functions
- Integration tests cover component interactions
- Critical paths: authentication, request execution, translation

### Coverage Goals
- Focus on correctness over percentage
- Cover edge cases and error paths
- Test concurrent access patterns

## Test Data Management

### Inline Test Data
- JSON payloads embedded in test functions
- Reusable via helper functions

### Fixture Files
- Located in testdata/ directories when needed
- Loaded via `os.ReadFile` or `embed`

### Mock Data Builders
- Factory functions for test objects
```go
func newTestAuth(apiKey, baseURL string) *cliproxyauth.Auth {
    return &cliproxyauth.Auth{
        Attributes: map[string]string{
            "api_key":  apiKey,
            "base_url": baseURL,
        },
    }
}
```
