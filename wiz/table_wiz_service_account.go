package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizServiceAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_service_account",
		Description: "Wiz Service Account",
		List: &plugin.ListConfig{
			Hydrate: listWizServiceAccounts,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "authentication_source", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getServiceAccount,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the service account."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the service account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the service account. Possible values are: THIRD_PARTY, SENSOR, KUBERNETES_ADMISSION_CONTROLLER, BROKER."},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "The client ID of the service account."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the service account was created."},
			{Name: "authentication_source", Type: proto.ColumnType_STRING, Description: "The authentication source of the service account. Possible values: LEGACY, MODERN."},
			{Name: "last_rotated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the service account was last rotated."},
			{Name: "assigned_projects", Type: proto.ColumnType_JSON, Description: "A list of projects where the service account was assigned. If null, the service account has access to all projects."},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "A list of actions the service account is allowed to perform."},
		},
	}
}

//// LIST FUNCTION

func listWizServiceAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_service_account.listWizServiceAccounts", "connection_error", err)
		return nil, err
	}

	options := &api.ListServiceAccountsRequestConfiguration{}

	// Default set to 60.
	// This is the maximum number of items can be requested.
	// If the value is >60, API returns following error:
	// 'first' must less than or equal to 60
	pageLimit := 60

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Check for additional filters
	if d.EqualsQualString("name") != "" {
		options.Name = d.EqualsQualString("name")
	}
	if d.EqualsQualString("type") != "" {
		options.Type = d.EqualsQualString("type")
	}
	if d.EqualsQualString("authentication_source") != "" {
		options.Source = d.EqualsQualString("authentication_source")
	}

	for {
		query, err := api.ListServiceAccounts(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_service_account.listWizServiceAccounts", "query_error", err)
			return nil, err
		}

		for _, serviceAccount := range query.ServiceAccounts.Nodes {
			d.StreamListItem(ctx, serviceAccount)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.ServiceAccounts.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.ServiceAccounts.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_service_account.getServiceAccount", "connection_error", err)
		return nil, err
	}

	query, err := api.GetServiceAccount(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_service_account.getServiceAccount", "query_error", err)
		return nil, err
	}

	return query.ServiceAccount, nil
}
