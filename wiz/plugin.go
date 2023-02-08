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
		DefaultTransform: transform.FromCamel().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"wiz_cloud_config_rule":           tableWizCloudConfigRule(ctx),
			"wiz_cloud_configuration_finding": tableWizCloudConfigurationFinding(ctx),
			"wiz_control":                     tableWizControl(ctx),
			"wiz_issue":                       tableWizIssue(ctx),
			"wiz_project":                     tableWizProject(ctx),
			"wiz_security_framework":          tableWizSecurityFramework(ctx),
			"wiz_service_account":             tableWizServiceAccount(ctx),
			"wiz_subscription":                tableWizSubscription(ctx),
			"wiz_user":                        tableWizUser(ctx),
			"wiz_user_role":                   tableWizUserRole(ctx),
			"wiz_vulnerability_finding":       tableWizVulnerabilityFinding(ctx),
		},
	}
	return p
}
