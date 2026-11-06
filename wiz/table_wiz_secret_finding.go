package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizSecretFinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_secret_finding",
		Description: "Wiz Secret Finding",
		List: &plugin.ListConfig{
			Hydrate: listWizSecretFindings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "first_seen_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
				{Name: "is_default_branch", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "last_updated_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
				{Name: "origin", Require: plugin.Optional},
				{Name: "path", Require: plugin.Optional},
				{Name: "repository_id", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "resource_type", Require: plugin.Optional},
				{Name: "rule_id", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "validation_status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizSecretFinding,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the secret finding."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the secret finding."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the secret, e.g. CLOUD_KEY, PASSWORD, PRIVATE_KEY, SAAS_API_KEY."},
			{Name: "confidence", Type: proto.ColumnType_STRING, Description: "The confidence of the detection."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The severity of the finding. Possible values are: CRITICAL, HIGH, MEDIUM, LOW, INFORMATIONAL."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the finding. Possible values are: OPEN, RESOLVED, REJECTED. All statuses are returned unless filtered. The API treats IN_PROGRESS as OPEN for findings, so it is never returned."},
			{Name: "resolution_reason", Type: proto.ColumnType_STRING, Description: "The reason the finding was resolved or rejected."},
			{Name: "validation_status", Type: proto.ColumnType_STRING, Description: "Whether the secret was validated as live. Possible values are: VALID, INVALID, UNKNOWN."},
			{Name: "last_validated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the secret was validated."},
			{Name: "is_encrypted", Type: proto.ColumnType_BOOL, Description: "True, if the secret is encrypted.", Transform: transform.FromField("IsEncrypted")},
			{Name: "is_managed", Type: proto.ColumnType_BOOL, Description: "True, if the secret is managed by a secret manager.", Transform: transform.FromField("IsManaged")},
			{Name: "has_leaked", Type: proto.ColumnType_BOOL, Description: "True, if the secret has leaked.", Transform: transform.FromField("HasLeaked")},
			{Name: "origin", Type: proto.ColumnType_STRING, Description: "The source that reported the finding, e.g. WIZ or a third-party scanner."},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "The external identifier of the finding."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path of the file where the secret was found."},
			{Name: "line_number", Type: proto.ColumnType_INT, Description: "The line number where the secret was found."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the finding."},
			{Name: "remediation", Type: proto.ColumnType_STRING, Description: "Instructions for remediating the finding."},
			{Name: "first_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "The first time the secret was seen."},
			{Name: "last_seen_at", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the secret was seen."},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the secret was modified."},
			{Name: "last_updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the finding was updated."},
			{Name: "resolved_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the finding was resolved."},
			{Name: "code_to_cloud_pipeline_stage", Type: proto.ColumnType_STRING, Description: "The pipeline stage where the secret was found."},
			{Name: "rule_id", Type: proto.ColumnType_STRING, Description: "The ID of the secret detection rule.", Transform: transform.FromField("Rule.Id")},
			{Name: "rule", Type: proto.ColumnType_JSON, Description: "The secret detection rule that found the secret."},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the resource holding the secret.", Transform: transform.FromField("Resource.Id")},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "The type of the resource holding the secret, e.g. REPOSITORY_BRANCH, VIRTUAL_MACHINE.", Transform: transform.FromField("Resource.Type")},
			{Name: "resource", Type: proto.ColumnType_JSON, Description: "The resource holding the secret."},
			{Name: "repository_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the code repository, for secrets found in code.", Transform: transform.FromField("Resource.TypedProperties.Repository.Id")},
			{Name: "repository_name", Type: proto.ColumnType_STRING, Description: "The name of the code repository, for secrets found in code.", Transform: transform.FromField("Resource.TypedProperties.Repository.Name")},
			{Name: "is_default_branch", Type: proto.ColumnType_BOOL, Description: "True, if the secret was found on the default branch of a repository.", Transform: transform.FromField("Resource.TypedProperties.IsDefaultBranch")},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "The projects the finding belongs to."},
			{Name: "vcs_code_owners", Type: proto.ColumnType_JSON, Description: "The code owners of the file in the version control system."},
			{Name: "vcs_details", Type: proto.ColumnType_JSON, Description: "Version control details (initial commit) of the secret."},
			{Name: "related_issue_analytics", Type: proto.ColumnType_JSON, Description: "Counts of related issues by severity."},
			{Name: "service_tickets", Type: proto.ColumnType_JSON, Description: "The service tickets (e.g. Jira) linked to the finding through its posture issues and its first 50 related issues, deduplicated. Selecting this column makes the query slower and caps the page size at 500.", Transform: transform.FromValue().Transform(findingServiceTickets)},
			{Name: "has_service_ticket", Type: proto.ColumnType_BOOL, Description: "True, if a service ticket is linked to the finding through its posture issues or its first 50 related issues. Selecting this column makes the query slower and caps the page size at 500.", Transform: transform.FromValue().Transform(findingHasServiceTicket)},
		},
	}
}

//// LIST FUNCTION

func listWizSecretFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// The findings APIs treat IN_PROGRESS as OPEN, so a filter on IN_PROGRESS alone matches nothing
	status, onlyInProgress := findingStatusQual(d)
	if onlyInProgress {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_secret_finding.listWizSecretFindings", "connection_error", err)
		return nil, err
	}

	options := &api.ListSecretFindingsRequestConfiguration{}

	// Default set to 500.
	// This is the maximum number of items can be requested, and also the cap when service tickets are selected.
	pageLimit := 500

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Check for additional filters
	options.FirstSeenAt = timeRangeQual(d, "first_seen_at")
	options.IsDefaultBranch = boolQual(d, "is_default_branch")
	options.LastUpdatedAt = timeRangeQual(d, "last_updated_at")
	options.Origin = stringListQual(d, "origin")
	options.Path = stringListQual(d, "path")
	options.RepositoryId = stringListQual(d, "repository_id")
	options.ResourceId = stringListQual(d, "resource_id")
	options.ResourceType = stringListQual(d, "resource_type")
	options.RuleId = stringListQual(d, "rule_id")
	options.Severity = stringListQual(d, "severity")
	options.Status = status
	options.Type = stringListQual(d, "type")
	options.ValidationStatus = stringListQual(d, "validation_status")
	options.IncludeServiceTickets = columnRequested(d, serviceTicketColumns...)

	for {
		query, err := api.ListSecretFindings(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_secret_finding.listWizSecretFindings", "query_error", err)
			return nil, err
		}

		for _, finding := range query.SecretFindings.Nodes {
			d.StreamListItem(ctx, finding)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.SecretFindings.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.SecretFindings.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizSecretFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_secret_finding.getWizSecretFinding", "connection_error", err)
		return nil, err
	}

	query, err := api.GetSecretFinding(context.Background(), conn, id, columnRequested(d, serviceTicketColumns...))
	if err != nil {
		plugin.Logger(ctx).Error("wiz_secret_finding.getWizSecretFinding", "query_error", err)
		return nil, err
	}

	// The API returns null for an unknown ID
	if query.SecretFinding == nil {
		return nil, nil
	}

	return *query.SecretFinding, nil
}
