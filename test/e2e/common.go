// Package e2e содержит e2e тесты.
package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
)

const defaultBaseURL = "http://localhost:18080"

// baseURL возвращает базовый URL сервиса для e2e-тестов.
func baseURL() string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	return defaultBaseURL
}

// doRequest — helper, который отправляет HTTP-запрос и проверяет статус.
func doRequest(
	t *testing.T,
	method, path string,
	body any,
	expectedStatus int,
	out any,
) {
	t.Helper()

	url := baseURL() + path

	var reqBody *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(data)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != expectedStatus {
		var errResp httperr.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		t.Fatalf(
			"%s %s: unexpected status %d (want %d), error: %+v",
			method, path, resp.StatusCode, expectedStatus, errResp,
		)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("decode response for %s %s: %v", method, path, err)
		}
	}
}
