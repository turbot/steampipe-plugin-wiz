package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryUserList = `
query ListUsers($first: Int, $after: String, $filter: UserFilters) {
  users(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      name
      id
      email
      createdAt
      lastLoginAt
      isSuspended
      isAnalyticsEnabled
			identityProviderType
      role {
        id
        name
        description
        scopes
        isProjectScoped
      }
      tenant {
        id
      }
			effectiveAssignedProjects {
				id
			}
      ipAddress
    }
  }
}
`

	queryGetUser = `
query GetUser($id: ID!) {
  user(id: $id) {
    name
    id
    email
    createdAt
    lastLoginAt
    isSuspended
    isAnalyticsEnabled
    role {
      id
      name
      description
      scopes
      isProjectScoped
    }
    tenant {
      id
    }
    ipAddress
  }
}
`
)

// Assigned project information
type UserQueryProject struct {
	Id string `json:"id"`
}

// Tenant object
type Tenant struct {
	Id string `json:"id"`
}

// User object
type User struct {
	CreatedAt                 string             `json:"createdAt"`
	EffectiveAssignedProjects []UserQueryProject `json:"effectiveAssignedProjects"`
	Email                     string             `json:"email"`
	Id                        string             `json:"id"`
	IdentityProviderType      string             `json:"identityProviderType"`
	IpAddress                 string             `json:"ipAddress"`
	IsAnalyticsEnabled        bool               `json:"isAnalyticsEnabled"`
	IsSuspended               bool               `json:"isSuspended"`
	LastLoginAt               string             `json:"lastLoginAt"`
	Name                      string             `json:"name"`
	Role                      UserRole           `json:"role"`
	Tenant                    Tenant             `json:"tenant"`
}

// Relay-style node for the user
type UserConnection struct {
	Nodes      []User   `json:"nodes"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount"`
}

// ListUsersResponse is returned by ListUsers on success
type ListUsersResponse struct {
	Users UserConnection `json:"users"`
}

// Fields used to filter the user response
type ListUsersRequestConfiguration struct {
	// Optional - filter by provider type.
	//
	// Possible values are: WIZ, SAML.
	AuthProviderType string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
}

// ListUsers returns a paginated list of the portal users
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListUsers(
	ctx context.Context,
	client *Client,
	options *ListUsersRequestConfiguration,
) (*ListUsersResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryUserList)

	// Check for optional filters
	filter := map[string]string{}
	if options.AuthProviderType != "" {
		filter["authProviderType"] = options.AuthProviderType
	}
	req.Var("filter", filter)

	// Check for options and set it
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
	var data ListUsersResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		// err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}

// GetUserResponse is returned by GetUser on success
type GetUserResponse struct {
	User User `json:"user"`
}

// GetUser returns a specific user that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetUser(
	ctx context.Context,
	client *Client,
	id string,
) (*GetUserResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryGetUser)

	// Set the required variables
	req.Var("id", id)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data GetUserResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		// err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
