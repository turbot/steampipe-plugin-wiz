package wiz

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizCloudConfigRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_cloud_config_rule",
		Description: "Wiz Cloud Configuration Rule",
		List: &plugin.ListConfig{
			Hydrate: listWizCloudConfigRules,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "cloud_provider", Require: plugin.Optional},
				{Name: "enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "has_auto_remediation", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "service_type", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "supports_nrt", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getCloudConfigRule,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the cloud configuration rule."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the cloud configuration rule."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Rule enabled status."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Rule severity will outcome to finding severity. This filed initial value is set as the severity of the CSPM rule. Possible values are: CRITICAL, HIGH, INFORMATIONAL, LOW, MEDIUM."},
			{Name: "cloud_provider", Type: proto.ColumnType_STRING, Description: "The cloud provider this rule is relevant to. Possible values are: AWS, GCP, OCI, Alibaba, Azure, Kubernetes, OpenShift, vSphere."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the rule was created."},
			{Name: "built_in", Type: proto.ColumnType_BOOL, Description: "Indicates whether the rule is built-in or custom."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the rule."},
			{Name: "short_id", Type: proto.ColumnType_STRING, Description: "A short unique identifier for the rule."},
			{Name: "function_as_role", Type: proto.ColumnType_BOOL, Description: "Make this rule also function as a control which means findings by this control will also trigger Issues."},
			{Name: "has_auto_remediation", Type: proto.ColumnType_BOOL, Description: "If true, the rule will automatically remediate the failed resources as per remediation steps defined in the rule."},
			{Name: "remediation_instructions", Type: proto.ColumnType_STRING, Description: "A set of instructions provided for the remediation."},
			{Name: "service_type", Type: proto.ColumnType_STRING, Description: "The service this rule is relevant to."},
			{Name: "supports_nrt", Type: proto.ColumnType_BOOL, Description: "Indicates the support of 'near real time' updates."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the rule was last modified."},
			{Name: "control", Type: proto.ColumnType_JSON, Description: "Specifies the control information, in case this rule also functions as a control."},
			{Name: "created_by", Type: proto.ColumnType_JSON, Description: "Specifies the user object that created the rule."},
			{Name: "scoped_accounts", Type: proto.ColumnType_JSON, Description: "A list of target cloud accounts where the rule is applied to. If empty, the rule will run on all environment."},
			{Name: "target_native_types", Type: proto.ColumnType_JSON, Description: "The identifier types of the objects targeted by this rule, as seen on the cloud provider service. e.g. 'ec2'."},
		},
	}
}

//// LIST FUNCTION

func listWizCloudConfigRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_config_rule.listWizCloudConfigRules", "connection_error", err)
		return nil, err
	}

	options := &api.ListCloudConfigRulesRequestConfiguration{}

	// Default set to 500.
	// This is the maximum number of items can be requested.
	pageLimit := 500

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Check for additional filters
	if d.EqualsQualString("cloud_provider") != "" {
		options.CloudProvider = d.EqualsQualString("cloud_provider")
	}
	if d.EqualsQualString("service_type") != "" {
		options.ServiceType = d.EqualsQualString("service_type")
	}
	if d.EqualsQualString("severity") != "" {
		options.Severity = d.EqualsQualString("severity")
	}
	if d.EqualsQuals["enabled"] != nil {
		options.Enabled = types.Bool(d.EqualsQuals["enabled"].GetBoolValue())
	}
	if d.EqualsQuals["has_auto_remediation"] != nil {
		options.HasAutoRemediation = types.Bool(d.EqualsQuals["has_auto_remediation"].GetBoolValue())
	}
	if d.EqualsQuals["supports_nrt"] != nil {
		options.SupportsNRT = types.Bool(d.EqualsQuals["supports_nrt"].GetBoolValue())
	}

	filterQuals := []string{
		"enabled",
		"has_auto_remediation",
		"supports_nrt",
	}

	// Check for not equal qual
	for _, qual := range filterQuals {
		if d.Quals[qual] != nil {
			for _, q := range d.Quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					switch qual {
					case "enabled":
						options.Enabled = types.Bool(!value)
					case "has_auto_remediation":
						options.HasAutoRemediation = types.Bool(!value)
					case "supports_nrt":
						options.SupportsNRT = types.Bool(!value)
					}
					break
				}
			}
		}
	}

	for {
		query, err := api.ListCloudConfigRules(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_cloud_config_rule.listWizCloudConfigRules", "query_error", err)
			return nil, err
		}

		for _, rule := range query.CloudConfigRules.Nodes {
			d.StreamListItem(ctx, rule)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.CloudConfigRules.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.CloudConfigRules.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudConfigRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_config_rule.getCloudConfigRule", "connection_error", err)
		return nil, err
	}

	query, err := api.GetCloudConfigRule(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_config_rule.getCloudConfigRule", "query_error", err)
		return nil, err
	}

	return query.CloudConfigRule, nil
}
