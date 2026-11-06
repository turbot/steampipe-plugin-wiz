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
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "builtin", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "is_project_scoped", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the user role."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the user role."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the user role."},
			{Name: "is_project_scoped", Type: proto.ColumnType_BOOL, Description: "If true, the role is scoped to a specific project."},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "A list of operation can be performed using the role."},
			{Name: "builtin", Type: proto.ColumnType_BOOL, Description: "If true, the role is a built-in Wiz role."},
			{Name: "is_tenant_group_scoped", Type: proto.ColumnType_BOOL, Description: "If true, the role is scoped to a tenant group."},
			{Name: "user_count", Type: proto.ColumnType_INT, Description: "The number of users assigned to the role."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the role was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the role was last updated."},
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
	options.Builtin = boolQual(d, "builtin")
	options.IsProjectScoped = boolQual(d, "is_project_scoped")

	// The API returns all roles in a single response
	query, err := api.ListUserRoles(context.Background(), conn, options)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_user_role.listWizUserRoles", "query_error", err)
		return nil, err
	}

	for _, role := range query.UserRoles {
		d.StreamListItem(ctx, role)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
