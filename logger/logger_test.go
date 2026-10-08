package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type captureLogger struct{ buf bytes.Buffer }

func (c *captureLogger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(&c.buf, format, v...)
}

const (
	testToken   = "eyJhbGciOiJSUzI1NiJ9.secret-access-token"
	testSession = "SESSION-ID-0123456789"
)

func newTestRequest(t *testing.T, body string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "https://api.zsapi.net/zia/api/v1/urlFilteringRules", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Cookie", "JSESSIONID="+testSession+"; other=value")
	req.Header.Set("JSessionID", testSession)
	req.Header.Set("auth-token", testToken)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertNoSecrets(t *testing.T, out string) {
	t.Helper()
	for _, secret := range []string{testToken, testSession} {
		if strings.Contains(out, secret) {
			t.Fatalf("logged output contains a secret %q:\n%s", secret, out)
		}
	}
}

func TestLogRequest_MasksSensitiveHeaders(t *testing.T) {
	l := &captureLogger{}
	req := newTestRequest(t, `{"name":"rule1"}`)

	LogRequest(l, req, "req-1", nil, true)
	out := l.buf.String()

	assertNoSecrets(t, out)
	for _, want := range []string{
		"Authorization: Bearer ********",
		"Cookie: JSESSIONID=********; other=********",
		"Jsessionid: ********",
		"Auth-Token: ********",
		"Content-Type: application/json",
		`{"name":"rule1"}`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("logged output missing %q:\n%s", want, out)
		}
	}

	// The request itself must be unchanged.
	if got := req.Header.Get("Authorization"); got != "Bearer "+testToken {
		t.Errorf("request Authorization header changed: %q", got)
	}
	body, _ := io.ReadAll(req.Body)
	if string(body) != `{"name":"rule1"}` {
		t.Errorf("request body consumed or changed: %q", string(body))
	}
}

func TestLogRequestSensitive_MasksSensitiveHeaders(t *testing.T) {
	l := &captureLogger{}
	req := newTestRequest(t, `{"password":"p@ss"}`)

	LogRequestSensitive(l, req, "req-2", []string{"p@ss"})
	out := l.buf.String()

	assertNoSecrets(t, out)
	if strings.Contains(out, "p@ss") {
		t.Errorf("explicit sensitive content was not masked:\n%s", out)
	}
}

func TestLogResponse_MasksSetCookie(t *testing.T) {
	l := &captureLogger{}
	u, _ := url.Parse("https://api.zsapi.net/zia/api/v1/authenticatedSession")
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Set-Cookie":   []string{"JSESSIONID=" + testSession + "; Path=/; Secure; HttpOnly"},
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(`{"authType":"API"}`)),
		Request: &http.Request{Method: http.MethodGet, URL: u},
	}

	LogResponse(l, resp, time.Now(), "req-3")
	out := l.buf.String()

	assertNoSecrets(t, out)
	if !strings.Contains(out, "Set-Cookie: JSESSIONID=********; Path=/; Secure; HttpOnly") {
		t.Errorf("Set-Cookie not masked as expected:\n%s", out)
	}
	if !strings.Contains(out, `{"authType":"API"}`) {
		t.Errorf("response body missing from log:\n%s", out)
	}
}

func TestMaskSensitiveHeaders_LeavesBodyUntouched(t *testing.T) {
	dump := "POST /x HTTP/1.1\r\nAuthorization: Bearer abc\r\n\r\nAuthorization: Bearer not-a-header"
	got := string(maskSensitiveHeaders([]byte(dump)))
	want := "POST /x HTTP/1.1\r\nAuthorization: Bearer ********\r\n\r\nAuthorization: Bearer not-a-header"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
