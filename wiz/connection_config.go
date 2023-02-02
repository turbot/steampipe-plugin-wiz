package wiz

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type wizConfig struct {
	ApiToken     *string `cty:"api_token"`
	ClientId     *string `cty:"client_id"`
	ClientSecret *string `cty:"client_secret"`
	Url          *string `cty:"url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_token": {
		Type: schema.TypeString,
	},
	"client_id": {
		Type: schema.TypeString,
	},
	"client_secret": {
		Type: schema.TypeString,
	},
	"url": {
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
