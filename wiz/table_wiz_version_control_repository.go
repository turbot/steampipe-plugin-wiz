package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizVersionControlRepository(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_version_control_repository",
		Description: "Wiz Version Control Repository",
		List: &plugin.ListConfig{
			Hydrate: listWizVersionControlRepositories,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "archived", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "organization_id", Require: plugin.Optional},
				{Name: "platform", Require: plugin.Optional},
				{Name: "public", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "visibility", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizVersionControlRepository,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the repository."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the organization."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the repository."},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "The version control platform, e.g. GITHUB, GITLAB, AZURE_DEVOPS, BITBUCKET."},
			{Name: "public", Type: proto.ColumnType_BOOL, Description: "True, if the repository is public.", Transform: transform.FromField("Public")},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "True, if the repository is archived.", Transform: transform.FromField("Archived")},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the repository."},
			{Name: "default_branch", Type: proto.ColumnType_STRING, Description: "The name of the default branch."},
			{Name: "organization_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the organization owning the repository.", Transform: transform.FromField("Organization.Id")},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization owning the repository.", Transform: transform.FromField("Organization.Name")},
			{Name: "organization", Type: proto.ColumnType_JSON, Description: "The organization owning the repository."},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "The projects the repository belongs to."},
		},
	}
}

//// LIST FUNCTION

func listWizVersionControlRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_version_control_repository.listWizVersionControlRepositories", "connection_error", err)
		return nil, err
	}

	options := &api.ListVersionControlRepositoriesRequestConfiguration{}

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
	options.Archived = boolQual(d, "archived")
	options.OrganizationId = stringListQual(d, "organization_id")
	options.Platform = stringListQual(d, "platform")
	options.Public = boolQual(d, "public")
	options.Visibility = stringListQual(d, "visibility")

	for {
		query, err := api.ListVersionControlRepositories(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_version_control_repository.listWizVersionControlRepositories", "query_error", err)
			return nil, err
		}

		for _, repository := range query.Repositories.Nodes {
			d.StreamListItem(ctx, repository)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Repositories.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Repositories.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizVersionControlRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_version_control_repository.getWizVersionControlRepository", "connection_error", err)
		return nil, err
	}

	query, err := api.GetVersionControlRepository(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_version_control_repository.getWizVersionControlRepository", "query_error", err)
		return nil, err
	}

	// The API returns null for an unknown ID
	if query.Repository == nil {
		return nil, nil
	}

	return *query.Repository, nil
}
