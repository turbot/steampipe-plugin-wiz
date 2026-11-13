package wiz

import (
	"context"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

// Filter values are only sent to the API when they look valid. A qual whose
// value is only known at execution time (for example from a subquery) can
// reach the plugin as a placeholder, and Postgres filters the rows anyway.
var (
	sbomResourceTypePattern = regexp.MustCompile(`^[A-Z][A-Z_]*$`)
	sbomResourceIdPattern   = regexp.MustCompile(`^[0-9a-fA-F-]{36}$`)
)

//// TABLE DEFINITION

func tableWizSBOMResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_sbom_resource",
		Description: "Wiz SBOM Resource",
		List: &plugin.ListConfig{
			Hydrate: listWizSBOMResources,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "id", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the resource, such as owner/repository/branch for a repository branch.", Transform: transform.FromField("Resource.Name")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The Wiz identifier of the resource.", Transform: transform.FromField("Resource.Id")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the resource, such as REPOSITORY_BRANCH, CONTAINER_IMAGE, CONTAINER, VIRTUAL_MACHINE or SERVERLESS.", Transform: transform.FromField("Resource.Type")},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "The identifier of the resource in its provider.", Transform: transform.FromField("Resource.ExternalId")},
			{Name: "native_type", Type: proto.ColumnType_STRING, Description: "The resource type as named by its provider.", Transform: transform.FromField("Resource.NativeType")},
			{Name: "cloud_platform", Type: proto.ColumnType_STRING, Description: "The platform hosting the resource, such as AWS, GCP, Kubernetes or GitHub.", Transform: transform.FromField("Resource.CloudPlatform")},
			{Name: "artifact_count", Type: proto.ColumnType_INT, Description: "The number of artifacts (packages, libraries, technologies) in the SBOM of the resource.", Transform: transform.FromField("Analytics.Artifacts.Count")},
			{Name: "version_count", Type: proto.ColumnType_INT, Description: "The number of distinct artifact versions in the SBOM of the resource.", Transform: transform.FromField("Analytics.Versions.Count")},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "The projects the resource belongs to.", Transform: transform.FromField("Resource.Projects")},
		},
	}
}

//// LIST FUNCTION

func listWizSBOMResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_sbom_resource.listWizSBOMResources", "connection_error", err)
		return nil, err
	}

	options := &api.ListSBOMResourcesRequestConfiguration{}
	if id := d.EqualsQualString("id"); sbomResourceIdPattern.MatchString(id) {
		options.ResourceId = id
	}
	if resourceType := d.EqualsQualString("type"); sbomResourceTypePattern.MatchString(resourceType) {
		options.ResourceType = resourceType
	}

	pageLimit := 500
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	for {
		query, err := api.ListSBOMResources(ctx, conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_sbom_resource.listWizSBOMResources", "query_error", err)
			return nil, err
		}

		for _, resource := range query.SBOMArtifactsGroupedByValues.Nodes {
			if resource.Resource == nil {
				continue
			}
			d.StreamListItem(ctx, resource)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.SBOMArtifactsGroupedByValues.PageInfo.HasNextPage {
			break
		}
		options.EndCursor = query.SBOMArtifactsGroupedByValues.PageInfo.EndCursor
	}

	return nil, nil
}
