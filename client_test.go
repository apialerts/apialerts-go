package apialerts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- Test server helpers ---

func serverWithResponse(statusCode int, body map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func slowServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

func captureServer(t *testing.T, statusCode int, body map[string]any) (*httptest.Server, *http.Request, chan *http.Request) {
	captured := make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured <- r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
	return server, nil, captured
}

// --- Validation tests ---

func TestSendMissingMessage(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync("http://unused", "test_api_key", Event{})
	if err == nil || err.Error() != "message is required" {
		t.Errorf("expected 'message is required', got: %v", err)
	}
}

func TestSendMissingApiKey(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync("http://unused", "", Event{Message: "hello"})
	if err == nil {
		t.Error("expected error for missing API key, got nil")
	}
}

// --- HTTP response tests ---

func TestSend200Success(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	Configure("test_api_key")

	result, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Workspace != "Acme Corp" {
		t.Errorf("expected workspace 'Acme Corp', got '%s'", result.Workspace)
	}
	if result.Channel != "general" {
		t.Errorf("expected channel 'general', got '%s'", result.Channel)
	}
}

func TestSend200WithWarnings(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
		"errors":    []string{"unknown field: foo", "tag limit reached"},
	})
	defer server.Close()

	Configure("test_api_key")

	result, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(result.Warnings))
	}
	if result.Warnings[0] != "unknown field: foo" {
		t.Errorf("unexpected warning: %s", result.Warnings[0])
	}
}

func TestSend200EmptyWarnings(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	Configure("test_api_key")

	result, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(result.Warnings))
	}
}

func TestSend400BadRequest(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusBadRequest, nil)
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil || err.Error() != "bad request" {
		t.Errorf("expected 'bad request' error, got: %v", err)
	}
}

func TestSend401Unauthorized(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusUnauthorized, nil)
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil || !strings.Contains(err.Error(), "unauthorized") {
		t.Errorf("expected 'unauthorized' error, got: %v", err)
	}
}

func TestSend403Forbidden(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusForbidden, nil)
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil || err.Error() != "forbidden" {
		t.Errorf("expected 'forbidden' error, got: %v", err)
	}
}

func TestSend429RateLimit(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusTooManyRequests, nil)
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil || err.Error() != "rate limit exceeded" {
		t.Errorf("expected 'rate limit exceeded' error, got: %v", err)
	}
}

func TestSend500UnexpectedError(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusInternalServerError, nil)
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil || err.Error() != "unexpected error" {
		t.Errorf("expected 'unexpected error', got: %v", err)
	}
}

func TestSendNetworkError(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync("http://127.0.0.1:1", "test_api_key", Event{Message: "hello"})
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

func TestSendInvalidJsonResponse(t *testing.T) {
	resetInstance()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	Configure("test_api_key")

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil {
		t.Error("expected JSON parse error, got nil")
	}
}

func TestSendTimeout(t *testing.T) {
	resetInstance()
	server := slowServer(2 * time.Second)
	defer server.Close()

	ConfigureWithConfig("test_api_key", Config{Timeout: 100 * time.Millisecond})

	_, err := instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

// --- Header tests ---

func TestRequestHeaders(t *testing.T) {
	resetInstance()
	server, _, captured := captureServer(t, http.StatusOK, map[string]any{
		"workspace": "test",
		"channel":   "test",
	})
	defer server.Close()

	Configure("my_api_key")

	instance.sendToUrlWithApiKeyAsync(server.URL, "my_api_key", Event{Message: "hello"})

	req := <-captured
	if req.Header.Get("Authorization") != "Bearer my_api_key" {
		t.Errorf("expected Authorization 'Bearer my_api_key', got '%s'", req.Header.Get("Authorization"))
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("X-Integration") != IntegrationName {
		t.Errorf("expected X-Integration '%s', got '%s'", IntegrationName, req.Header.Get("X-Integration"))
	}
	if req.Header.Get("X-Version") != IntegrationVersion {
		t.Errorf("expected X-Version '%s', got '%s'", IntegrationVersion, req.Header.Get("X-Version"))
	}
}

func TestSetOverridesHeaders(t *testing.T) {
	resetInstance()
	server, _, captured := captureServer(t, http.StatusOK, map[string]any{
		"workspace": "test",
		"channel":   "test",
	})
	defer server.Close()

	Configure("test_api_key")
	SetOverrides("cli", "1.2.0", "")

	instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})

	req := <-captured
	if req.Header.Get("X-Integration") != "cli" {
		t.Errorf("expected X-Integration 'cli', got '%s'", req.Header.Get("X-Integration"))
	}
	if req.Header.Get("X-Version") != "1.2.0" {
		t.Errorf("expected X-Version '1.2.0', got '%s'", req.Header.Get("X-Version"))
	}
}

func TestSetOverridesBaseURL(t *testing.T) {
	resetInstance()
	server := serverWithResponse(http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	Configure("test_api_key")
	SetOverrides("cli", "1.2.0", server.URL)

	result, err := SendAsync(Event{Message: "hello"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Workspace != "Acme Corp" {
		t.Errorf("expected workspace 'Acme Corp', got '%s'", result.Workspace)
	}
}

// --- Payload tests ---

func TestRequestPayload(t *testing.T) {
	resetInstance()
	var decoded map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&decoded)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"workspace": "test", "channel": "test"})
	}))
	defer server.Close()

	Configure("test_api_key")

	instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{
		Event:   "user.purchase",
		Title:   "New Sale",
		Message: "hello",
		Channel: "payments",
		Tags:    []string{"billing", "prod"},
		Link:    "https://example.com",
	})

	if decoded["message"] != "hello" {
		t.Errorf("expected message 'hello', got '%v'", decoded["message"])
	}
	if decoded["channel"] != "payments" {
		t.Errorf("expected channel 'payments', got '%v'", decoded["channel"])
	}
	if decoded["event"] != "user.purchase" {
		t.Errorf("expected event 'user.purchase', got '%v'", decoded["event"])
	}
	if decoded["title"] != "New Sale" {
		t.Errorf("expected title 'New Sale', got '%v'", decoded["title"])
	}
	if decoded["link"] != "https://example.com" {
		t.Errorf("expected link 'https://example.com', got '%v'", decoded["link"])
	}
	tags, ok := decoded["tags"].([]interface{})
	if !ok || len(tags) != 2 {
		t.Errorf("expected 2 tags, got '%v'", decoded["tags"])
	}
}

func TestRequestPayloadWithData(t *testing.T) {
	resetInstance()
	var decoded map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&decoded)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"workspace": "test", "channel": "test"})
	}))
	defer server.Close()

	Configure("test_api_key")

	instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{
		Message: "hello",
		Data:    map[string]any{"plan": "pro", "amount": 49.99},
	})

	data, ok := decoded["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in payload, got '%v'", decoded["data"])
	}
	if data["plan"] != "pro" {
		t.Errorf("expected data.plan 'pro', got '%v'", data["plan"])
	}
}

func TestRequestPayloadOmitsEmptyData(t *testing.T) {
	resetInstance()
	var decoded map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&decoded)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"workspace": "test", "channel": "test"})
	}))
	defer server.Close()

	Configure("test_api_key")

	instance.sendToUrlWithApiKeyAsync(server.URL, "test_api_key", Event{Message: "hello"})

	if _, exists := decoded["data"]; exists {
		t.Errorf("expected 'data' to be omitted when nil, but it was present")
	}
}
