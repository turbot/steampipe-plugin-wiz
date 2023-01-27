package wiz

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type wizConfig struct {
	ApiToken *string `cty:"api_token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_token": {
		Type: schema.TypeString,
	},
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
