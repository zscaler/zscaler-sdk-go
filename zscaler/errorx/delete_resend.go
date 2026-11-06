package errorx

import (
	"context"
	"net/http"
	"sync/atomic"
)

// A DELETE can be applied by the API although the client does not see a
// success: the API reported an edit lock or org barrier conflict for it
// (observed while a configuration activation was running), a gateway returned
// a 5xx while the backend completed it, or the connection failed before the
// response arrived. The clients retry those outcomes, and the next attempt
// then gets 404 because the object is already gone. The object is deleted,
// which is what the caller asked for, so that 404 must not be reported as a
// failure.
//
// A client handles this as follows:
//   - TrackDeleteAttempts marks the context of a DELETE request once, before
//     the first attempt, so every attempt of that request (including a
//     request rebuilt for a resend) shares one record;
//   - TrackAttempts wraps the transport, so every attempt is recorded, whether
//     it is a retry of the retryable HTTP client or a resend of the client's
//     own request loop;
//   - DeleteAlreadyApplied tells the request layer to report the final 404 as
//     a success.
//
// A 404 is still an error when no earlier attempt can have been applied: on
// the first attempt, or after attempts that the API turned away without
// processing them (429 rate limit, 401/403 session errors). The object did
// not exist.

type deleteAttemptsKey struct{}

type deleteAttempts struct {
	// mayHaveApplied is set once an attempt ends without a response, or with
	// a response that does not prove the DELETE was not applied.
	mayHaveApplied atomic.Bool
}

// TrackDeleteAttempts returns ctx prepared to record the attempts of a DELETE
// request. For any other method, or when ctx is already tracked, it returns
// ctx unchanged.
func TrackDeleteAttempts(ctx context.Context, method string) context.Context {
	if method != http.MethodDelete || ctx.Value(deleteAttemptsKey{}) != nil {
		return ctx
	}
	return context.WithValue(ctx, deleteAttemptsKey{}, &deleteAttempts{})
}

// TrackDeleteRequestAttempts is TrackDeleteAttempts for a request that is
// already built. Use it only where the request is not rebuilt for a resend.
func TrackDeleteRequestAttempts(req *http.Request) *http.Request {
	ctx := TrackDeleteAttempts(req.Context(), req.Method)
	if ctx == req.Context() {
		return req
	}
	return req.WithContext(ctx)
}

// TrackAttempts wraps base so that every attempt of a tracked DELETE request
// is recorded. Requests without tracking pass through unchanged.
func TrackAttempts(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return attemptTrackingTransport{base: base}
}

type attemptTrackingTransport struct {
	base http.RoundTripper
}

func (t attemptTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if a, ok := req.Context().Value(deleteAttemptsKey{}).(*deleteAttempts); ok && mayHaveApplied(resp, err) {
		a.mayHaveApplied.Store(true)
	}
	return resp, err
}

// mayHaveApplied reports whether an attempt that ended with resp and err can
// have deleted the object without the client seeing a success.
func mayHaveApplied(resp *http.Response, err error) bool {
	if err != nil || resp == nil {
		return true
	}
	switch {
	case resp.StatusCode == http.StatusConflict, resp.StatusCode == http.StatusPreconditionFailed:
		return true
	case resp.StatusCode >= 500:
		return true
	}
	return false
}

// DeleteAlreadyApplied reports whether resp is a 404 for a DELETE that an
// earlier attempt can have applied. ctx must be the context of the tracked
// request.
func DeleteAlreadyApplied(ctx context.Context, method string, resp *http.Response) bool {
	if method != http.MethodDelete || resp == nil || resp.StatusCode != http.StatusNotFound {
		return false
	}
	a, ok := ctx.Value(deleteAttemptsKey{}).(*deleteAttempts)
	return ok && a.mayHaveApplied.Load()
}
