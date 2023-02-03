package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryConfigurationFindingList = `
query ListConfigurationFindings($first: Int, $after: String, $filter: ConfigurationFindingFilters) {
  configurationFindings(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      result
      resource {
        name
        id
        type
        status
        region
        cloudPlatform
        nativeType
        projects {
          id
        }
        tags {
          key
          value
        }
      }
      rule {
        id
      }
      severity
      subscription {
        id
      }
      remediation
      analyzedAt
      status
      resolutionReason
      note {
        text
      }
    }
  }
}
`

	queryConfigurationFindingGet = `
query GetConfigurationFinding($id: ID!) {
  configurationFinding(id: $id) {
    id
    result
    resource {
      name
      id
      type
      status
      region
      cloudPlatform
      nativeType
      projects {
        id
      }
      tags {
        key
        value
      }
    }
    rule {
      id
    }
    severity
    subscription {
      id
    }
    remediation
    analyzedAt
    status
    resolutionReason
    note {
      text
    }
  }
}
`
)

// Configuration finding object
type ConfigurationFinding struct {
	AnalyzedAt       string
	Id               string
	Remediation      string
	ResolutionReason string
	Resource         ConfigurationFindingResource
	Result           string
	Rule             ConfigurationFindingQueryRule
	Severity         string
	Status           string
	Subscription     ConfigurationFindingQuerySubscription
}

// Configuration finding resource object
type ConfigurationFindingResource struct {
	CloudPlatform string
	Id            string
	Name          string
	NativeType    string
	Projects      []ConfigurationFindingResourceQueryProject
	Region        string
	Status        string
	Tags          []ConfigurationFindingResourceTag
	Type          string
}

// Project information
type ConfigurationFindingResourceQueryProject struct {
	Id string
}

// Resource tag object
type ConfigurationFindingResourceTag struct {
	Key   string
	Value string
}

// Cloud configuration rule information
type ConfigurationFindingQueryRule struct {
	Id string
}

// Subscription information
type ConfigurationFindingQuerySubscription struct {
	Id string
}

// Relay-style node for cloud configuration findings
type ConfigurationFindingConnection struct {
	Nodes      []ConfigurationFinding `json:"nodes"`
	PageInfo   PageInfo               `json:"pageInfo"`
	TotalCount int                    `json:"totalCount"`
}

// ListConfigurationFindingsResponse is returned by ListConfigurationFindings on success
type ListConfigurationFindingsResponse struct {
	ConfigurationFindings ConfigurationFindingConnection `json:"configurationFindings"`
}

// Fields used to filter the cloud configuration findings response
type ListConfigurationFindingsRequestConfiguration struct {
	// Optional - filter findings by result.
	//
	// Possible values are: ERROR, FAIL, NOT_ASSESSED, PASS.
	Result string

	// Optional - filter findings by severity.
	//
	// Possible values: CRITICAL, HIGH, LOW, MEDIUM, NONE.
	Severity string

	// Optional - filter findings by status.
	//
	// Possible values: IN_PROGRESS, OPEN, REJECTED, RESOLVED.
	Status string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
}

// ListConfigurationFindings returns a paginated list of cloud configuration findings
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListConfigurationFindings(
	ctx context.Context,
	client *Client,
	options *ListConfigurationFindingsRequestConfiguration,
) (*ListConfigurationFindingsResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryConfigurationFindingList)

	// Check for optional filters
	filter := map[string]string{}
	if options.Result != "" {
		filter["result"] = options.Result
	}
	if options.Severity != "" {
		filter["severity"] = options.Severity
	}
	if options.Status != "" {
		filter["status"] = options.Status
	}
	req.Var("filter", filter)

	// Check for paging options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}

	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data ListConfigurationFindingsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
}

// GetConfigurationFindingResponse is returned by GetConfigurationFinding on success
type GetConfigurationFindingResponse struct {
	ConfigurationFinding ConfigurationFinding `json:"configurationFinding"`
}

// GetConfigurationFinding returns a specific finding that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetConfigurationFinding(
	ctx context.Context,
	client *Client,
	id string,
) (*GetConfigurationFindingResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryConfigurationFindingGet)

	// Set the required variables
	req.Var("id", id)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data GetConfigurationFindingResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
}
