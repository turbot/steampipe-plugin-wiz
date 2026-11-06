package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryVersionControlRepositoryList = `
query ListVersionControlRepositories($first: Int, $after: String, $filter: VersionControlRepositoryFilters) {
  versionControlRepositories(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      url
      platform
      public
      archived
      visibility
      defaultBranch
      organization {
        id
        externalId
        name
      }
      projects {
        id
        name
        slug
      }
    }
  }
}
`

	queryVersionControlRepositoryGet = `
query GetVersionControlRepository($id: ID!) {
  repository(id: $id) {
    id
    name
    url
    platform
    public
    archived
    visibility
    defaultBranch
    organization {
      id
      externalId
      name
    }
    projects {
      id
      name
      slug
    }
  }
}
`
)

// Organization owning a repository
type VersionControlOrganization struct {
	ExternalId string `json:"externalId,omitempty"`
	Id         string `json:"id"`
	Name       string `json:"name"`
}

// Version control repository object
type VersionControlRepository struct {
	Archived      *bool                       `json:"archived"`
	DefaultBranch string                      `json:"defaultBranch"`
	Id            string                      `json:"id"`
	Name          string                      `json:"name"`
	Organization  *VersionControlOrganization `json:"organization"`
	Platform      string                      `json:"platform"`
	Projects      []FindingProject            `json:"projects"`
	Public        *bool                       `json:"public"`
	Url           string                      `json:"url"`
	Visibility    string                      `json:"visibility"`
}

// Relay-style node for the version control repository
type VersionControlRepositoryConnection struct {
	Nodes      []VersionControlRepository `json:"nodes"`
	PageInfo   PageInfo                   `json:"pageInfo"`
	TotalCount int                        `json:"totalCount"`
}

// ListVersionControlRepositoriesResponse is returned by ListVersionControlRepositories on success
type ListVersionControlRepositoriesResponse struct {
	Repositories VersionControlRepositoryConnection `json:"versionControlRepositories"`
}

// Fields used to filter the version control repository response
type ListVersionControlRepositoriesRequestConfiguration struct {
	// Optional - return only archived (true) or active (false) repositories.
	Archived *bool

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter by the organization ID.
	OrganizationId []string

	// Optional - filter by version control platform, e.g. GITHUB, GITLAB.
	Platform []string

	// Optional - return only public (true) or private (false) repositories.
	Public *bool

	// Optional - filter by repository visibility.
	Visibility []string
}

// ListVersionControlRepositories returns a paginated list of the version control repositories
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListVersionControlRepositories(
	ctx context.Context,
	client *Client,
	options *ListVersionControlRepositoriesRequestConfiguration,
) (*ListVersionControlRepositoriesResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryVersionControlRepositoryList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.Archived != nil {
		filter["archived"] = *options.Archived
	}
	if len(options.OrganizationId) > 0 {
		filter["organizationId"] = options.OrganizationId
	}
	if len(options.Platform) > 0 {
		filter["platform"] = options.Platform
	}
	if options.Public != nil {
		filter["public"] = *options.Public
	}
	if len(options.Visibility) > 0 {
		filter["visibility"] = options.Visibility
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
	var responseData ListVersionControlRepositoriesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetVersionControlRepositoryResponse is returned by GetVersionControlRepository on success
type GetVersionControlRepositoryResponse struct {
	Repository *VersionControlRepository `json:"repository"`
}

// GetVersionControlRepository returns a specific repository that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetVersionControlRepository(
	ctx context.Context,
	client *Client,
	id string,
) (*GetVersionControlRepositoryResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryVersionControlRepositoryGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetVersionControlRepositoryResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
