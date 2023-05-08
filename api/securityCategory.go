package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	querySecurityCategoryList = `
query ListCategories($first: Int, $after: String, $filter: SecurityCategoryFilters) {
	securityCategories(first: $first, after: $after, filterBy: $filter) {
		pageInfo {
			hasNextPage
			endCursor
		}
		totalCount
		nodes {
			name
			id
			description
			framework {
				id
			}
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

	querySecurityCategoryGet = `
query GetCategory($id: ID!) {
	securityCategory(id: $id) {
		name
		id
		description
		framework {
			id
		}
		subCategories {
			id
      title
      description
      resolutionRecommendation
		}
	}
}
`
)

// Security category object
type SecurityCategory struct {
	Description   string                         `json:"description"`
	Framework     SecurityCategoryQueryFramework `json:"framework"`
	Id            string                         `json:"id"`
	Name          string                         `json:"name"`
	SubCategories []SecuritySubCategory          `json:"subCategories"`
}

// Security sub-category object
type SecuritySubCategory struct {
	Description              string `json:"description"`
	Id                       string `json:"id"`
	ResolutionRecommendation string `json:"resolutionRecommendation"`
	Title                    string `json:"title"`
}

type SecurityCategoryQueryFramework struct {
	Id string `json:"id"`
}

// Relay-style node for the security category
type SecurityCategoryConnection struct {
	Nodes      []SecurityCategory `json:"nodes"`
	PageInfo   PageInfo           `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

type SecurityCategoryQueryCategory struct {
	Id string `json:"id"`
}

// ListSecurityCategoriesResponse is returned by ListSecurityCategories on success
type ListSecurityCategoriesResponse struct {
	SecurityCategories SecurityCategoryConnection `json:"securityCategories"`
}

// Fields used to filter the security category response
type ListSecurityCategoriesRequestConfiguration struct {
	// Optional - filter security categories of specific framework ID.
	FrameworkId string

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int
}

// ListSecurityCategories returns a paginated list of security category
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSecurityCategories(
	ctx context.Context,
	client *Client,
	options *ListSecurityCategoriesRequestConfiguration,
) (*ListSecurityCategoriesResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySecurityCategoryList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.FrameworkId != "" {
		filter["securityFramework"] = options.FrameworkId
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
	var responseData ListSecurityCategoriesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetSecurityCategoryResponse is returned by GetSecurityCategory on success
type GetSecurityCategoryResponse struct {
	SecurityCategory SecurityCategory `json:"securityCategory"`
}

// GetSecurityCategory returns a specific security category that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetSecurityCategory(
	ctx context.Context,
	client *Client,
	id string,
) (*GetSecurityCategoryResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySecurityCategoryGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetSecurityCategoryResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
