package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryProjectList = `
query ListProjects($first: Int, $after: String) {
  projects(first: $first, after: $after) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      description
      businessUnit
      archived
      slug
      securityScore
      riskProfile {
        businessImpact
      }
      profileCompletion
      repositoryCount
      cloudAccountCount
      cloudOrganizationCount
      kubernetesClusterCount
      workloadCount
      teamMemberCount
      entityCount
      technologyCount
      projectOwners {
        id
      }
      identifiers
      resourceTagLinks {
        environment
        resourceTags {
          key
          value
        }
      }
    }
  }
}
`
	queryProjectGet = `
query GetProject($id: ID!) {
  project(id: $id) {
    id
    name
    description
    businessUnit
    archived
    slug
    securityScore
    riskProfile {
      businessImpact
    }
    profileCompletion
    repositoryCount
    cloudAccountCount
    cloudOrganizationCount
    kubernetesClusterCount
    workloadCount
    teamMemberCount
    entityCount
    technologyCount
    projectOwners {
      id
    }
    identifiers
    resourceTagLinks {
      environment
      resourceTags {
        key
        value
      }
    }
  }
}
`
)

// Owner information
type ProjectOwner struct {
	Id string `json:"id"`
}

// Project resource tags link object
type ProjectResourceTagLinks struct {
	Environment  string        `json:"environment"`
	ResourceTags []ResourceTag `json:"resourceTags"`
}

// Risk profile object
type ProjectRiskProfile struct {
	BusinessImpact string `json:"businessImpact"`
}

// Resource tag object
type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Project object
type Project struct {
	Archived               bool                      `json:"archived"`
	BusinessUnit           string                    `json:"businessUnit"`
	CloudAccountCount      int                       `json:"cloudAccountCount"`
	CloudOrganizationCount int                       `json:"cloudOrganizationCount"`
	Description            string                    `json:"description"`
	EntityCount            int                       `json:"entityCount"`
	Id                     string                    `json:"id"`
	Identifiers            []string                  `json:"identifiers"`
	KubernetesClusterCount int                       `json:"kubernetesClusterCount"`
	Name                   string                    `json:"name"`
	ProfileCompletion      int                       `json:"profileCompletion"`
	ProjectOwners          []ProjectOwner            `json:"projectOwners"`
	RepositoryCount        int                       `json:"repositoryCount"`
	ResourceTagLinks       []ProjectResourceTagLinks `json:"resourceTagLinks"`
	RiskProfile            ProjectRiskProfile        `json:"riskProfile"`
	SecurityScore          int                       `json:"securityScore"`
	Slug                   string                    `json:"slug"`
	TeamMemberCount        int                       `json:"teamMemberCount"`
	TechnologyCount        int                       `json:"technologyCount"`
	WorkloadCount          int                       `json:"workloadCount"`
}

// Relay-style node for the project
type ProjectConnection struct {
	Nodes      []Project `json:"nodes"`
	PageInfo   PageInfo  `json:"pageInfo"`
	TotalCount int       `json:"totalCount"`
}

// ListProjectsResponse is returned by ListProjects on success
type ListProjectsResponse struct {
	Projects ProjectConnection `json:"projects"`
}

// Fields used to filter the project response
type ListProjectsRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int
}

// ListProjects returns a paginated list of the project
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListProjects(
	ctx context.Context,
	client *Client,
	options *ListProjectsRequestConfiguration,
) (*ListProjectsResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryProjectList)

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
	var data ListProjectsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
}

// GetProjectResponse is returned by GetProject on success
type GetProjectResponse struct {
	Project Project `json:"project"`
}

// GetProject returns a specific project that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetProject(
	ctx context.Context,
	client *Client,
	id string,
) (*GetProjectResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryProjectGet)

	// Set the required variables
	req.Var("id", id)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+*client.Token)

	var err error
	var data GetProjectResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		return nil, err
	}

	return &data, err
}
