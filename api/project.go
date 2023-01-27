package api

import (
	"context"

	"github.com/machinebox/graphql"
)

type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProjectOwner struct {
	Id string `json:"id"`
}

type ProjectResourceTagLinks struct {
	Environment  string        `json:"environment"`
	ResourceTags []ResourceTag `json:"resourceTags"`
}

type ProjectRiskProfile struct {
	BusinessImpact string `json:"businessImpact"`
}

type Project struct {
	Id                     string                    `json:"id"`
	Name                   string                    `json:"name"`
	Slug                   string                    `json:"slug"`
	Description            string                    `json:"description"`
	Archived               bool                      `json:"archived"`
	BusinessUnit           string                    `json:"businessUnit"`
	SecurityScore          int                       `json:"securityScore"`
	CloudAccountCount      int                       `json:"cloudAccountCount"`
	CloudOrganizationCount int                       `json:"cloudOrganizationCount"`
	EntityCount            int                       `json:"entityCount"`
	KubernetesClusterCount int                       `json:"kubernetesClusterCount"`
	ProfileCompletion      int                       `json:"profileCompletion"`
	RepositoryCount        int                       `json:"repositoryCount"`
	TeamMemberCount        int                       `json:"teamMemberCount"`
	TechnologyCount        int                       `json:"technologyCount"`
	WorkloadCount          int                       `json:"workloadCount"`
	Identifiers            []string                  `json:"identifiers"`
	ResourceTagLinks       []ProjectResourceTagLinks `json:"resourceTagLinks"`
	ProjectOwners          []ProjectOwner            `json:"projectOwners"`
	RiskProfile            ProjectRiskProfile        `json:"riskProfile"`
}

// Relay-style node for the project
type Projects struct {
	Nodes      []Project `json:"nodes"`
	PageInfo   PageInfo  `json:"pageInfo"`
	TotalCount int       `json:"totalCount"`
}

// ListProjectsResponse is returned by ListProjects on success
type ListProjectsResponse struct {
	Projects Projects `json:"projects"`
}

// GetProjectResponse is returned by GetProject on success
type GetProjectResponse struct {
	Project Project `json:"project"`
}

type ListProjectsRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
}

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
		// err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
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
		// err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
