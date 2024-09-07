package zidentity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ZIDENTITY_CLIENT_ID     = "ZIDENTITY_CLIENT_ID"
	ZIDENTITY_CLIENT_SECRET = "ZIDENTITY_CLIENT_SECRET"
	ZIDENTITY_VANITY_DOMAIN = "ZIDENTITY_VANITY_DOMAIN" // Updated to vanity domain
)

type AuthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Credentials struct {
	AuthToken         *AuthToken
	ClientID          string
	ClientSecret      string
	Oauth2ProviderUrl string // This can still exist for backward compatibility
}

func Authenticate(clientID, clientSecret, vanityDomain, userAgent string, httpClient *http.Client) (*AuthToken, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("no client credentials were provided")
	}

	// Ensure the vanity domain is provided and does not include protocol
	if !strings.HasPrefix(vanityDomain, "https://") && vanityDomain != "" {
		vanityDomain = fmt.Sprintf("https://%s.zslogin.net/oauth2/v1/token", vanityDomain)
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", clientSecret)
	data.Set("client_id", clientID)
	data.Set("audience", "https://api.zscaler.com")
	authUrl := vanityDomain

	req, err := http.NewRequest("POST", authUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to signin the user %s, err: %v", clientID, err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to signin the user %s, err: %v", clientID, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to signin the user %s, err: %v", clientID, err)
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("[ERROR] Failed to signin the user %s, got http status:%d, response body:%s", clientID, resp.StatusCode, respBody)
	}
	var a AuthToken
	err = json.Unmarshal(respBody, &a)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to signin the user %s, err: %v", clientID, err)
	}

	return &a, nil
}
