package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

type accessToken struct {
	Token string `json:"access_token"`
}

// getClient:: returns Wiz client after authentication
func getClient(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "wiz"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*api.Client), nil
	}

	// Get the config
	wizConfig := GetConfig(d.Connection)

	var clientId, clientSecret, token string

	/* Credential precedence
	 * api_token set in config; if empty,
	 * client_id and client_secret set in config
	 */
	if wizConfig.ApiToken != nil {
		token = *wizConfig.ApiToken
	}

	if wizConfig.ClientId != nil && wizConfig.ClientSecret != nil {
		clientId = *wizConfig.ClientId
		clientSecret = *wizConfig.ClientSecret
	}

	// Return if no credential specified
	if token == "" && (clientId == "" || clientSecret == "") {
		return nil, fmt.Errorf("either api_token, or client_id, client_secret must be configured")
	}

	// Using client_id and client_secret credentials
	if token == "" {
		accessTokenResponse, err := GetAPIToken(ctx, d)
		if err != nil {
			return nil, err
		}
		token = accessTokenResponse.Token
	}

	// Start with an empty Wiz config
	config := api.ClientConfig{ApiToken: &token}

	// Create the client
	client, err := api.CreateClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

// GetAPIToken retrieves a new API token using the clientId and secret
func GetAPIToken(ctx context.Context, d *plugin.QueryData) (*accessToken, error) {

	// have we already created and cached the token?
	cacheKey := "wiz.session_token"
	if ts, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return ts.(*accessToken), nil
	}

	plugin.Logger(ctx).Debug("Creating session token", "connection", d.Connection.Name)

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

	plugin.Logger(ctx).Debug("GetAPIToken", "Getting token...")

	// Make a request to get the token
	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error getting token response. Status: %s, err: %v", resp.Status, err)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}

	// Parse the response body of type accessToken{}
	at := accessToken{}
	jsonErr := json.Unmarshal(bodyBytes, &at)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed parsing JSON body: %v", jsonErr)
	}

	return &at, nil
}