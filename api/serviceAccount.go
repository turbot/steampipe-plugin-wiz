package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryServiceAccountList = `
query ListServiceAccounts($first: Int, $after: String, $filter: ServiceAccountFilters) {
  serviceAccounts(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      type
      createdAt
      lastRotatedAt
      clientId
      scopes
      assignedProjects {
        id
      }
      authenticationSource
    }
  }
}
`

	queryServiceAccountGet = `
query GetServiceAccount($id: ID!) {
  serviceAccount(id: $id) {
    id
    name
    createdAt
    lastRotatedAt
    clientId
    scopes
    assignedProjects {
      id
    }
    authenticationSource
    type
  }
}
`
)

// Assigned project information
type ServiceAccountQueryProject struct {
	Id string `json:"id"`
}

// Service account object
type ServiceAccount struct {
	AssignedProjects     []ServiceAccountQueryProject `json:"assignedProjects"`
	AuthenticationSource string                       `json:"authenticationSource"`
	ClientId             string                       `json:"clientId"`
	CreatedAt            string                       `json:"createdAt"`
	Id                   string                       `json:"id"`
	LastRotatedAt        string                       `json:"lastRotatedAt"`
	Name                 string                       `json:"name"`
	Scopes               []string                     `json:"scopes"`
	Type                 string                       `json:"type"`
}

// Relay-style node for the service account
type ServiceAccountConnection struct {
	Nodes      []ServiceAccount `json:"nodes"`
	PageInfo   PageInfo         `json:"pageInfo"`
	TotalCount int              `json:"totalCount"`
}

// ListServiceAccountsResponse is returned by ListServiceAccounts on success
type ListServiceAccountsResponse struct {
	ServiceAccounts ServiceAccountConnection `json:"serviceAccounts"`
}

// Fields used to filter the service account response
type ListServiceAccountsRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 60.
	Limit int

	// Optional - the name of the service account.
	Name string

	// Optional - the service account authentication source.
	//
	// Possible values: LEGACY, MODERN.
	Source string

	// Optional - the service account type.
	//
	// Possible values are: THIRD_PARTY, SENSOR, KUBERNETES_ADMISSION_CONTROLLER, BROKER.
	Type string
}

// ListServiceAccounts returns a paginated list of the service accounts
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListServiceAccounts(
	ctx context.Context,
	client *Client,
	options *ListServiceAccountsRequestConfiguration,
) (*ListServiceAccountsResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryServiceAccountList)

	// Check for optional filters
	filter := map[string]string{}
	if options.Name != "" {
		filter["name"] = options.Name
	}
	if options.Source != "" {
		filter["source"] = options.Source
	}
	if options.Type != "" {
		filter["type"] = options.Type
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
	var responseData ListServiceAccountsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetServiceAccountResponse is returned by GetServiceAccount on success
type GetServiceAccountResponse struct {
	ServiceAccount ServiceAccount `json:"serviceAccount"`
}

// GetServiceAccount returns a specific service account that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetServiceAccount(
	ctx context.Context,
	client *Client,
	id string,
) (*GetServiceAccountResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryServiceAccountGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetServiceAccountResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
