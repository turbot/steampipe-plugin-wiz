package api

import (
	"context"

	"github.com/machinebox/graphql"
)

type Tenant struct {
	Id string `json:"id"`
}

// User object
type User struct {
	CreatedAt          string   `json:"createdAt"`
	Email              string   `json:"email"`
	Id                 string   `json:"id"`
	IpAddress          string   `json:"ipAddress"`
	IsAnalyticsEnabled bool     `json:"isAnalyticsEnabled"`
	IsSuspended        bool     `json:"isSuspended"`
	LastLoginAt        string   `json:"lastLoginAt"`
	Name               string   `json:"name"`
	Role               UserRole `json:"role"`
	Tenant             Tenant   `json:"tenant"`
}

// Relay-style node for the user
type Users struct {
	Nodes      []User   `json:"nodes"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount"`
}

// ListUsersResponse is returned by ListUsers on success
type ListUsersResponse struct {
	Users Users `json:"users"`
}

// GetUserResponse is returned by GetUser on success
type GetUserResponse struct {
	User User `json:"user"`
}

type ListUsersRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
}

// Define the query
const (
	queryUserList = `
query ListUsers($first: Int, $after: String) {
  users(first: $first, after: $after) {
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
