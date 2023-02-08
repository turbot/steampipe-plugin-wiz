package wiz

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizControl(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_control",
		Description: "Wiz Control",
		List: &plugin.ListConfig{
			Hydrate: listWizControls,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "has_auto_remediation", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getControl,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the control."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The wiz identifier for the control."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The control severity. Possible values are: CRITICAL, HIGH, INFORMATIONAL, LOW, MEDIUM."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "True, if the control is enabled."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The control description."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The control type. Possible values are: CLOUD_CONFIGURATION, SECURITY_GRAPH."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the control was created."},
			{Name: "enabled_for_hbi", Type: proto.ColumnType_BOOL, Description: "If true, the control has High Business Impact (HBI).", Transform: transform.FromField("EnabledForHBI")},
			{Name: "enabled_for_mbi", Type: proto.ColumnType_BOOL, Description: "If true, the control has Medium Business Impact (MBI).", Transform: transform.FromField("EnabledForMBI")},
			{Name: "enabled_for_lbi", Type: proto.ColumnType_BOOL, Description: "If true, the control has Low Business Impact (LBI).", Transform: transform.FromField("EnabledForLBI")},
			{Name: "enabled_for_unattributed", Type: proto.ColumnType_BOOL, Description: "Enables control for projects which are not LBI/MBI/HBI. All controls should have this set to true by default."},
			{Name: "has_auto_remediation", Type: proto.ColumnType_BOOL, Description: "If true, the cloud configuration has auto remediation enabled."},
			{Name: "last_run_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the control was last run."},
			{Name: "last_run_error", Type: proto.ColumnType_STRING, Description: "The error that the controls gets during the last run, if any."},
			{Name: "last_successful_run_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time of a successful control run."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The project ID this control is scoped to."},
			{Name: "resolution_recommendation", Type: proto.ColumnType_STRING, Description: "The guidance on how the user should address an issue that was created by this control."},
			{Name: "supports_nrt", Type: proto.ColumnType_BOOL, Description: "", Transform: transform.FromField("SupportsNRT")},
			{Name: "created_by", Type: proto.ColumnType_JSON, Description: "The owner information of the control."},
			{Name: "source_cloud_configuration_rule", Type: proto.ColumnType_JSON, Description: "The information about the cloud configuration rule."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The list of tags associated with the control"},
			{Name: "query", Type: proto.ColumnType_JSON, Description: "The query that the control runs. If query is null, this is a built in control with custom logic."},
		},
	}
}

//// LIST FUNCTION

func listWizControls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_control.listWizControls", "connection_error", err)
		return nil, err
	}

	options := &api.ListControlsRequestConfiguration{}

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
	if d.EqualsQuals["enabled"] != nil {
		options.Enabled = types.Bool(d.EqualsQuals["enabled"].GetBoolValue())
	}
	if d.EqualsQuals["has_auto_remediation"] != nil {
		options.HasAutoRemediation = types.Bool(d.EqualsQuals["has_auto_remediation"].GetBoolValue())
	}
	if d.EqualsQualString("project_id") != "" {
		options.Project = d.EqualsQualString("project_id")
	}
	if d.EqualsQualString("severity") != "" {
		options.Severity = d.EqualsQualString("severity")
	}
	if d.EqualsQualString("type") != "" {
		options.Type = d.EqualsQualString("type")
	}

	// Check for not equal qual
	filterQuals := []string{"enabled", "has_auto_remediation"}
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
					}
					break
				}
			}
		}
	}

	for {
		query, err := api.ListControls(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_control.listWizControls", "query_error", err)
			return nil, err
		}

		for _, control := range query.Controls.Nodes {
			d.StreamListItem(ctx, control)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Controls.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Controls.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getControl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_control.getControl", "connection_error", err)
		return nil, err
	}

	query, err := api.GetControl(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_control.getControl", "query_error", err)
		return nil, err
	}

	return query.Control, nil
}
