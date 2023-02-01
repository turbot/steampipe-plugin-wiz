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
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Enabled     bool               `json:"enabled"`
	BuiltIn     bool               `json:"builtin"`
	Description string             `json:"description"`
	PolicyTypes []string           `json:"policyTypes"`
	Categories  []SecurityCategory `json:"categories"`
}

// Relay-style node for the security framework
type SecurityFrameworkConnection struct {
	Nodes      []SecurityFramework `json:"nodes"`
	PageInfo   PageInfo            `json:"pageInfo"`
	TotalCount int                 `json:"totalCount"`
}

// ListSecurityFrameworksResponse is returned by ListSecurityFrameworks on success
type ListSecurityFrameworksResponse struct {
	SecurityFrameworks SecurityFrameworkConnection `json:"securityFrameworks"`
}

// Fields used to filter the security framework response
type ListSecurityFrameworksRequestConfiguration struct {
	// Optional - filter security frameworks which are enabled/disabled.
	Enabled *bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
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

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data ListSecurityFrameworksResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
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

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data GetSecurityFrameworkResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
}
