package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizSubscription(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_subscription",
		Description: "Wiz Subscription",
		List: &plugin.ListConfig{
			Hydrate: listWizSubscriptions,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "cloud_provider", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSubscription,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name for the account."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the account."},
			{Name: "cloud_provider", Type: proto.ColumnType_STRING, Description: "The type of the cloud provider. Possible values are: AWS, GCP, OCI, Alibaba, Azure, Kubernetes, OpenShift, vSphere."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Cloud Account connectivity status as affected by configured connectors.  Possible values: CONNECTED, DISABLED, DISCONNECTED, DISCOVERED, ERROR, INITIAL_SCANNING, PARTIALLY_CONNECTED."},
			{Name: "last_scanned_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the cloud account was last scanned."},
			{Name: "container_count", Type: proto.ColumnType_INT, Description: "Number of containers that are part of this cloud account."},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "External subscription ID from a cloud provider (subscriptionId in security graph)."},
			{Name: "resource_count", Type: proto.ColumnType_INT, Description: "Number of resources that are part of this cloud account."},
			{Name: "virtual_machine_count", Type: proto.ColumnType_INT, Description: "Number of virtual machines that are part of this cloud account."},
			{Name: "linked_projects", Type: proto.ColumnType_JSON, Description: "A list of projects, this cloud account is assigned to."},
		},
	}
}

//// LIST FUNCTION

func listWizSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_subscription.listWizSubscriptions", "connection_error", err)
		return nil, err
	}

	options := &api.ListSubscriptionsRequestConfiguration{}

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
	if d.EqualsQualString("status") != "" {
		options.Status = d.EqualsQualString("status")
	}

	for {
		query, err := api.ListSubscriptions(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_subscription.listWizSubscriptions", "query_error", err)
			return nil, err
		}

		for _, subscription := range query.Subscriptions.Nodes {
			d.StreamListItem(ctx, subscription)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Subscriptions.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Subscriptions.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSubscription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_subscription.getSubscription", "connection_error", err)
		return nil, err
	}

	query, err := api.GetSubscription(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_subscription.getSubscription", "query_error", err)
		return nil, err
	}

	return query.Subscription, nil
}
