package wiz

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizSecurityFramework(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_security_framework",
		Description: "Wiz Security Framework",
		List: &plugin.ListConfig{
			Hydrate: listWizSecurityFrameworks,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "enabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSecurityFramework,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the security framework."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the security framework."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "If true, the security framework is enabled."},
			{Name: "built_in", Type: proto.ColumnType_BOOL, Description: "If true, the security framework is managed by Wiz."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the security framework."},
			{Name: "policy_types", Type: proto.ColumnType_JSON, Description: "A list of security framework policy type."},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: "A list of security category."},
		},
	}
}

//// LIST FUNCTION

func listWizSecurityFrameworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_framework.listWizSecurityFrameworks", "connection_error", err)
		return nil, err
	}

	options := &api.ListSecurityFrameworksRequestConfiguration{}

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
		options.Enabled = types.ToBoolPtr(d.EqualsQuals["enabled"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["enabled"] != nil {
		for _, q := range d.Quals["enabled"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				options.Enabled = types.ToBoolPtr(false)
				if !value {
					options.Enabled = types.ToBoolPtr(true)
				}
			}
		}
	}

	for {
		query, err := api.ListSecurityFrameworks(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_security_framework.listWizSecurityFrameworks", "query_error", err)
			return nil, err
		}

		for _, framework := range query.SecurityFrameworks.Nodes {
			d.StreamListItem(ctx, framework)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.SecurityFrameworks.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.SecurityFrameworks.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityFramework(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_framework.getSecurityFramework", "connection_error", err)
		return nil, err
	}

	query, err := api.GetSecurityFramework(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_framework.getSecurityFramework", "query_error", err)
		return nil, err
	}

	return query.SecurityFramework, nil
}
