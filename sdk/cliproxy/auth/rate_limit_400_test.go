package auth

import (
	"net/http"
	"testing"
)

func TestIsRateLimitDisguisedAs400(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "ModelArts.81101 error",
			err:      &Error{Message: `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81101\",\"message\":\"Too many requests, the rate limit is 25000000 tokens per minute.\",\"param\":null,\"type\":\"TooManyRequests\"}}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "ModelArts.81011 sensitive info error",
			err:      &Error{Message: `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81011\",\"message\":\"Input text May contain sensitive information, please try again.\",\"param\":null,\"type\":\"Forbidden\"}}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "TooManyRequests type in error",
			err:      &Error{Message: `{"error":{"type":"TooManyRequests","message":"Rate limit exceeded"}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "TooManyRequests type with single quotes",
			err:      &Error{Message: `{'error':{'type':'TooManyRequests','message':'Rate limit exceeded'}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "too many requests message",
			err:      &Error{Message: "Error: too many requests, please try again later", HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "rate limit exceeded message",
			err:      &Error{Message: "Rate limit exceeded, please wait before retrying", HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "rate_limit exceeded message",
			err:      &Error{Message: "rate_limit exceeded for user", HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "invalid_request_error should not be rate limit",
			err:      &Error{Message: `{"error":{"type":"invalid_request_error","message":"Invalid parameter"}}`, HTTPStatus: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "random 400 error",
			err:      &Error{Message: "Bad request: missing required field", HTTPStatus: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "rate limit without exceeded",
			err:      &Error{Message: "The rate limit is 100 requests per minute", HTTPStatus: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "Decode server is overloaded",
			err:      &Error{Message: `{"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}","code":400,"message":"模型服务调用失败"}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "decode server is overloaded lowercase",
			err:      &Error{Message: "decode server is overloaded, please retry", HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "Decode server overloaded partial match",
			err:      &Error{Message: "Error: Decode server overloaded", HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRateLimitDisguisedAs400(tt.err)
			if got != tt.expected {
				t.Errorf("isRateLimitDisguisedAs400() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRequestInvalidErrorWithRateLimit400(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ModelArts.81101 should be retryable (not invalid)",
			err:      &Error{Message: `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81101\",\"message\":\"Too many requests, the rate limit is 25000000 tokens per minute.\",\"param\":null,\"type\":\"TooManyRequests\"}}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: false, // Should return false so it can be retried
		},
		{
			name:     "invalid_request_error should be invalid",
			err:      &Error{Message: `{"error":{"type":"invalid_request_error","message":"Invalid parameter"}}`, HTTPStatus: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "422 should be invalid",
			err:      &Error{Message: "Unprocessable entity", HTTPStatus: http.StatusUnprocessableEntity},
			expected: true,
		},
		{
			name:     "429 should not be invalid",
			err:      &Error{Message: "Too many requests", HTTPStatus: http.StatusTooManyRequests},
			expected: false,
		},
		{
			name:     "500 should not be invalid",
			err:      &Error{Message: "Internal server error", HTTPStatus: http.StatusInternalServerError},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRequestInvalidError(tt.err)
			if got != tt.expected {
				t.Errorf("isRequestInvalidError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestError matches the real-world ModelArts error format
func TestModelArts81101ErrorFormat(t *testing.T) {
	// This is the exact error format from the user's log
	realWorldError := `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81101\",\"message\":\"Too many requests, the rate limit is 25000000 tokens per minute.\",\"param\":null,\"type\":\"TooManyRequests\"},\"error_code\":\"ModelArts.81101\",\"error_msg\":\"Too many requests, the rate limit is 25000000 tokens per minute.\",\"span_id\":\"c3ec6e9dca2be8560fffd57d9e5176cb\"}","code":400,"message":"模型服务调用失败","status":"FAILED_RESPONSE"},"requestId":"a904d4a7c6fb4124dcab913d47262e35-aHdke","result":null}`

	err := &Error{
		Message:    realWorldError,
		HTTPStatus: http.StatusBadRequest,
	}

	// Should detect as rate limit
	if !isRateLimitDisguisedAs400(err) {
		t.Error("ModelArts.81101 error should be detected as rate limit")
	}

	// Should NOT be treated as invalid request (so it can be retried)
	if isRequestInvalidError(err) {
		t.Error("ModelArts.81101 error should NOT be treated as invalid request")
	}
}

// TestDecodeServerOverloadedErrorFormat tests the exact error format from the user's log
func TestDecodeServerOverloadedErrorFormat(t *testing.T) {
	// This is the exact error format from the user's log
	realWorldError := `{"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}","code":400,"message":"模型服务调用失败","status":"FAILED_RESPONSE"},"requestId":"8b53d485fd30ba9f00bc7bb5ba5eb7bf-eUNmU","result":null}`

	err := &Error{
		Message:    realWorldError,
		HTTPStatus: http.StatusBadRequest,
	}

	// Should detect as rate limit
	if !isRateLimitDisguisedAs400(err) {
		t.Error("Decode server overloaded error should be detected as rate limit")
	}

	// Should NOT be treated as invalid request (so it can be retried)
	if isRequestInvalidError(err) {
		t.Error("Decode server overloaded error should NOT be treated as invalid request")
	}
}

// TestIsRequestInvalidErrorWithDecodeServerOverloaded tests that DecodeServerOverloaded is retryable
func TestIsRequestInvalidErrorWithDecodeServerOverloaded(t *testing.T) {
	err := &Error{
		Message:    `{"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}","code":400,"message":"模型服务调用失败"}}`,
		HTTPStatus: http.StatusBadRequest,
	}

	// Should return false so it can be retried
	if isRequestInvalidError(err) {
		t.Error("Decode server overloaded error should NOT be treated as invalid request")
	}
}

// TestDetectRateLimitErrorType tests the error type discrimination
func TestDetectRateLimitErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "ModelArts.81101 error",
			err:      &Error{Message: `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81101\",\"message\":\"Too many requests\"}}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: "ModelArts81101",
		},
		{
			name:     "ModelArts.81011 sensitive info error",
			err:      &Error{Message: `{"error":{"cause":"{\"error\":{\"code\":\"ModelArts.81011\",\"message\":\"Input text May contain sensitive information, please try again.\"}}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: "ModelArts81011",
		},
		{
			name:     "Decode server is overloaded",
			err:      &Error{Message: `{"error":{"cause":"{\"detail\":\"Decode server is overloaded\"}"}}`, HTTPStatus: http.StatusBadRequest},
			expected: "DecodeServerOverloaded",
		},
		{
			name:     "Decode server overloaded partial",
			err:      &Error{Message: "Decode server overloaded", HTTPStatus: http.StatusBadRequest},
			expected: "DecodeServerOverloaded",
		},
		{
			name:     "TooManyRequests type without ModelArts code",
			err:      &Error{Message: `{"error":{"type":"TooManyRequests","message":"Rate limit"}}`, HTTPStatus: http.StatusBadRequest},
			expected: "",
		},
		{
			name:     "Random error",
			err:      &Error{Message: "Something went wrong", HTTPStatus: http.StatusBadRequest},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectRateLimitErrorType(tt.err)
			if got != tt.expected {
				t.Errorf("detectRateLimitErrorType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
