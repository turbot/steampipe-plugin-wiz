package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizProject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_project",
		Description: "Wiz Project",
		List: &plugin.ListConfig{
			Hydrate: listWizProjects,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizProject,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the project."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the project."},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "The project slug."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the project."},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "If tru, the project was archived."},
			{Name: "business_unit", Type: proto.ColumnType_STRING, Description: "The project business unit."},
			{Name: "security_score", Type: proto.ColumnType_INT, Description: "Security score is based on the number of successful assessments ran on this project out of the total assessments."},
			{Name: "cloud_account_count", Type: proto.ColumnType_INT, Description: "The count of cloud account associated with the project."},
			{Name: "cloud_organization_count", Type: proto.ColumnType_INT, Description: "The count of cloud organization associated with the project."},
			{Name: "entity_count", Type: proto.ColumnType_INT, Description: "The count of the entity."},
			{Name: "kubernetes_cluster_count", Type: proto.ColumnType_INT, Description: "The count of Kubernetes cluster created within the project."},
			{Name: "profile_completion", Type: proto.ColumnType_INT, Description: "The profile completion percentage of the project."},
			{Name: "repository_count", Type: proto.ColumnType_INT, Description: "The count of the repository."},
			{Name: "team_member_count", Type: proto.ColumnType_INT, Description: "The count of the project team member."},
			{Name: "technology_count", Type: proto.ColumnType_INT, Description: "The count of the technology."},
			{Name: "workload_count", Type: proto.ColumnType_INT, Description: "The count of the workload."},
			{Name: "identifiers", Type: proto.ColumnType_JSON, Description: "A list of project identifiers."},
			{Name: "project_owners", Type: proto.ColumnType_JSON, Description: "A list of project owners."},
			{Name: "resource_tag_links", Type: proto.ColumnType_STRING, Description: "A list of resource tags."},
			{Name: "risk_profile", Type: proto.ColumnType_STRING, Description: "Specifies the project risk-profile."},
		},
	}
}

//// LIST FUNCTION

func listWizProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_project.listWizProjects", "connection_error", err)
		return nil, err
	}

	options := &api.ListProjectsRequestConfiguration{}

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
		query, err := api.ListProjects(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_project.listWizProjects", "query_error", err)
			return nil, err
		}

		for _, project := range query.Projects.Nodes {
			d.StreamListItem(ctx, project)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Projects.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Projects.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_project.getWizProject", "connection_error", err)
		return nil, err
	}

	query, err := api.GetProject(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_project.getWizProject", "query_error", err)
		return nil, err
	}

	return query.Project, nil
}
