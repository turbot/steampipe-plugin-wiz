package wiz

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

// getClient:: returns vanta client after authentication
func getClient(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "wiz"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*api.Client), nil
	}

	// Get the config
	wizConfig := GetConfig(d.Connection)

	var token string
	if wizConfig.ApiToken != nil {
		token = *wizConfig.ApiToken
	}

	// Return if no credential specified
	if token == "" {
		return nil, fmt.Errorf("api_token must be configured")
	}

	// Start with an empty Wiz config
	config := api.ClientConfig{ApiToken: wizConfig.ApiToken}

	// Create the client
	client, err := api.CreateClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
