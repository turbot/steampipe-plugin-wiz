package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizSASTFinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_sast_finding",
		Description: "Wiz SAST Finding",
		List: &plugin.ListConfig{
			Hydrate: listWizSASTFindings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "ai_verdict", Require: plugin.Optional},
				{Name: "created_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
				{Name: "file_path", Require: plugin.Optional},
				{Name: "is_ai_powered", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "is_default_branch", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "origin", Require: plugin.Optional},
				{Name: "repository_branch_id", Require: plugin.Optional},
				{Name: "repository_id", Require: plugin.Optional},
				{Name: "rule_id", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "updated_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizSASTFinding,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the SAST finding."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the SAST finding."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the SAST finding."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The severity of the finding. Possible values are: CRITICAL, HIGH, MEDIUM, LOW, INFORMATIONAL."},
			{Name: "original_severity", Type: proto.ColumnType_STRING, Description: "The severity reported by the scanner, before any override."},
			{Name: "impact", Type: proto.ColumnType_STRING, Description: "The impact of the weakness."},
			{Name: "likelihood", Type: proto.ColumnType_STRING, Description: "The likelihood of the weakness being exploited."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the finding. Possible values are: OPEN, RESOLVED, REJECTED. All statuses are returned unless filtered. The API treats IN_PROGRESS as OPEN for findings, so it is never returned."},
			{Name: "resolution_reason", Type: proto.ColumnType_STRING, Description: "The reason the finding was resolved or rejected."},
			{Name: "rejection_expired_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the rejection of the finding expires."},
			{Name: "origin", Type: proto.ColumnType_STRING, Description: "The source that reported the finding, e.g. WIZ or a third-party scanner."},
			{Name: "is_ai_powered", Type: proto.ColumnType_BOOL, Description: "True, if the finding was detected by AI-powered scanning.", Transform: transform.FromField("IsAiPowered")},
			{Name: "ai_verdict", Type: proto.ColumnType_STRING, Description: "The verdict of the Wiz AI analysis. Possible values are: INCONCLUSIVE, EXPLOITABLE, TEST_CODE, SECURITY_CHALLENGE, NON_EXPLOITABLE, NO_RISK.", Transform: transform.FromField("AiAnalysis.Verdict")},
			{Name: "ai_analysis", Type: proto.ColumnType_JSON, Description: "The Wiz AI analysis of the finding (verdict, likelihood, confidence, explanation)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the finding was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the finding was last updated."},
			{Name: "first_detected_at_source", Type: proto.ColumnType_TIMESTAMP, Description: "The first time the source scanner detected the finding."},
			{Name: "file_path", Type: proto.ColumnType_STRING, Description: "The path of the file where the finding was detected."},
			{Name: "start_line", Type: proto.ColumnType_INT, Description: "The first line of the finding in the file."},
			{Name: "start_column", Type: proto.ColumnType_INT, Description: "The first column of the finding in the file."},
			{Name: "end_line", Type: proto.ColumnType_INT, Description: "The last line of the finding in the file."},
			{Name: "end_column", Type: proto.ColumnType_INT, Description: "The last column of the finding in the file."},
			{Name: "snippet", Type: proto.ColumnType_STRING, Description: "The code snippet where the finding was detected."},
			{Name: "remediation_instructions", Type: proto.ColumnType_STRING, Description: "Instructions for remediating the finding."},
			{Name: "source_url", Type: proto.ColumnType_STRING, Description: "The URL of the finding in the source scanner, for third-party findings."},
			{Name: "code_library_language", Type: proto.ColumnType_JSON, Description: "The programming languages of the finding."},
			{Name: "repository_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the repository.", Transform: transform.FromField("Repository.Id")},
			{Name: "repository_name", Type: proto.ColumnType_STRING, Description: "The name of the repository.", Transform: transform.FromField("Repository.Name")},
			{Name: "repository_platform", Type: proto.ColumnType_STRING, Description: "The version control platform of the repository.", Transform: transform.FromField("Repository.Platform")},
			{Name: "repository_branch_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the repository branch.", Transform: transform.FromField("RepositoryBranch.Id")},
			{Name: "repository_branch_name", Type: proto.ColumnType_STRING, Description: "The name of the repository branch.", Transform: transform.FromField("RepositoryBranch.Name")},
			{Name: "is_default_branch", Type: proto.ColumnType_BOOL, Description: "Filter-only column: use is_default_branch = true (or false) to return findings on (or off) the default branch of a repository. Only the = operator is supported.", Transform: transform.FromQual("is_default_branch")},
			{Name: "rule_id", Type: proto.ColumnType_STRING, Description: "The Wiz ID of the SAST rule.", Transform: transform.FromField("Rule.Id")},
			{Name: "rule_external_id", Type: proto.ColumnType_STRING, Description: "The external ID of the SAST rule.", Transform: transform.FromField("Rule.ExternalId")},
			{Name: "rule", Type: proto.ColumnType_JSON, Description: "The SAST rule that generated the finding."},
			{Name: "weaknesses", Type: proto.ColumnType_JSON, Description: "The weaknesses (e.g. CWE) of the finding."},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "The projects the finding belongs to."},
			{Name: "vcs_details", Type: proto.ColumnType_JSON, Description: "Version control details of the finding (commit and code owners)."},
			{Name: "related_issue_analytics", Type: proto.ColumnType_JSON, Description: "Counts of related issues by severity."},
			{Name: "service_tickets", Type: proto.ColumnType_JSON, Description: "The service tickets (e.g. Jira) linked to the finding through its posture issues and its first 50 related issues, deduplicated. Selecting this column makes the query slower and caps the page size at 500.", Transform: transform.FromValue().Transform(findingServiceTickets)},
			{Name: "has_service_ticket", Type: proto.ColumnType_BOOL, Description: "True, if a service ticket is linked to the finding through its posture issues or its first 50 related issues. Selecting this column makes the query slower and caps the page size at 500.", Transform: transform.FromValue().Transform(findingHasServiceTicket)},
		},
	}
}

