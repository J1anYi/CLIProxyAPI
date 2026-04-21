package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/interfaces"
	coreauth "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/auth"
	sdkconfig "github.com/router-for-me/CLIProxyAPI/v6/sdk/config"
)

func TestWriteErrorResponse_AddonHeadersDisabledByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler := NewBaseAPIHandlers(nil, nil)
	handler.WriteErrorResponse(c, &interfaces.ErrorMessage{
		StatusCode: http.StatusTooManyRequests,
		Error:      errors.New("rate limit"),
		Addon: http.Header{
			"Retry-After":  {"30"},
			"X-Request-Id": {"req-1"},
		},
	})

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	if got := recorder.Header().Get("Retry-After"); got != "" {
		t.Fatalf("Retry-After should be empty when passthrough is disabled, got %q", got)
	}
	if got := recorder.Header().Get("X-Request-Id"); got != "" {
		t.Fatalf("X-Request-Id should be empty when passthrough is disabled, got %q", got)
	}
}

func TestWriteErrorResponse_AddonHeadersEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Writer.Header().Set("X-Request-Id", "old-value")

	handler := NewBaseAPIHandlers(&sdkconfig.SDKConfig{PassthroughHeaders: true}, nil)
	handler.WriteErrorResponse(c, &interfaces.ErrorMessage{
		StatusCode: http.StatusTooManyRequests,
		Error:      errors.New("rate limit"),
		Addon: http.Header{
			"Retry-After":  {"30"},
			"X-Request-Id": {"new-1", "new-2"},
		},
	})

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	if got := recorder.Header().Get("Retry-After"); got != "30" {
		t.Fatalf("Retry-After = %q, want %q", got, "30")
	}
	if got := recorder.Header().Values("X-Request-Id"); !reflect.DeepEqual(got, []string{"new-1", "new-2"}) {
		t.Fatalf("X-Request-Id = %#v, want %#v", got, []string{"new-1", "new-2"})
	}
}

func TestEnrichAuthSelectionError_DefaultsTo503WithContext(t *testing.T) {
	in := &coreauth.Error{Code: "auth_not_found", Message: "no auth available"}
	out := enrichAuthSelectionError(in, []string{"claude"}, "claude-sonnet-4-6")

	var got *coreauth.Error
	if !errors.As(out, &got) || got == nil {
		t.Fatalf("expected coreauth.Error, got %T", out)
	}
	if got.StatusCode() != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", got.StatusCode(), http.StatusServiceUnavailable)
	}
	if !strings.Contains(got.Message, "providers=claude") {
		t.Fatalf("message missing provider context: %q", got.Message)
	}
	if !strings.Contains(got.Message, "model=claude-sonnet-4-6") {
		t.Fatalf("message missing model context: %q", got.Message)
	}
	if !strings.Contains(got.Message, "/v0/management/auth-files") {
		t.Fatalf("message missing management hint: %q", got.Message)
	}
}

func TestEnrichAuthSelectionError_PreservesExplicitStatus(t *testing.T) {
	in := &coreauth.Error{Code: "auth_unavailable", Message: "no auth available", HTTPStatus: http.StatusTooManyRequests}
	out := enrichAuthSelectionError(in, []string{"gemini"}, "gemini-2.5-pro")

	var got *coreauth.Error
	if !errors.As(out, &got) || got == nil {
		t.Fatalf("expected coreauth.Error, got %T", out)
	}
	if got.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", got.StatusCode(), http.StatusTooManyRequests)
	}
}

func TestEnrichAuthSelectionError_IgnoresOtherErrors(t *testing.T) {
	in := errors.New("boom")
	out := enrichAuthSelectionError(in, []string{"claude"}, "claude-sonnet-4-6")
	if out != in {
		t.Fatalf("expected original error to be returned unchanged")
	}
}

func TestIsContextLengthError(t *testing.T) {
	tests := []struct {
		name     string
		errText  string
		expected bool
	}{
		{
			name:     "exceeds maximum context length",
			errText:  "Requested token count exceeds the model's maximum context length of 202752 tokens",
			expected: true,
		},
		{
			name:     "requested token count exceeds",
			errText:  "Requested token count exceeds the model's maximum context length",
			expected: true,
		},
		{
			name:     "case insensitive",
			errText:  "EXCEEDS THE MODEL'S MAXIMUM CONTEXT LENGTH",
			expected: true,
		},
		{
			name:     "other error",
			errText:  "rate limit exceeded",
			expected: false,
		},
		{
			name:     "empty string",
			errText:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isContextLengthError(tt.errText); got != tt.expected {
				t.Errorf("isContextLengthError(%q) = %v, want %v", tt.errText, got, tt.expected)
			}
		})
	}
}

