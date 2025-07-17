package wiz

import (
	"errors"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"os"
	"strings"
)

type wizConfig struct {
	ClientId     *string `hcl:"client_id"`
	ClientSecret *string `hcl:"client_secret"`
	Url          *string `hcl:"url"`
	AuthUrl      *string `hcl:"auth_url"`
}

func (wc wizConfig) validate() error {
	var errorStrings []string
	if wc.Url == nil || len(*wc.Url) == 0 {
		errorStrings = append(errorStrings, "url must be configured")
	}

	if wc.ClientId == nil || len(*wc.ClientId) == 0 {
		errorStrings = append(errorStrings, "client_id must be configured")
	}

	if wc.ClientSecret == nil || len(*wc.ClientSecret) == 0 {
		errorStrings = append(errorStrings, "client_secret must be configured")
	}

	if wc.AuthUrl == nil || len(*wc.AuthUrl) == 0 {
		errorStrings = append(errorStrings, "auth_url must be configured")
	}

	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, "\n"))
	}

	return nil
}

func ConfigInstance() interface{} {
	return &wizConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) wizConfig {
	config := wizConfig{}
	if connection != nil && connection.Config != nil {
		config, _ = connection.Config.(wizConfig)
	}

	// If credentials set in the config, override it
	if config.ClientId == nil {
		config.ClientId = ptr(os.Getenv("WIZ_AUTH_CLIENT_ID"))
	}
	if config.ClientSecret == nil {
		config.ClientSecret = ptr(os.Getenv("WIZ_AUTH_CLIENT_SECRET"))
	}
	if config.Url == nil {
		config.Url = ptr(os.Getenv("WIZ_URL"))
	}
	if config.AuthUrl == nil {
		config.AuthUrl = ptr(os.Getenv("WIZ_AUTH_URL"))
	}

	if len(*config.AuthUrl) == 0 {
		config.AuthUrl = ptr("https://auth.app.wiz.io/oauth/token")
	}

	return config
}

func ptr(x string) *string {
	return &x
}
