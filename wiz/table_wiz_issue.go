package wiz

import (
	"context"
	"slices"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizIssue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_issue",
		Description: "Wiz Issue",
		List: &plugin.ListConfig{
			Hydrate: listWizIssues,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "control_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "created_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
				{Name: "framework_category_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "has_service_ticket", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "project_id", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "resolution_reason", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizIssue,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the issue."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current status of the issue. Possible values are: OPEN, IN_PROGRESS, RESOLVED, REJECTED."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The control severity. Possible values are: INFORMATIONAL, LOW, MEDIUM, HIGH, CRITICAL."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the issue."},
			{Name: "due_at", Type: proto.ColumnType_TIMESTAMP, Description: "The issue due date."},
			{Name: "rejection_expired_at", Type: proto.ColumnType_TIMESTAMP, Description: "The issue rejection expired at date."},
			{Name: "resolution_reason", Type: proto.ColumnType_STRING, Description: "The reason for issue resolution. Possible values are: OBJECT_DELETED, ISSUE_FIXED, CONTROL_CHANGED, CONTROL_DISABLED, CONTROL_DELETED, FALSE_POSITIVE, EXCEPTION, WONT_FIX."},
			{Name: "resolved_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was resolved."},
			{Name: "status_changed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue status was last changed."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was last updated."},
			{Name: "control_id", Type: proto.ColumnType_STRING, Description: "The control ID through which this issue is generated. As a filter, use a single value with the = operator.", Transform: transform.FromField("SourceRules").Transform(issueControlId)},
			{Name: "framework_category_id", Type: proto.ColumnType_STRING, Description: "Filter-only column: use framework_category_id = '<id>' to return the issues of a framework category. List filters (in) on this column are not supported.", Transform: transform.FromValue().TransformP(scalarQualValue, "framework_category_id")},
			{Name: "entity", Type: proto.ColumnType_JSON, Description: "The graph entity to which this issue is related."},
			{Name: "notes", Type: proto.ColumnType_JSON, Description: "The issue related notes."},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "A list of projects to which the issue is related."},
			{Name: "service_tickets", Type: proto.ColumnType_JSON, Description: "The service tickets (e.g. Jira) linked to the issue, with their integration."},
			{Name: "has_service_ticket", Type: proto.ColumnType_BOOL, Description: "True, if a service ticket is linked to the issue.", Transform: transform.FromField("ServiceTickets").Transform(issueHasServiceTicket)},
			{Name: "source_rules", Type: proto.ColumnType_JSON, Description: "The rules (controls, cloud configuration rules or cloud event rules) that generated the issue."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the issue. Possible values are: TOXIC_COMBINATION, THREAT_DETECTION, CLOUD_CONFIGURATION, ATTACK_SURFACE, RISK_TOXIC_COMBINATION."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the issue in the Wiz portal."},
			{Name: "first_event_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time of the first event of the issue, for threat detection issues."},
			{Name: "last_event_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time of the last event of the issue, for threat detection issues."},
			{Name: "resolution_note", Type: proto.ColumnType_STRING, Description: "The note added when the issue was resolved or rejected."},
			{Name: "open_reason", Type: proto.ColumnType_STRING, Description: "The reason the issue was opened."},
			{Name: "has_code_remediation", Type: proto.ColumnType_BOOL, Description: "True, if a code remediation is available for the issue.", Transform: transform.FromField("HasCodeRemediation")},
			{Name: "has_code_to_cloud_mapping", Type: proto.ColumnType_BOOL, Description: "True, if the issue's resource is mapped to the code that deploys it.", Transform: transform.FromField("HasCodeToCloudMapping")},
			{Name: "validated_as_exploitable", Type: proto.ColumnType_BOOL, Description: "True, if the issue was validated as exploitable.", Transform: transform.FromField("ValidatedAsExploitable")},
			{Name: "environments", Type: proto.ColumnType_JSON, Description: "The environments of the issue's resource, e.g. PRODUCTION, STAGING."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "Filter-only column: use project_id = '<id>' or project_id in (...) to return the issues of projects. Returns the filter value that matches one of the issue's projects.", Transform: transform.FromField("Projects").Transform(issueProjectId)},
			{Name: "entity_snapshot", Type: proto.ColumnType_JSON, Description: "A snapshot of the resource the issue is about."},
			{Name: "entity_type", Type: proto.ColumnType_STRING, Description: "The graph entity type of the resource, e.g. VIRTUAL_MACHINE, BUCKET.", Transform: transform.FromField("EntitySnapshot.Type")},
			{Name: "entity_name", Type: proto.ColumnType_STRING, Description: "The name of the resource.", Transform: transform.FromField("EntitySnapshot.Name")},
			{Name: "entity_native_type", Type: proto.ColumnType_STRING, Description: "The cloud provider type of the resource.", Transform: transform.FromField("EntitySnapshot.NativeType")},
			{Name: "entity_cloud_platform", Type: proto.ColumnType_STRING, Description: "The cloud platform of the resource.", Transform: transform.FromField("EntitySnapshot.CloudPlatform")},
			{Name: "entity_external_id", Type: proto.ColumnType_STRING, Description: "The external ID of the resource.", Transform: transform.FromField("EntitySnapshot.ExternalId")},
			{Name: "entity_provider_id", Type: proto.ColumnType_STRING, Description: "The provider unique ID of the resource, e.g. an ARN.", Transform: transform.FromField("EntitySnapshot.ProviderId")},
			{Name: "entity_region", Type: proto.ColumnType_STRING, Description: "The region of the resource.", Transform: transform.FromField("EntitySnapshot.Region")},
			{Name: "entity_subscription_external_id", Type: proto.ColumnType_STRING, Description: "The external ID of the cloud account of the resource.", Transform: transform.FromField("EntitySnapshot.SubscriptionExternalId")},
			{Name: "entity_subscription_name", Type: proto.ColumnType_STRING, Description: "The name of the cloud account of the resource.", Transform: transform.FromField("EntitySnapshot.SubscriptionName")},
			{Name: "entity_tags", Type: proto.ColumnType_JSON, Description: "The tags of the resource.", Transform: transform.FromField("EntitySnapshot.Tags")},
			{Name: "kubernetes_cluster_name", Type: proto.ColumnType_STRING, Description: "The Kubernetes cluster of the resource, if any.", Transform: transform.FromField("EntitySnapshot.KubernetesClusterName")},
			{Name: "kubernetes_namespace_name", Type: proto.ColumnType_STRING, Description: "The Kubernetes namespace of the resource, if any.", Transform: transform.FromField("EntitySnapshot.KubernetesNamespaceName")},
			{Name: "application_services", Type: proto.ColumnType_JSON, Description: "The application services the issue belongs to. Selecting this column makes the query slower."},
			{Name: "related_application_services", Type: proto.ColumnType_JSON, Description: "The application services related to the issue. Selecting this column makes the query slower."},
			{Name: "risk_factors", Type: proto.ColumnType_JSON, Description: "The risk factors of the issue (policy and severity). Selecting this column makes the query slower."},
		},
	}
}

