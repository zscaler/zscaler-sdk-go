package logger

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type nopLogger struct{}

func (l *nopLogger) Printf(format string, v ...interface{}) {}

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
---[ ZSCALER SDK REQUEST ]-------------------------------
%s
---------------------------------------------------------`

	logRespMsg = `[DEBUG] Response "%s %s" details:
---[ ZSCALER SDK RESPONSE ]--------------------------------
%s
-------------------------------------------------------`
)

func WriteLog(logger Logger, format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func LogRequest(logger Logger, req *http.Request) {
	if logger != nil && req != nil {
		out, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			WriteLog(logger, logReqMsg, req.Method, req.URL, string(out))
		}
	}
}

func LogResponse(logger Logger, resp *http.Response) {
	if logger != nil && resp != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			WriteLog(logger, logRespMsg, resp.Request.Method, resp.Request.URL, string(out))
		}
	}
}
