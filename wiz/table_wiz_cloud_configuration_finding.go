package wiz

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizCloudConfigurationFinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_cloud_configuration_finding",
		Description: "Wiz Cloud Configuration Finding",
		List: &plugin.ListConfig{
			Hydrate: listWizCloudConfigurationFindings,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "result", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "rule_id", Require: plugin.Optional},
				{Name: "analyzed_at", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getCloudConfigurationFinding,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The name of the resource.", Transform: transform.FromField("Resource.Name")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the finding."},
			{Name: "result", Type: proto.ColumnType_STRING, Description: "The outcome of the finding. Possible values are: ERROR, FAIL, NOT_ASSESSED, PASS."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "The finding severity. Possible values: CRITICAL, HIGH, LOW, MEDIUM, NONE."},
			{Name: "analyzed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the finding was detected."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the finding. Possible values: IN_PROGRESS, OPEN, REJECTED, RESOLVED."},
			{Name: "remediation", Type: proto.ColumnType_STRING, Description: "Specifies the steps to mitigate the issue that match this rule."},
			{Name: "resolution_reason", Type: proto.ColumnType_STRING, Description: "The status resolution reason of the finding."},
			{Name: "rule_id", Type: proto.ColumnType_STRING, Description: "Specifies the rule against which the finding is generated.", Transform: transform.FromField("Rule.Id")},
			{Name: "resource", Type: proto.ColumnType_JSON, Description: "Specifies the configuration of the resource detected through the finding."},
			{Name: "subscription", Type: proto.ColumnType_JSON, Description: "Specifies the cloud account where the rule was applied and the finding is generated."},
		},
	}
}

//// LIST FUNCTION

func listWizCloudConfigurationFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_configuration_finding.listWizCloudConfigurationFindings", "connection_error", err)
		return nil, err
	}

	options := &api.ListConfigurationFindingsRequestConfiguration{}

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
	if d.EqualsQualString("result") != "" {
		options.Result = d.EqualsQualString("result")
	}
	if d.EqualsQualString("rule_id") != "" {
		options.RuleId = d.EqualsQualString("rule_id")
	}
	if d.EqualsQualString("severity") != "" {
		options.Severity = d.EqualsQualString("severity")
	}
	if d.EqualsQualString("status") != "" {
		options.Status = d.EqualsQualString("status")
	}

	// Filter using time range
	if d.Quals["analyzed_at"] != nil {
		options.AnalyzedAt = api.ConfigurationFindingDateFilter{}
		for _, q := range d.Quals["analyzed_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=":
				options.AnalyzedAt.After = givenTime.Add(-2 * time.Second).Format(time.RFC3339)
				options.AnalyzedAt.Before = givenTime.Add(2 * time.Second).Format(time.RFC3339)
			case ">=":
				options.AnalyzedAt.After = givenTime.Add(-1 * time.Second).Format(time.RFC3339)
			case ">":
				options.AnalyzedAt.After = givenTime.Format(time.RFC3339)
			case "<=":
				options.AnalyzedAt.Before = givenTime.Add(1 * time.Second).Format(time.RFC3339)
			case "<":
				options.AnalyzedAt.Before = givenTime.Format(time.RFC3339)
			}
		}
	}

	for {
		query, err := api.ListConfigurationFindings(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_cloud_configuration_finding.listWizCloudConfigurationFindings", "query_error", err)
			return nil, err
		}

		for _, finding := range query.ConfigurationFindings.Nodes {
			d.StreamListItem(ctx, finding)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.ConfigurationFindings.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.ConfigurationFindings.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudConfigurationFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Check the quals
	id := d.EqualsQualString("id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_configuration_finding.getCloudConfigurationFinding", "connection_error", err)
		return nil, err
	}

	query, err := api.GetConfigurationFinding(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_cloud_configuration_finding.getCloudConfigurationFinding", "query_error", err)
		return nil, err
	}

	return query.ConfigurationFinding, nil
}
