package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	querySecurityFrameworkList = `
query ListSecurityFrameworks($first: Int, $after: String, $filter: SecurityFrameworkFilters) {
  securityFrameworks(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      enabled
      builtin
      description
      policyTypes
      categories {
        id
        name
        description
        subCategories {
          id
          title
          description
          resolutionRecommendation
        }
      }
    }
  }
}
`
	querySecurityFrameworkGet = `
query GetSecurityFramework($id: ID!) {
  securityFramework(id: $id) {
    id
    name
    enabled
    builtin
    description
    policyTypes
    categories {
      id
      name
      description
      subCategories {
        id
        title
        description
        resolutionRecommendation
      }
    }
  }
}
`
)

// Security framework object
type SecurityFramework struct {
	BuiltIn     bool                             `json:"builtin"`
	Categories  []SecurityFrameworkQueryCategory `json:"categories"`
	Description string                           `json:"description"`
	Enabled     bool                             `json:"enabled"`
	Id          string                           `json:"id"`
	Name        string                           `json:"name"`
	PolicyTypes []string                         `json:"policyTypes"`
}

// Relay-style node for the security framework
type SecurityFrameworkConnection struct {
	Nodes      []SecurityFramework `json:"nodes"`
	PageInfo   PageInfo            `json:"pageInfo"`
	TotalCount int                 `json:"totalCount"`
}

type SecurityFrameworkQueryCategory struct {
	Id string `json:"id"`
}

// ListSecurityFrameworksResponse is returned by ListSecurityFrameworks on success
type ListSecurityFrameworksResponse struct {
	SecurityFrameworks SecurityFrameworkConnection `json:"securityFrameworks"`
}

// Fields used to filter the security framework response
type ListSecurityFrameworksRequestConfiguration struct {
	// Optional - filter security frameworks which are enabled/disabled.
	Enabled *bool

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int
}

// ListSecurityFrameworks returns a paginated list of compliance frameworks
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSecurityFrameworks(
	ctx context.Context,
	client *Client,
	options *ListSecurityFrameworksRequestConfiguration,
) (*ListSecurityFrameworksResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySecurityFrameworkList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.Enabled != nil {
		filter["enabled"] = *options.Enabled
	}
	req.Var("filter", filter)

	// Check for paging options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}

	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// execute api call
	var responseData ListSecurityFrameworksResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetSecurityFrameworkResponse is returned by GetSecurityFramework on success
type GetSecurityFrameworkResponse struct {
	SecurityFramework SecurityFramework `json:"securityFramework"`
}

// GetSecurityFramework returns a specific security framework that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetSecurityFramework(
	ctx context.Context,
	client *Client,
	id string,
) (*GetSecurityFrameworkResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySecurityFrameworkGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetSecurityFrameworkResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