//// LIST FUNCTION

func listWizSASTFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// The findings APIs treat IN_PROGRESS as OPEN, so a filter on IN_PROGRESS alone matches nothing
	status, onlyInProgress := findingStatusQual(d)
	if onlyInProgress {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_sast_finding.listWizSASTFindings", "connection_error", err)
		return nil, err
	}

	options := &api.ListSASTFindingsRequestConfiguration{}

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
	options.AiVerdict = stringListQual(d, "ai_verdict")
	options.CreatedAt = timeRangeQual(d, "created_at")
	options.FilePath = stringListQual(d, "file_path")
	options.IsAiPowered = boolQual(d, "is_ai_powered")
	options.IsDefaultBranch = boolQual(d, "is_default_branch")
	options.Origin = stringListQual(d, "origin")
	options.RepositoryBranchId = stringListQual(d, "repository_branch_id")
	options.RepositoryId = stringListQual(d, "repository_id")
	options.IncludeServiceTickets = columnRequested(d, serviceTicketColumns...)
	options.RuleId = stringListQual(d, "rule_id")
	options.Severity = stringListQual(d, "severity")
	options.Status = status
	options.UpdatedAt = timeRangeQual(d, "updated_at")

	for {
		query, err := api.ListSASTFindings(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_sast_finding.listWizSASTFindings", "query_error", err)
			return nil, err
		}

		for _, finding := range query.SASTFindings.Nodes {
			d.StreamListItem(ctx, finding)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.SASTFindings.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.SASTFindings.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizSASTFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_sast_finding.getWizSASTFinding", "connection_error", err)
		return nil, err
	}

	query, err := api.GetSASTFinding(context.Background(), conn, id, columnRequested(d, serviceTicketColumns...))
	if err != nil {
		plugin.Logger(ctx).Error("wiz_sast_finding.getWizSASTFinding", "query_error", err)
		return nil, err
	}

	return query.SASTFinding, nil
}
