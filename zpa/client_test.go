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

	"github.com/zscaler/zscaler-sdk-go/logger"
)

type dummyStruct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

const (
	getResponse  = `{"id": 1234,"name":"name with &amp;amp;, &amp;lt; and &amp;gt;","description":"description with &amp;amp;, &amp;lt; and &amp;gt;"}`
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
				ID:          1234,
				Name:        "name with &, < and >",
				Description: "description with &, < and >",
			},
		},
	}

	for _, tt := range tests {
		client = NewClient(setupMuxConfig())
		logger.WriteLog(client.Config.Logger, "Server URL: %v", client.Config.BaseURL)
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
	type args struct {
		config *Config
	}
	tests := []struct {
		name  string
		args  args
		cloud string
		wantC *Config
	}{
		// NewClient test cases
		{
			name:  "Successful Client creation with default config values",
			args:  struct{ config *Config }{config: nil},
			cloud: "",
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
			name:  "Production cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "production",
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
			name:  "Beta cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "beta",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.zpabeta.net",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "Gov cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "gov",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.zpagov.net",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "GovUS cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "govus",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.zpagov.us",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "Preview cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "preview",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.zpapreview.net",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "QA cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "qa",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.qa.zpath.net",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "Arbitrary cloud support",
			args:  struct{ config *Config }{config: nil},
			cloud: "https://config.somecloud.net",
			wantC: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "config.somecloud.net",
				},
				ClientID:     "ClientID",
				ClientSecret: "ClientSecret",
				CustomerID:   "CustomerID",
				UserAgent:    "userAgent",
			},
		},
		{
			name:  "Successful Client creation with custom config values",
			cloud: "",
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
			os.Setenv(ZPA_CLOUD, tt.cloud)
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
