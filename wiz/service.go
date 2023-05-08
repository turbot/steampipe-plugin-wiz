package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

type accessToken struct {
	AccessToken string `json:"access_token"`
}

// getClient:: returns Wiz client after authentication
func getClient(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	conn, err := clientCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*api.Client), nil
}

// Get the cached version of the client
var clientCached = plugin.HydrateFunc(clientUncached).Memoize()

// clientUncached returns the Wiz client and cached the data
func clientUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	// Get the config
	wizConfig := GetConfig(d.Connection)

	// Get the credentials
	clientId, clientSecret, url := getCredentialsByPrecedence(d)

	// No creds
	if url == "" {
		return nil, fmt.Errorf("url must be configured")
	}

	// Return if no credential specified
	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("client_id and client_secret must be configured")
	}

	// Using client_id and client_secret credentials
	accessTokenResponse, err := GetAccessToken(ctx, d)
	if err != nil {
		return nil, err
	}
	accessToken := accessTokenResponse.AccessToken

	// Start with an empty Wiz config
	config := api.ClientConfig{
		Token: &accessToken,
		Url:   wizConfig.Url,
	}

	// Create the client
	client, err := api.CreateClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}

	return client, nil
}

// GetAccessToken retrieves a new access token using the clientId and secret
func GetAccessToken(ctx context.Context, d *plugin.QueryData) (*accessToken, error) {
	tokenResponse, err := accessTokenCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return tokenResponse.(*accessToken), nil
}

// Get the cached version of the token response
var accessTokenCached = plugin.HydrateFunc(accessTokenUncached).Memoize()

// accessTokenUncached returns the access token after authenticating using clientId and secret
func accessTokenUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("Creating access token", "connection", d.Connection.Name)

	// Get Wiz config
	wizConfig := GetConfig(d.Connection)

	clientId := wizConfig.ClientId
	clientSecret := wizConfig.ClientSecret

	auth_data := url.Values{}
	auth_data.Set("grant_type", "client_credentials")
	auth_data.Set("audience", "wiz-api")
	auth_data.Set("client_id", *clientId)
	auth_data.Set("client_secret", *clientSecret)

	// Create client
	client := &http.Client{}

	// Create request
	r, err := http.NewRequest(http.MethodPost, "https://auth.app.wiz.io/oauth/token", strings.NewReader(auth_data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	r.Header.Add("Encoding", "UTF-8")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	plugin.Logger(ctx).Debug("GetAccessToken", "Getting access token...")

	// Make a request to get the token
	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error getting token response. Status: %s, err: %v", resp.Status, err)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Logger(ctx).Error("Failed reading response body", "body", string(bodyBytes), "error", err)
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}

	// Parse the response body of type accessToken{}
	at := accessToken{}
	jsonErr := json.Unmarshal(bodyBytes, &at)
	if jsonErr != nil {
		plugin.Logger(ctx).Error("failed to parse token response", jsonErr)
		return nil, fmt.Errorf("failed to parse token response: %v", jsonErr)
	}

	return &at, nil
}

/*
Returns credentials by precedence.

Precedence of credentials:
  - Credentials set in config
  - Value set using WIZ_AUTH_CLIENT_ID, WIZ_AUTH_CLIENT_SECRET, and WIZ_URL env variables
*/
func getCredentialsByPrecedence(d *plugin.QueryData) (clientId string, clientSecret string, url string) {
	// Get wiz config
	wizConfig := GetConfig(d.Connection)

	// Check for env vars
	clientId = os.Getenv("WIZ_AUTH_CLIENT_ID")
	clientSecret = os.Getenv("WIZ_AUTH_CLIENT_SECRET")
	url = os.Getenv("WIZ_URL")

	// If credentials set in the config, override it
	if wizConfig.ClientId != nil {
		clientId = *wizConfig.ClientId
	}
	if wizConfig.ClientSecret != nil {
		clientSecret = *wizConfig.ClientId
	}
	if wizConfig.Url != nil {
		url = *wizConfig.Url
	}

	return
}
