package wiz

import (
	"context"
	"slices"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

// boolQual returns the value of a boolean "=" or "<>" qual, or nil if unset
func boolQual(d *plugin.QueryData, column string) *bool {
	if d.Quals[column] == nil {
		return nil
	}
	for _, q := range d.Quals[column].Quals {
		value := q.Value.GetBoolValue()
		switch q.Operator {
		case "=":
			return &value
		case "<>":
			value = !value
			return &value
		}
	}
	return nil
}

// stringListQual returns the values of an "=" or "in" qual on a string column
func stringListQual(d *plugin.QueryData, column string) []string {
	qual := d.EqualsQuals[column]
	if qual == nil {
		return nil
	}
	if list := qual.GetListValue(); list != nil {
		values := []string{}
		for _, v := range list.Values {
			if s := v.GetStringValue(); s != "" {
				values = append(values, s)
			}
		}
		return values
	}
	if s := qual.GetStringValue(); s != "" {
		return []string{s}
	}
	return nil
}

// singleStringQual returns the value of an "=" qual for API filters that take a single value.
// An "in" list with several values is not pushed down; Postgres filters those rows itself.
func singleStringQual(d *plugin.QueryData, column string) string {
	if values := stringListQual(d, column); len(values) == 1 {
		return values[0]
	}
	return ""
}

// findingStatusQual returns the status qual values of a findings table without IN_PROGRESS,
// which the findings APIs treat as OPEN. onlyInProgress is true when the qual asked for
// IN_PROGRESS alone, so the caller can return no rows without calling the API.
func findingStatusQual(d *plugin.QueryData) (statuses []string, onlyInProgress bool) {
	values := stringListQual(d, "status")
	if len(values) == 0 {
		return nil, false
	}
	for _, v := range values {
		if v != "IN_PROGRESS" {
			statuses = append(statuses, v)
		}
	}
	return statuses, len(statuses) == 0
}

// timeRangeQual converts comparison quals on a timestamp column into an API date filter.
// "in" quals are not pushed down; Postgres filters those rows itself.
func timeRangeQual(d *plugin.QueryData, column string) *api.DateFilter {
	if d.Quals[column] == nil {
		return nil
	}
	filter := &api.DateFilter{}
	for _, q := range d.Quals[column].Quals {
		if q.Value.GetListValue() != nil {
			return nil
		}
		givenTime := q.Value.GetTimestampValue().AsTime()
		switch q.Operator {
		case "=":
			filter.After = givenTime.Add(-2 * time.Second).Format(time.RFC3339Nano)
			filter.Before = givenTime.Add(2 * time.Second).Format(time.RFC3339Nano)
		case ">=":
			filter.After = givenTime.Add(-1 * time.Second).Format(time.RFC3339Nano)
		case ">":
			filter.After = givenTime.Format(time.RFC3339Nano)
		case "<=":
			filter.Before = givenTime.Add(1 * time.Second).Format(time.RFC3339Nano)
		case "<":
			filter.Before = givenTime.Format(time.RFC3339Nano)
		}
	}
	if filter.After == "" && filter.Before == "" {
		return nil
	}
	return filter
}

// columnRequested reports whether the query selects any of the given columns
func columnRequested(d *plugin.QueryData, columns ...string) bool {
	for _, c := range columns {
		if slices.Contains(d.QueryContext.Columns, c) {
			return true
		}
	}
	return false
}

// serviceTicketColumns are the columns that need the (slower) service ticket query
var serviceTicketColumns = []string{"service_tickets", "has_service_ticket"}

// Service ticket as returned in the service_tickets column
type serviceTicketRow struct {
	ExternalId  string `json:"external_id"`
	Id          string `json:"id"`
	Integration string `json:"integration,omitempty"`
	Name        string `json:"name"`
	Url         string `json:"url"`
}

type linkedServiceTickets interface {
	LinkedServiceTickets() []api.ServiceTicket
}

// findingServiceTickets returns the deduplicated service tickets linked to a finding
func findingServiceTickets(_ context.Context, d *transform.TransformData) (interface{}, error) {
	finding, ok := d.HydrateItem.(linkedServiceTickets)
	if !ok {
		return nil, nil
	}
	rows := []serviceTicketRow{}
	for _, t := range finding.LinkedServiceTickets() {
		row := serviceTicketRow{ExternalId: t.ExternalId, Id: t.Id, Name: t.Name, Url: t.Url}
		if t.Integration != nil {
			row.Integration = t.Integration.Name
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// findingHasServiceTicket reports whether any service ticket is linked to a finding
func findingHasServiceTicket(_ context.Context, d *transform.TransformData) (interface{}, error) {
	finding, ok := d.HydrateItem.(linkedServiceTickets)
	if !ok {
		return nil, nil
	}
	return len(finding.LinkedServiceTickets()) > 0, nil
}
