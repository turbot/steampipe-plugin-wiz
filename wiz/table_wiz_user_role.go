package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizUserRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_user_role",
		Description: "Wiz User Role",
		List: &plugin.ListConfig{
			Hydrate: listWizUserRoles,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the user role."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the user role."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the user role."},
			{Name: "is_project_scoped", Type: proto.ColumnType_BOOL, Description: "If true, the role is scoped to a specific project."},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "A list of operation can be performed using the role."},
		},
	}
}

//// LIST FUNCTION

func listWizUserRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_user_role.listWizUserRoles", "connection_error", err)
		return nil, err
	}

	options := &api.ListUserRolesRequestConfiguration{}

	// Default set to 100.
	// This is the maximum number of items can be requested
	pageLimit := 100

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	for {
		query, err := api.ListUserRoles(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_user_role.listWizUserRoles", "query_error", err)
			return nil, err
		}

		for _, role := range query.UserRoles.Nodes {
			d.StreamListItem(ctx, role)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.UserRoles.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.UserRoles.PageInfo.EndCursor
	}

	return nil, nil
}
