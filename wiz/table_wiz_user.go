package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_user",
		Description: "Wiz User",
		List: &plugin.ListConfig{
			Hydrate: listWizUsers,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "identity_provider_type", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizUser,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the user."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email address of the user."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the user was created."},
			{Name: "last_login_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the user was last login to the console."},
			{Name: "is_suspended", Type: proto.ColumnType_BOOL, Description: "If true, the user is suspended."},
			{Name: "identity_provider_type", Type: proto.ColumnType_STRING, Description: "The auth provider type of the user. Possible values are: WIZ, SAML."},
			{Name: "is_analytics_enabled", Type: proto.ColumnType_BOOL, Description: "If true, the user analytics is enabled."},
			{Name: "ip_address", Type: proto.ColumnType_IPADDR, Description: "The IP address of the user."},
			{Name: "role", Type: proto.ColumnType_JSON, Description: "Specifies the role assigned to the user."},
			{Name: "effective_assigned_projects", Type: proto.ColumnType_JSON, Description: "A list of project ids this user was last logged in with. Null value means all projects are allowed."},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: "Specifies the tenant the user is a member of.", Transform: transform.FromField("Tenant.Id")},
		},
	}
}

//// LIST FUNCTION

func listWizUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_user.listWizUsers", "connection_error", err)
		return nil, err
	}

	options := &api.ListUsersRequestConfiguration{}

	// Default set to 600.
	// This is the maximum number of items can be requested
	pageLimit := 60

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Check for additional filters
	if d.EqualsQualString("identity_provider_type") != "" {
		options.AuthProviderType = d.EqualsQualString("identity_provider_type")
	}

	for {
		query, err := api.ListUsers(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_user.listWizUsers", "query_error", err)
			return nil, err
		}

		for _, user := range query.Users.Nodes {
			d.StreamListItem(ctx, user)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Users.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Users.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_user.getWizUser", "connection_error", err)
		return nil, err
	}

	query, err := api.GetUser(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_user.getWizUser", "query_error", err)
		return nil, err
	}

	return query.User, nil
}