func TestExtractTokenCounts(t *testing.T) {
	tests := []struct {
		name           string
		errText        string
		wantInput      int
		wantMax        int
		wantFound      bool
	}{
		{
			name:      "full Anthropic error format",
			errText:   "Requested token count exceeds the model's maximum context length of 202752 tokens. You requested a total of 204909 tokens: 172909 tokens from the input messages and 32000 tokens for the completion.",
			wantInput: 172909,
			wantMax:   202752,
			wantFound: true,
		},
		{
			name:      "only context length",
			errText:   "exceeds the model's maximum context length of 4096 tokens",
			wantInput: 0,
			wantMax:   4096,
			wantFound: false,
		},
		{
			name:      "total tokens requested",
			errText:   "You requested a total of 53428 tokens, max 4096 tokens",
			wantInput: 53428,
			wantMax:   4096,
			wantFound: true,
		},
		{
			name:      "no token info",
			errText:   "some other error message",
			wantInput: 0,
			wantMax:   0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputTokens, maxTokens, found := extractTokenCounts(tt.errText)
			if inputTokens != tt.wantInput {
				t.Errorf("inputTokens = %d, want %d", inputTokens, tt.wantInput)
			}
			if maxTokens != tt.wantMax {
				t.Errorf("maxTokens = %d, want %d", maxTokens, tt.wantMax)
			}
			if found != tt.wantFound {
				t.Errorf("found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

func TestBuildContextLengthExceededError(t *testing.T) {
	data := buildContextLengthExceededError(172909, 202752)

	var resp contextLengthExceededResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Type != "error" {
		t.Errorf("type = %q, want %q", resp.Type, "error")
	}
	if resp.Error.Type != "invalid_request_error" {
		t.Errorf("error.type = %q, want %q", resp.Error.Type, "invalid_request_error")
	}
	if !strings.Contains(resp.Error.Message, "input_length and max_tokens exceed context limit") {
		t.Errorf("error.message should contain 'input_length and max_tokens exceed context limit', got %q", resp.Error.Message)
	}
	if !strings.Contains(resp.Error.Message, "172909") {
		t.Errorf("error.message should contain input token count, got %q", resp.Error.Message)
	}
}

func TestBuildErrorResponseBody_ContextLengthError(t *testing.T) {
	errText := `{"error": {"object": "error", "message": "Requested token count exceeds the model's maximum context length of 202752 tokens. You requested a total of 204909 tokens: 172909 tokens from the input messages and 32000 tokens for the completion.", "type": "BadRequestError", "param": null, "code": 400}}`

	data := BuildErrorResponseBody(http.StatusBadRequest, errText)

	var resp contextLengthExceededResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Type != "error" {
		t.Errorf("type = %q, want %q", resp.Type, "error")
	}
	if resp.Error.Type != "invalid_request_error" {
		t.Errorf("error.type = %q, want %q", resp.Error.Type, "invalid_request_error")
	}
}

func TestBuildErrorResponseBody_ContextLengthErrorNoTokens(t *testing.T) {
	errText := "exceeds the model's maximum context length"

	data := BuildErrorResponseBody(http.StatusBadRequest, errText)

	var resp contextLengthExceededResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Type != "error" {
		t.Errorf("type = %q, want %q", resp.Type, "error")
	}
	if resp.Error.Type != "invalid_request_error" {
		t.Errorf("error.type = %q, want %q", resp.Error.Type, "invalid_request_error")
	}
}

func TestBuildErrorResponseBody_NonContextLengthError(t *testing.T) {
	errText := "rate limit exceeded"

	data := BuildErrorResponseBody(http.StatusTooManyRequests, errText)

	var resp ErrorResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error.Type != "rate_limit_error" {
		t.Errorf("error.type = %q, want %q", resp.Error.Type, "rate_limit_error")
	}
}
