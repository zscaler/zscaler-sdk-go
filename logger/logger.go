package logger

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type nopLogger struct{}

func (l *nopLogger) Printf(format string, v ...interface{}) {}

func NewNopLogger() Logger {
	return &nopLogger{}
}

type defaultLogger struct {
	logger  *log.Logger
	Verbose bool
}

func (l *defaultLogger) Printf(format string, v ...interface{}) {
	trimedF := strings.TrimSpace(format)
	if (strings.HasPrefix(trimedF, "[DEBUG]") || strings.HasPrefix(trimedF, "[TRACE]")) && !l.Verbose {
		return
	}

	l.logger.Printf(format, v...)
}

func GetDefaultLogger(loggerPrefix string) Logger {
	loggingEnabled, _ := strconv.ParseBool(os.Getenv("ZSCALER_SDK_LOG"))
	if !loggingEnabled {
		return &nopLogger{}
	}
	verbose, _ := strconv.ParseBool(os.Getenv("ZSCALER_SDK_VERBOSE"))
	return &defaultLogger{
		logger:  log.New(os.Stdout, loggerPrefix, log.LstdFlags|log.Lshortfile),
		Verbose: verbose,
	}
}

const (
	logReqMsg = `[DEBUG] Request "%s %s" details:
---[ ZSCALER SDK REQUEST | ID:%s ]-------------------------------
%s
---------------------------------------------------------`

	logRespMsg = `[DEBUG] Response "%s %s" details:
---[ ZSCALER SDK RESPONSE | ID:%s | Duration:%s ]--------------------------------
%s
-------------------------------------------------------`
)

const maskedValue = "********"

// sensitiveHeaders are headers whose values carry credentials or session
// identifiers. Their values are masked in logged request and response dumps.
var sensitiveHeaders = map[string]bool{
	"authorization":       true,
	"proxy-authorization": true,
	"cookie":              true,
	"set-cookie":          true,
	"jsessionid":          true,
	"auth-token":          true,
	"x-api-key":           true,
}

// maskSensitiveHeaders masks the values of sensitive headers in an HTTP dump
// produced by httputil. Only the logged text is changed; the request or
// response itself is never modified. The header name, the Authorization scheme
// and cookie names are kept so the log remains useful.
func maskSensitiveHeaders(dump []byte) []byte {
	s := string(dump)
	headerEnd := strings.Index(s, "\r\n\r\n")
	if headerEnd < 0 {
		headerEnd = len(s)
	}
	lines := strings.Split(s[:headerEnd], "\r\n")
	for i := 1; i < len(lines); i++ { // line 0 is the request or status line
		name, value, ok := strings.Cut(lines[i], ":")
		key := strings.ToLower(strings.TrimSpace(name))
		if !ok || !sensitiveHeaders[key] {
			continue
		}
		lines[i] = name + ": " + maskHeaderValue(key, strings.TrimSpace(value))
	}
	return []byte(strings.Join(lines, "\r\n") + s[headerEnd:])
}

func maskHeaderValue(name, value string) string {
	switch name {
	case "authorization", "proxy-authorization":
		// Keep the scheme (e.g. "Bearer") and mask the credential.
		if scheme, _, ok := strings.Cut(value, " "); ok {
			return scheme + " " + maskedValue
		}
		return maskedValue
	case "cookie":
		// "a=1; b=2" -> "a=********; b=********"
		pairs := strings.Split(value, ";")
		for i, p := range pairs {
			if k, _, ok := strings.Cut(p, "="); ok {
				pairs[i] = k + "=" + maskedValue
			} else {
				pairs[i] = maskedValue
			}
		}
		return strings.Join(pairs, ";")
	case "set-cookie":
		// "JSESSIONID=abc; Path=/; Secure" -> "JSESSIONID=********; Path=/; Secure"
		first, attrs, hasAttrs := strings.Cut(value, ";")
		if k, _, ok := strings.Cut(first, "="); ok {
			first = k + "=" + maskedValue
		} else {
			first = maskedValue
		}
		if hasAttrs {
			return first + ";" + attrs
		}
		return first
	default:
		return maskedValue
	}
}

func WriteLog(logger Logger, format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func LogRequestSensitive(logger Logger, req *http.Request, reqID string, sensitiveContent []string) {
	if logger != nil && req != nil {
		out, err := httputil.DumpRequestOut(req, true)
		out = maskSensitiveHeaders(out)
		for _, s := range sensitiveContent {
			out = []byte(strings.ReplaceAll(string(out), s, maskedValue))
		}
		if err == nil {
			WriteLog(logger, logReqMsg, req.Method, req.URL, reqID, string(out))
		}
	}
}

func LogRequest(logger Logger, req *http.Request, reqID string, otherHeaderParams map[string]string, body bool) {
	if logger != nil && req != nil {
		l, ok := logger.(*defaultLogger)
		if ok && l.Verbose {
			for k, v := range otherHeaderParams {
				req.Header.Add(k, v)
			}
		}
		out, err := httputil.DumpRequestOut(req, body)
		if err == nil {
			WriteLog(logger, logReqMsg, req.Method, req.URL, reqID, string(maskSensitiveHeaders(out)))
		}
	}
}

func LogResponse(logger Logger, resp *http.Response, start time.Time, reqID string) {
	if logger != nil && resp != nil {
		// Dump the entire response
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			WriteLog(logger, logRespMsg, resp.Request.Method, resp.Request.URL, reqID, time.Since(start).String(), string(maskSensitiveHeaders(out)))
		} else {
			WriteLog(logger, logRespMsg, resp.Request.Method, resp.Request.URL, reqID, time.Since(start).String(), "Got error:"+err.Error())
		}
	}
}
