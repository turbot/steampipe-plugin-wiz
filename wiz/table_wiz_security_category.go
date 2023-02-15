package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizSecurityCategory(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_security_category",
		Description: "Wiz Security Category",
		List: &plugin.ListConfig{
			Hydrate: listWizSecurityCategories,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "framework_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSecurityCategory,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the category."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the category."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the category."},
			{Name: "framework_id", Type: proto.ColumnType_STRING, Description: "The ID of security framework, the category is part of.", Transform: transform.FromField("Framework.Id")},
			{Name: "sub_categories", Type: proto.ColumnType_JSON, Description: "A list of security category."},
		},
	}
}

//// LIST FUNCTION

func listWizSecurityCategories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_category.listWizSecurityCategories", "connection_error", err)
		return nil, err
	}

	options := &api.ListSecurityCategoriesRequestConfiguration{}

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
	if d.EqualsQualString("framework_id") != "" {
		options.FrameworkId = d.EqualsQualString("framework_id")
	}

	for {
		query, err := api.ListSecurityCategories(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_security_category.listWizSecurityCategories", "query_error", err)
			return nil, err
		}

		for _, category := range query.SecurityCategories.Nodes {
			d.StreamListItem(ctx, category)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.SecurityCategories.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.SecurityCategories.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCategory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_category.getSecurityCategory", "connection_error", err)
		return nil, err
	}

	query, err := api.GetSecurityCategory(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_security_category.getSecurityCategory", "query_error", err)
		return nil, err
	}

	return query.SecurityCategory, nil
}
