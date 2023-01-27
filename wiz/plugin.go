package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-wiz"

// Plugin creates this (wiz) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: pluginName,
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromCamel().Transform(transform.NullIfZeroValue),
		TableMap: map[string]*plugin.Table{
			"wiz_user":      tableWizUser(ctx),
			"wiz_user_role": tableWizUserRole(ctx),
		},
	}
	return p
}
