package wiz

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type wizConfig struct {
	ClientId     *string `hcl:"client_id"`
	ClientSecret *string `hcl:"client_secret"`
	Url          *string `hcl:"url"`
}

func ConfigInstance() interface{} {
	return &wizConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) wizConfig {
	if connection == nil || connection.Config == nil {
		return wizConfig{}
	}
	config, _ := connection.Config.(wizConfig)
	return config
}