//// LIST FUNCTION

func listWizIssues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_issue.listWizIssues", "connection_error", err)
		return nil, err
	}

	options := &api.ListIssuesRequestConfiguration{}

	options.IncludeExtras = columnRequested(d, issueExtraColumns...)

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
	options.Filter = &api.IssueFilter{}
	options.Filter.SourceControl = stringListQual(d, "control_id")
	options.Filter.FrameworkCategory = stringListQual(d, "framework_category_id")
	options.Filter.ResolutionReason = stringListQual(d, "resolution_reason")
	options.Filter.HasServiceTicket = boolQual(d, "has_service_ticket")
	options.Filter.ProjectId = stringListQual(d, "project_id")
	options.Filter.Severity = stringListQual(d, "severity")
	options.Filter.Status = stringListQual(d, "status")
	options.Filter.Type = stringListQual(d, "type")

	// Filter using time range
	options.Filter.CreatedAt = timeRangeQual(d, "created_at")

	for {
		query, err := api.ListIssues(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_issue.listWizIssues", "query_error", err)
			return nil, err
		}

		for _, issue := range query.Issues.Nodes {
			d.StreamListItem(ctx, issue)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Issues.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Issues.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWizIssue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_issue.getWizIssue", "connection_error", err)
		return nil, err
	}

	query, err := api.GetIssue(context.Background(), conn, id, columnRequested(d, issueExtraColumns...))
	if err != nil {
		plugin.Logger(ctx).Error("wiz_issue.getWizIssue", "query_error", err)
		return nil, err
	}

	return query.Issue, nil
}

// issueExtraColumns are the columns that need the (slower) extended issue query
var issueExtraColumns = []string{"application_services", "related_application_services", "risk_factors"}

//// TRANSFORM FUNCTIONS

// issueHasServiceTicket reports whether any service ticket is linked to the issue
func issueHasServiceTicket(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tickets, ok := d.Value.([]api.IssueServiceTickets)
	if !ok {
		return false, nil
	}
	return len(tickets) > 0, nil
}

// issueControlId returns the source rule ID matching a control_id qual, so the Postgres
// re-check keeps the rows the API matched; otherwise the ID of the first control.
func issueControlId(_ context.Context, d *transform.TransformData) (interface{}, error) {
	rules, ok := d.Value.([]api.IssueSourceRule)
	if !ok {
		return nil, nil
	}
	wanted := transformQualValues(d, "control_id")
	for _, rule := range rules {
		if slices.Contains(wanted, rule.Id) {
			return rule.Id, nil
		}
	}
	for _, rule := range rules {
		if rule.Typename == "Control" {
			return rule.Id, nil
		}
	}
	return nil, nil
}

// issueProjectId returns the project_id qual value that matches one of the issue's projects,
// so the Postgres re-check keeps the rows the API matched. Nil when project_id is not filtered.
func issueProjectId(_ context.Context, d *transform.TransformData) (interface{}, error) {
	projects, ok := d.Value.([]api.IssueQueryProject)
	if !ok {
		return nil, nil
	}
	wanted := transformQualValues(d, "project_id")
	for _, project := range projects {
		if slices.Contains(wanted, project.Id) {
			return project.Id, nil
		}
	}
	return nil, nil
}
