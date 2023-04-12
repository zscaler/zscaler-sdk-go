package zcbc

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStruct struct {
	ID int `json:"id"`
}

var (
	mux    *http.ServeMux
	server *httptest.Server
)

const (
	getResponse  = `{"id": 1234}`
	authResponse = `{
	"authType": "authType",
    "obfuscateApiKey": true,
    "passwordExpiryTime": 10000,
    "passwordExpiryDays": 10000,
    "source": "source",
    "jSessionID": ""
}`
)

func TestClient_Request(t *testing.T) {
	defer teardown()
	type args struct {
		method string
		url    string
		body   []byte
		v      interface{}
	}
	tests := []struct {
		name       string
		args       args
		muxHandler func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		wantVal    *dummyStruct
	}{
		// NewRequestDo test cases
		{
			name: "GET happy path",
			args: struct {
				method string
				url    string
				body   []byte
				v      interface{}
			}{
				method: "GET",
				url:    "/test",
				body:   nil,
				v:      new(dummyStruct),
			},
			muxHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				_, err := w.Write([]byte(getResponse))
				if err != nil {
					t.Fatal(err)
				}
				// panic(fmt.Sprintf("%v", r.Header))
			},
			wantVal: &dummyStruct{
				ID: 1234,
			},
		},
	}

	for _, tt := range tests {
		client := setupMuxAndClient()
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc("/authenticatedSession", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				w.Header().Add("Set-Cookie", "JSESSIONID=JSESSIONID;")
				_, err := w.Write([]byte(authResponse))
				if err != nil {
					log.Fatal(err)
				}
			})
			mux.HandleFunc(tt.args.url, tt.muxHandler)
			resp, err := client.Request(tt.args.url, tt.args.method, tt.args.body, client.GetContentType())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewRequestDo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = json.Unmarshal(resp, &tt.args.v)
			if !reflect.DeepEqual(tt.args.v, tt.wantVal) {
				t.Errorf("returned %#v; want %#v", tt.args.v, tt.wantVal)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		userName  string
		password  string
		apiKey    string
		ziaCloud  string
		UserAgent string
	}
	tests := []struct {
		name  string
		args  args
		wantC *Client
	}{
		{
			name: "Successful Client creation with custom config values",
			args: struct {
				userName  string
				password  string
				apiKey    string
				ziaCloud  string
				UserAgent string
			}{
				ziaCloud:  "ziaCloud",
				userName:  "userName",
				password:  "password",
				apiKey:    "apiKey",
				UserAgent: "UserAgent",
			},
			wantC: &Client{
				URL:       fmt.Sprintf("https://zsapi.%s.net/%s", "ziaCloud", ziaAPIVersion),
				userName:  "userName",
				password:  "password",
				apiKey:    "apiKey",
				UserAgent: "UserAgent",
			},
		},
		{
			name: "Successful Client creation with custom config values",
			args: struct {
				userName  string
				password  string
				apiKey    string
				ziaCloud  string
				UserAgent string
			}{
				ziaCloud:  "zspreview",
				userName:  "userName",
				password:  "password",
				apiKey:    "apiKey",
				UserAgent: "UserAgent",
			},
			wantC: &Client{
				URL:       fmt.Sprintf("https://admin.%s.net/%s", "zspreview", ziaAPIVersion),
				userName:  "userName",
				password:  "password",
				apiKey:    "apiKey",
				UserAgent: "UserAgent",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewClient(tt.args.userName, tt.args.password, tt.args.apiKey, tt.args.ziaCloud, tt.args.UserAgent)
			if err != nil {
				t.Errorf("NewClient error = %v, wantErr nil", err)
				return
			}
			assert.Equal(t, gotC.URL, tt.wantC.URL)
			assert.NotNil(t, gotC.HTTPClient)
			assert.Equal(t, gotC.apiKey, tt.wantC.apiKey)
			assert.Equal(t, gotC.password, tt.wantC.password)
			assert.Equal(t, gotC.userName, tt.wantC.userName)
			assert.Equal(t, gotC.UserAgent, tt.wantC.UserAgent)
		})
	}
}

func setupMuxAndClient() *Client {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	cli, err := NewClient("username", "password", "apiKey------", "cloud", "")
	if err != nil {
		panic(err)
	}
	cli.URL = server.URL

	return cli
}

func teardown() {
	server.Close()
}
*/
