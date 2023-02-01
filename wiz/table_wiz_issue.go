package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
				{Name: "resolution_reason", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getWizIssue,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the issue."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current status of the issue."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The control severity."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the issue."},
			{Name: "due_at", Type: proto.ColumnType_TIMESTAMP, Description: "The issue due date."},
			{Name: "rejection_expired_at", Type: proto.ColumnType_TIMESTAMP, Description: "The issue rejection expired at date."},
			{Name: "resolution_reason", Type: proto.ColumnType_STRING, Description: "The reason for issue resolution."},
			{Name: "resolved_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was resolved."},
			{Name: "status_changed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue status was last changed."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the issue was last updated."},
			{Name: "control", Type: proto.ColumnType_JSON, Description: "A list of security controls that the issue violates."},
			{Name: "notes", Type: proto.ColumnType_JSON, Description: "The issue related notes."},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "A list of projects to which the issue is related."},
			{Name: "service_tickets", Type: proto.ColumnType_JSON, Description: "Specifies the related issues from ticket services."},
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

	// Default set to 60.
	// This is the maximum number of items can be requested.
	pageLimit := 500

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	//
	// Check for additional filters
	// TODO :: Add other supported filters: createdAt, dueAt, resolvedAt, statusChangedAt
	//
	options.Filter = &api.IssueFilter{}
	if d.EqualsQualString("resolution_reason") != "" {
		options.Filter.ResolutionReason = d.EqualsQualString("resolution_reason")
	}
	if d.EqualsQualString("severity") != "" {
		options.Filter.Severity = d.EqualsQualString("severity")
	}
	if d.EqualsQualString("status") != "" {
		options.Filter.Status = d.EqualsQualString("status")
	}

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
		plugin.Logger(ctx).Error("wiz_issue.getIssue", "connection_error", err)
		return nil, err
	}

	query, err := api.GetIssue(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_issue.getIssue", "query_error", err)
		return nil, err
	}

	return query.Issue, nil
}
