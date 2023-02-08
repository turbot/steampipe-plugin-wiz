package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryUserRoleList = `
query ListUserRoles($first: Int, $after: String) {
  userRoles(first: $first, after: $after) {
    pageInfo{
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      description
      scopes
      isProjectScoped
    }
  }
}
`
)

// User role object
type UserRole struct {
	Description     string   `json:"description"`
	Id              string   `json:"id"`
	IsProjectScoped bool     `json:"isProjectScoped"`
	Name            string   `json:"name"`
	Scopes          []string `json:"scopes"`
}

// Relay-style node for the user-role
type UserRoleConnection struct {
	Nodes      []UserRole `json:"nodes"`
	PageInfo   PageInfo   `json:"pageInfo"`
	TotalCount int        `json:"totalCount"`
}

// ListUserRolesResponse is returned by ListUserRoles on success
type ListUserRolesResponse struct {
	UserRoles UserRoleConnection `json:"userRoles"`
}

type ListUserRolesRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int
}

// ListUserRoles returns a paginated list of the available user roles
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListUserRoles(
	ctx context.Context,
	client *Client,
	options *ListUserRolesRequestConfiguration,
) (*ListUserRolesResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryUserRoleList)

	// Check for options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}

	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// execute api call
	var responseData ListUserRolesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
