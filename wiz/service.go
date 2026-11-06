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

const wizAuthUrl = "https://auth.app.wiz.io/oauth/token"

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
	// Get the credentials
	clientId, clientSecret, apiUrl := getCredentialsByPrecedence(d)

	// No creds
	if apiUrl == "" {
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
		Url:   &apiUrl,
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

	clientId, clientSecret, _ := getCredentialsByPrecedence(d)
	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("client_id and client_secret must be configured")
	}

	authData := url.Values{}
	authData.Set("grant_type", "client_credentials")
	authData.Set("audience", "wiz-api")
	authData.Set("client_id", clientId)
	authData.Set("client_secret", clientSecret)

	// Create request
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, wizAuthUrl, strings.NewReader(authData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	r.Header.Add("Encoding", "UTF-8")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	plugin.Logger(ctx).Debug("GetAccessToken", "Getting access token...")

	// Make a request to get the token
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error getting token response: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate access token, please verify your connection config. Status: %s", resp.Status)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.Logger(ctx).Error("Failed reading response body", "error", err)
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}

	// Parse the response body of type accessToken{}
	at := accessToken{}
	if err := json.Unmarshal(bodyBytes, &at); err != nil {
		plugin.Logger(ctx).Error("failed to parse token response", "error", err)
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}
	if at.AccessToken == "" {
		return nil, fmt.Errorf("token response did not contain an access token")
	}

	return &at, nil
}

/*
Returns credentials by precedence.

Precedence of credentials:
  - Credentials set in config
  - Value set using WIZ_AUTH_CLIENT_ID, WIZ_AUTH_CLIENT_SECRET, and WIZ_URL env variables
*/
func getCredentialsByPrecedence(d *plugin.QueryData) (clientId string, clientSecret string, apiUrl string) {
	// Get wiz config
	wizConfig := GetConfig(d.Connection)

	// Check for env vars
	clientId = os.Getenv("WIZ_AUTH_CLIENT_ID")
	clientSecret = os.Getenv("WIZ_AUTH_CLIENT_SECRET")
	apiUrl = os.Getenv("WIZ_URL")

	// If credentials set in the config, override it
	if wizConfig.ClientId != nil && *wizConfig.ClientId != "" {
		clientId = *wizConfig.ClientId
	}
	if wizConfig.ClientSecret != nil && *wizConfig.ClientSecret != "" {
		clientSecret = *wizConfig.ClientSecret
	}
	if wizConfig.Url != nil && *wizConfig.Url != "" {
		apiUrl = *wizConfig.Url
	}

	return
}
