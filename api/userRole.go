package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryUserRoleList = `
query ListUserRoles($filter: UserRoleFilters) {
  userRolesV2(filterBy: $filter) {
    id
    name
    description
    scopes
    isProjectScoped
    isTenantGroupScoped
    builtin
    createdAt
    updatedAt
    userCount
  }
}
`
)

// User role object
type UserRole struct {
	Builtin             bool     `json:"builtin,omitempty"`
	CreatedAt           string   `json:"createdAt,omitempty"`
	Description         string   `json:"description"`
	Id                  string   `json:"id"`
	IsProjectScoped     bool     `json:"isProjectScoped"`
	IsTenantGroupScoped bool     `json:"isTenantGroupScoped,omitempty"`
	Name                string   `json:"name"`
	Scopes              []string `json:"scopes"`
	UpdatedAt           string   `json:"updatedAt,omitempty"`
	UserCount           *int     `json:"userCount,omitempty"`
}

// ListUserRolesResponse is returned by ListUserRoles on success
type ListUserRolesResponse struct {
	UserRoles []UserRole `json:"userRolesV2"`
}

type ListUserRolesRequestConfiguration struct {
	// Optional - return only built-in (true) or custom (false) roles.
	Builtin *bool

	// Optional - return only project scoped (true) or global (false) roles.
	IsProjectScoped *bool
}

// ListUserRoles returns all the available user roles. The API returns them in a single response.
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

	filter := map[string]interface{}{}
	if options.Builtin != nil {
		filter["builtin"] = *options.Builtin
	}
	if options.IsProjectScoped != nil {
		filter["isProjectScoped"] = *options.IsProjectScoped
	}
	req.Var("filter", filter)

	// execute api call
	var responseData ListUserRolesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
