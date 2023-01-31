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
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"wiz_project":         tableWizProject(ctx),
			"wiz_service_account": tableWizServiceAccount(ctx),
			"wiz_user":            tableWizUser(ctx),
			"wiz_user_role":       tableWizUserRole(ctx),
		},
	}
	return p
}
