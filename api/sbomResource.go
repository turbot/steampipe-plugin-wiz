package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// SBOM resources are read from sbomArtifactsGroupedByValues grouped by
// RESOURCE: one row per resource (repository branch, container image,
// virtual machine, ...) for which Wiz holds a software bill of materials,
// with the number of artifacts and distinct versions in it.
const (
	querySBOMResourceList = `
query ListSBOMResources($first: Int, $after: String, $filter: SBOMArtifactFilters) {
  sbomArtifactsGroupedByValues(
    first: $first
    after: $after
    filterBy: $filter
    groupBy: [RESOURCE]
  ) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      resource {
        ... on SBOMArtifactAssetBase {
          id
          name
          type
          externalId
          nativeType
          cloudPlatform
          projects {
            id
            name
          }
        }
      }
      analytics {
        artifacts {
          count
        }
        versions {
          count
        }
      }
    }
  }
}
`
)

// SBOMResourceProject is a project the resource belongs to
type SBOMResourceProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// SBOMResourceAsset is the resource an SBOM describes
type SBOMResourceAsset struct {
	CloudPlatform string                `json:"cloudPlatform"`
	ExternalId    string                `json:"externalId"`
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	NativeType    string                `json:"nativeType"`
	Projects      []SBOMResourceProject `json:"projects"`
	Type          string                `json:"type"`
}

// SBOMResourceCount holds a count from the grouped analytics
type SBOMResourceCount struct {
	Count int `json:"count"`
}

// SBOMResourceAnalytics holds the artifact and version counts of a resource
type SBOMResourceAnalytics struct {
	Artifacts SBOMResourceCount `json:"artifacts"`
	Versions  SBOMResourceCount `json:"versions"`
}

// SBOMResource is one resource with an SBOM
type SBOMResource struct {
	Analytics SBOMResourceAnalytics `json:"analytics"`
	Resource  *SBOMResourceAsset    `json:"resource"`
}

// SBOMResourceConnection is the paginated result
type SBOMResourceConnection struct {
	Nodes    []SBOMResource `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// ListSBOMResourcesResponse is returned by ListSBOMResources on success
type ListSBOMResourcesResponse struct {
	SBOMArtifactsGroupedByValues SBOMResourceConnection `json:"sbomArtifactsGroupedByValues"`
}

// ListSBOMResourcesRequestConfiguration holds the request options
type ListSBOMResourcesRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call.
	Limit int

	// Only return the resource with this Wiz identifier.
	ResourceId string

	// Only return resources of this type, such as REPOSITORY_BRANCH or
	// CONTAINER_IMAGE.
	ResourceType string
}

// ListSBOMResources returns a paginated list of the resources with an SBOM
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSBOMResources(
	ctx context.Context,
	client *Client,
	options *ListSBOMResourcesRequestConfiguration,
) (*ListSBOMResourcesResponse, error) {
	req := graphql.NewRequest(querySBOMResourceList)

	resourceFilter := map[string]interface{}{}
	if options.ResourceId != "" {
		resourceFilter["id"] = []string{options.ResourceId}
	}
	if options.ResourceType != "" {
		resourceFilter["type"] = []string{options.ResourceType}
	}
	if len(resourceFilter) > 0 {
		req.Var("filter", map[string]interface{}{"resource": resourceFilter})
	}

	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}
	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	var responseData ListSBOMResourcesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
