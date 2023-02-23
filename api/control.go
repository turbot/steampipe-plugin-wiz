package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryControlList = `
query ListControls($first: Int, $after: String, $filter: ControlFilters) {
  controls(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      description
      type
      query
      severity
      enabled
      enabledForLBI
      enabledForMBI
      enabledForHBI
      enabledForUnattributed
      resolutionRecommendation
      createdBy {
        id
      }
      createdAt
      tags
      scopeProject {
        id
      }
      supportsNRT
      hasAutoRemediation
      sourceCloudConfigurationRule {
        id
      }
      lastRunAt
      lastRunError
      lastSuccessfulRunAt
      resolutionRecommendationPlainText
    }
  }
}
`

	queryControlGet = `
query GetControl($id: ID!) {
  control(id: $id) {
    id
    name
    type
    query
    severity
    enabled
    enabledForLBI
    enabledForMBI
    enabledForHBI
    enabledForUnattributed
    resolutionRecommendation
    createdBy {
      id
    }
    createdAt
    tags
    scopeProject {
      id
    }
    supportsNRT
    hasAutoRemediation
    sourceCloudConfigurationRule {
      id
    }
    lastRunAt
    lastRunError
    lastSuccessfulRunAt
    resolutionRecommendationPlainText
  }
}
`
)

// Control object
type Control struct {
	CreatedAt                         string                      `json:"createdAt"`
	CreatedBy                         ControlQueryUser            `json:"createdBy"`
	Description                       string                      `json:"description"`
	Enabled                           bool                        `json:"enabled"`
	EnabledForHBI                     bool                        `json:"enabledForHBI"`
	EnabledForLBI                     bool                        `json:"enabledForLBI"`
	EnabledForMBI                     bool                        `json:"enabledForMBI"`
	EnabledForUnattributed            bool                        `json:"enabledForUnattributed"`
	HasAutoRemediation                bool                        `json:"hasAutoRemediation"`
	Id                                string                      `json:"id"`
	LastRunAt                         string                      `json:"lastRunAt"`
	LastRunError                      string                      `json:"lastRunError"`
	LastSuccessfulRunAt               string                      `json:"lastSuccessfulRunAt"`
	Name                              string                      `json:"name"`
	Query                             interface{}                 `json:"query"`
	ResolutionRecommendation          string                      `json:"resolutionRecommendation"`
	ResolutionRecommendationPlainText string                      `json:"resolutionRecommendationPlainText"`
	ScopedProject                     ControlQueryProject         `json:"scopeProject"`
	Severity                          string                      `json:"severity"`
	SourceCloudConfigurationRule      ControlQueryCloudConfigRule `json:"sourceCloudConfigurationRule"`
	SupportsNRT                       bool                        `json:"supportsNRT"`
	Tags                              []string                    `json:"tags"`
	Type                              string                      `json:"type"`
}

// Cloud configuration rule object
type ControlQueryCloudConfigRule struct {
	Id string `json:"id"`
}

// Project object
type ControlQueryProject struct {
	Id string `json:"id"`
}

// User object
type ControlQueryUser struct {
	Id string `json:"id"`
}

// Relay-style node for the control
type ControlConnection struct {
	Nodes      []Control `json:"nodes"`
	PageInfo   PageInfo  `json:"pageInfo"`
	TotalCount int       `json:"totalCount"`
}

// ListControlsResponse is returned by ListControls on success
type ListControlsResponse struct {
	Controls ControlConnection `json:"controls"`
}

// Fields used to filter the control response
type ListControlsRequestConfiguration struct {
	// Optional - filter controls which are enabled/disabled.
	Enabled *bool

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - filter controls using any of securityFramework | securitySubCategory | securityCategory.
	FrameworkCategory string

	// Optional - filter controls which their related cloud configuration rule have auto remediation.
	HasAutoRemediation *bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter controls by project ID.
	Project string

	// Optional - filter controls by severity.
	//
	// Possible values are: CRITICAL, HIGH, INFORMATIONAL, LOW, MEDIUM.
	Severity string

	// Optional - filter controls by type.
	//
	// Possible values are: CLOUD_CONFIGURATION, SECURITY_GRAPH.
	Type string
}

// ListControls returns a paginated list of the security controls
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListControls(
	ctx context.Context,
	client *Client,
	options *ListControlsRequestConfiguration,
) (*ListControlsResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryControlList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.Enabled != nil {
		filter["enabled"] = options.Enabled
	}
	if options.FrameworkCategory != "" {
		filter["frameworkCategory"] = options.FrameworkCategory
	}
	if options.HasAutoRemediation != nil {
		filter["hasAutoRemediation"] = options.HasAutoRemediation
	}
	if options.Project != "" {
		filter["project"] = options.Project
	}
	if options.Severity != "" {
		filter["severity"] = options.Severity
	}
	if options.Type != "" {
		filter["type"] = options.Type
	}
	req.Var("filter", filter)

	// Check for options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}
	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// execute api call
	var responseData ListControlsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetControlResponse is returned by GetControl on success
type GetControlResponse struct {
	Control Control `json:"control"`
}

// GetControl returns a specific user that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetControl(
	ctx context.Context,
	client *Client,
	id string,
) (*GetControlResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryControlGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetControlResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
