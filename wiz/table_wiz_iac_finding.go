package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizIaCFinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_iac_finding",
		Description: "Wiz IaC Finding",
		List: &plugin.ListConfig{
			Hydrate: listWizIaCFindings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "file_path", Require: plugin.Optional},
				{Name: "is_default_branch", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "platform", Require: plugin.Optional},
				{Name: "repository_id", Require: plugin.Optional},
				{Name: "rule_id", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "updated_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizIaCFinding,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the IaC finding."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the IaC finding."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The severity of the finding. Possible values are: CRITICAL, HIGH, MEDIUM, LOW, INFORMATIONAL."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the finding. Possible values are: OPEN, RESOLVED, REJECTED. All statuses are returned unless filtered. The API treats IN_PROGRESS as OPEN for findings, so it is never returned."},
			{Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "The first time the finding was seen."},
			{Name: "last_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the finding was seen."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the finding was last updated."},
			{Name: "platform", Type: proto.ColumnType_STRING, Description: "The IaC platform, e.g. TERRAFORM, CLOUD_FORMATION, KUBERNETES, DOCKERFILE, GITHUB_ACTIONS."},
			{Name: "file_path", Type: proto.ColumnType_STRING, Description: "The path of the file where the finding was detected."},
			{Name: "start_line", Type: proto.ColumnType_INT, Description: "The first line of the finding in the file."},
			{Name: "end_line", Type: proto.ColumnType_INT, Description: "The last line of the finding in the file."},
			{Name: "file_url", Type: proto.ColumnType_STRING, Description: "The URL of the file in the version control system.", Transform: transform.FromField("FileURL")},
			{Name: "vcs_platform", Type: proto.ColumnType_STRING, Description: "The version control platform of the repository."},
			{Name: "cloud_platform", Type: proto.ColumnType_STRING, Description: "The cloud platform of the resource defined in the file."},
			{Name: "entity_type", Type: proto.ColumnType_STRING, Description: "The graph entity type of the resource defined in the file."},
			{Name: "native_type", Type: proto.ColumnType_STRING, Description: "The native type of the resource defined in the file."},
			{Name: "found_content", Type: proto.ColumnType_STRING, Description: "The content found in the file."},
			{Name: "expected_content", Type: proto.ColumnType_STRING, Description: "The content expected by the rule."},
			{Name: "match_content", Type: proto.ColumnType_STRING, Description: "The content that matched the rule."},
			{Name: "remediation_instructions", Type: proto.ColumnType_STRING, Description: "Instructions for remediating the finding."},
			{Name: "repository_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the repository.", Transform: transform.FromField("Repository.Id")},
			{Name: "repository_name", Type: proto.ColumnType_STRING, Description: "The name of the repository.", Transform: transform.FromField("Repository.Name")},
			{Name: "branch_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the repository branch.", Transform: transform.FromField("Branch.Id")},
			{Name: "branch_name", Type: proto.ColumnType_STRING, Description: "The name of the repository branch.", Transform: transform.FromField("Branch.Name")},
			{Name: "is_default_branch", Type: proto.ColumnType_BOOL, Description: "Filter-only column: use is_default_branch = true (or false) to return findings on (or off) the default branch of a repository. Only the = operator is supported.", Transform: transform.FromQual("is_default_branch")},
			{Name: "rule_id", Type: proto.ColumnType_STRING, Description: "The ID of the cloud configuration rule.", Transform: transform.FromField("Rule.Id")},
			{Name: "rule", Type: proto.ColumnType_JSON, Description: "The cloud configuration rule that generated the finding."},
			{Name: "resource_graph_entity", Type: proto.ColumnType_JSON, Description: "The graph entity the finding maps to."},
			{Name: "vcs_details", Type: proto.ColumnType_JSON, Description: "Version control details (code owners) of the finding."},
		},
	}
}

//// LIST FUNCTION

func listWizIaCFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// The findings APIs treat IN_PROGRESS as OPEN, so a filter on IN_PROGRESS alone matches nothing
	status, onlyInProgress := findingStatusQual(d)
	if onlyInProgress {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_iac_finding.listWizIaCFindings", "connection_error", err)
		return nil, err
	}

	options := &api.ListIaCFindingsRequestConfiguration{}

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
	options.FilePath = stringListQual(d, "file_path")
	options.IsDefaultBranch = boolQual(d, "is_default_branch")
	options.Platform = stringListQual(d, "platform")
	options.RepositoryId = stringListQual(d, "repository_id")
	options.RuleId = stringListQual(d, "rule_id")
	options.Severity = stringListQual(d, "severity")
	options.Status = status
	options.UpdatedAt = timeRangeQual(d, "updated_at")

	for {
		query, err := api.ListIaCFindings(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_iac_finding.listWizIaCFindings", "query_error", err)
			return nil, err
		}

		for _, finding := range query.IaCFindings.Nodes {
			d.StreamListItem(ctx, finding)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.IaCFindings.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.IaCFindings.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizIaCFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_iac_finding.getWizIaCFinding", "connection_error", err)
		return nil, err
	}

	query, err := api.GetIaCFinding(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_iac_finding.getWizIaCFinding", "query_error", err)
		return nil, err
	}

	// The API returns null for an unknown ID
	if query.IaCFinding == nil {
		return nil, nil
	}

	return *query.IaCFinding, nil
}
