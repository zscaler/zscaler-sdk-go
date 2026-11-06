package errorx

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// scriptedTransport answers each attempt with the next outcome of the script.
type scriptedTransport struct {
	outcomes []outcome
	sent     int
}

type outcome struct {
	status int
	err    error
}

func (t *scriptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	o := t.outcomes[t.sent]
	t.sent++
	if o.err != nil {
		return nil, o.err
	}
	return &http.Response{StatusCode: o.status, Body: http.NoBody, Request: req}, nil
}

// send runs every outcome of the script through the tracking transport, as the
// retries of one request would, and returns the context and the last response.
func send(t *testing.T, method string, ctx context.Context, outcomes ...outcome) (context.Context, *http.Response) {
	t.Helper()
	rt := TrackAttempts(&scriptedTransport{outcomes: outcomes})
	var resp *http.Response
	for range outcomes {
		req, err := http.NewRequestWithContext(ctx, method, "https://api.example.com/zia/api/v1/firewallFilteringRules/1495551", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, _ = rt.RoundTrip(req)
	}
	return ctx, resp
}

func TestDeleteAlreadyApplied(t *testing.T) {
	transportErr := errors.New("connection reset by peer")
	tests := []struct {
		name     string
		outcomes []outcome
		want     bool
	}{
		// SUP-4290: the DELETE was turned away (429), then applied but reported
		// as an edit lock conflict (409), and the retry got 404.
		{"429 then 409 then 404", []outcome{{status: 429}, {status: 409}, {status: 404}}, true},
		{"409 then 404", []outcome{{status: 409}, {status: 404}}, true},
		{"412 then 404", []outcome{{status: 412}, {status: 404}}, true},
		{"500 then 404", []outcome{{status: 500}, {status: 404}}, true},
		{"502 then 404", []outcome{{status: 502}, {status: 404}}, true},
		{"no response then 404", []outcome{{err: transportErr}, {status: 404}}, true},

		// The object did not exist: no earlier attempt can have deleted it.
		{"404 on the first attempt", []outcome{{status: 404}}, false},
		{"429 then 404", []outcome{{status: 429}, {status: 404}}, false},
		{"401 then 404", []outcome{{status: 401}, {status: 404}}, false},
		{"403 then 404", []outcome{{status: 403}, {status: 404}}, false},
		{"404 then 404", []outcome{{status: 404}, {status: 404}}, false},

		// Only a final 404 is reported as deleted.
		{"409 then 409", []outcome{{status: 409}, {status: 409}}, false},
		{"409 then 400", []outcome{{status: 409}, {status: 400}}, false},
		{"409 then 204", []outcome{{status: 409}, {status: 204}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, resp := send(t, http.MethodDelete, TrackDeleteAttempts(context.Background(), http.MethodDelete), tt.outcomes...)
			if got := DeleteAlreadyApplied(ctx, http.MethodDelete, resp); got != tt.want {
				t.Fatalf("DeleteAlreadyApplied = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteAlreadyApplied_OnlyTrackedDeletes(t *testing.T) {
	outcomes := []outcome{{status: 409}, {status: 404}}

	t.Run("other methods are never tracked", func(t *testing.T) {
		for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodPost} {
			ctx := TrackDeleteAttempts(context.Background(), method)
			if ctx.Value(deleteAttemptsKey{}) != nil {
				t.Fatalf("%s must not be tracked", method)
			}
			ctx, resp := send(t, method, ctx, outcomes...)
			if DeleteAlreadyApplied(ctx, method, resp) {
				t.Fatalf("a 404 for %s must not be reported as deleted", method)
			}
		}
	})

	t.Run("an untracked DELETE keeps its 404", func(t *testing.T) {
		ctx, resp := send(t, http.MethodDelete, context.Background(), outcomes...)
		if DeleteAlreadyApplied(ctx, http.MethodDelete, resp) {
			t.Fatal("a 404 for an untracked DELETE must not be reported as deleted")
		}
	})

	t.Run("nil response", func(t *testing.T) {
		ctx := TrackDeleteAttempts(context.Background(), http.MethodDelete)
		if DeleteAlreadyApplied(ctx, http.MethodDelete, nil) {
			t.Fatal("a nil response must not be reported as deleted")
		}
	})
}

// A request rebuilt for a resend shares the record of the first one, so the
// request loops of the clients can resend from the same context.
func TestTrackDeleteAttempts_SharedAcrossRebuiltRequests(t *testing.T) {
	ctx := TrackDeleteAttempts(context.Background(), http.MethodDelete)
	if again := TrackDeleteAttempts(ctx, http.MethodDelete); again != ctx {
		t.Fatal("tracking an already tracked context must return it unchanged")
	}

	rt := TrackAttempts(&scriptedTransport{outcomes: []outcome{{status: 409}, {status: 404}}})
	var resp *http.Response
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, "https://api.example.com/x/1", nil)
		resp, _ = rt.RoundTrip(req)
	}
	if !DeleteAlreadyApplied(ctx, http.MethodDelete, resp) {
		t.Fatal("rebuilt requests must share one record")
	}
}

func TestTrackDeleteRequestAttempts(t *testing.T) {
	del, _ := http.NewRequest(http.MethodDelete, "https://api.example.com/x/1", nil)
	tracked := TrackDeleteRequestAttempts(del)
	if tracked.Context().Value(deleteAttemptsKey{}) == nil {
		t.Fatal("a DELETE request must be tracked")
	}
	if TrackDeleteRequestAttempts(tracked) != tracked {
		t.Fatal("tracking an already tracked request must return it unchanged")
	}

	get, _ := http.NewRequest(http.MethodGet, "https://api.example.com/x/1", nil)
	if TrackDeleteRequestAttempts(get) != get {
		t.Fatal("a GET request must be returned unchanged")
	}
}
