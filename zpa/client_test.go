package zpa

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStruct struct {
	ID int `json:"id"`
}

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

const (
	getResponse  = `{"id": 1234}`
	authResponse = `{
	"token_type": "token_type",
	"access_token": "access_token"
}`
)

func TestClient_NewRequestDo(t *testing.T) {
	type args struct {
		method string
		url    string
		body   interface{}
		v      interface{}
	}
	tests := []struct {
		name       string
		args       args
		muxHandler func(w http.ResponseWriter, r *http.Request)
		wantResp   *http.Response
		wantErr    bool
		wantVal    *dummyStruct
	}{
		// NewRequestDo test cases
		{
			name: "GET happy path",
			args: struct {
				method string
				url    string
				body   interface{}
				v      interface{}
			}{
				method: "GET",
				url:    "/test",
				body:   nil,
				v:      new(dummyStruct),
			},
			muxHandler: func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(getResponse))
				w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
				if err != nil {
					t.Fatal(err)
				}
			},
			wantResp: &http.Response{
				StatusCode: 200,
			},
			wantVal: &dummyStruct{
				ID: 1234,
			},
		},
	}

	for _, tt := range tests {
		client = NewClient(setupMuxConfig())
		client.WriteLog("Server URL: %v", client.Config.BaseURL)
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.args.url, tt.muxHandler)
			res, err := client.NewRequestDo(tt.args.method, tt.args.url, nil, tt.args.body, tt.args.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewRequestDo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp.StatusCode != res.StatusCode {
				t.Errorf("Client.NewRequestDo() = %v, want %v", res, tt.wantResp)
			}

			if !reflect.DeepEqual(tt.args.v, tt.wantVal) {
				t.Errorf("returned %#v; want %#v", tt.args.v, tt.wantVal)
			}
		})
	}
	teardown()
}

func TestNewClient(t *testing.T) {
	os.Setenv(ZPA_CLIENT_ID, "ClientID")
	os.Setenv(ZPA_CLIENT_SECRET, "ClientSecret")
	os.Setenv(ZPA_CUSTOMER_ID, "CustomerID")
	os.Setenv(ZPA_CLOUD, "test")
	type args struct {
		config *Config
	}
	tests := []struct {
		name  string
		args  args
		wantC *Config
	}{
		// NewClient test cases
		{
			name: "Successful Client creation with default config values",
			args: struct{ config *Config }{config: nil},
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.private.zscaler.com",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name: "Successful Client creation with custom config values",
			args: struct{ config *Config }{config: &Config{
				BaseURL:      &url.URL{Host: "https://otherhost.com"},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			}},
			wantC: &Config{
				BaseURL:      &url.URL{Host: "https://otherhost.com"},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC := NewClient(tt.args.config)
			assert.Equal(t, gotC.Config.BaseURL.Host, tt.wantC.BaseURL.Host)
			assert.Equal(t, gotC.Config.BaseURL.Scheme, tt.wantC.BaseURL.Scheme)
			assert.Equal(t, gotC.Config.ClientID, tt.wantC.ClientID)
			assert.Equal(t, gotC.Config.ClientSecret, tt.wantC.ClientSecret)
		})
	}
}

func setupMuxConfig() *Config {
	mux = http.NewServeMux()
	mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Set-Cookie", "JSESSIONID=JSESSIONID;")
		_, err := w.Write([]byte(authResponse))
		if err != nil {
			log.Fatal(err)
		}
	})
	server = httptest.NewServer(mux)
	config, err := NewConfig("clientID", "clientID", "customerID", "cloud", "userAgent")
	if err != nil {
		panic(err)
	}
	url, _ := url.Parse(server.URL)
	config.BaseURL = url
	return config
}

func teardown() {
	server.Close()
}
