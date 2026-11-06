package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryIaCFindingList = `
query ListIaCFindings($first: Int, $after: String, $filter: IaCFindingFilters) {
  iacFindings(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...IaCFindingDetails
    }
  }
}
` + fragmentIaCFinding

	queryIaCFindingGet = `
query GetIaCFinding($id: ID!) {
  iacFinding(id: $id) {
    ...IaCFindingDetails
  }
}
` + fragmentIaCFinding
)

// Fields selected for every IaCFinding row
const fragmentIaCFinding = `
fragment IaCFindingDetails on IaCFinding {
  id
  name
  severity
  status
  firstSeenAt
  lastSeenAt
  updatedAt
  platform
  filePath
  startLine
  endLine
  fileURL
  vcsPlatform
  cloudPlatform
  entityType
  nativeType
  foundContent
  expectedContent
  matchContent
  remediationInstructions
  repository {
    id
    name
  }
  branch {
    id
    name
  }
  rule {
    id
    shortId
    name
    severity
    serviceType
    cloudProvider
    builtin
  }
  resourceGraphEntity {
    id
    type
    name
    providerUniqueId
  }
  vcsDetails {
    codeOwners {
      id
      name
    }
  }
}
`

// Cloud configuration rule reference
type IaCFindingRule struct {
	Builtin       bool   `json:"builtin"`
	CloudProvider string `json:"cloudProvider,omitempty"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	ServiceType   string `json:"serviceType,omitempty"`
	Severity      string `json:"severity,omitempty"`
	ShortId       string `json:"shortId,omitempty"`
}

// Graph entity the IaC finding maps to
type IaCFindingResource struct {
	Id               string `json:"id"`
	Name             string `json:"name,omitempty"`
	ProviderUniqueId string `json:"providerUniqueId,omitempty"`
	Type             string `json:"type,omitempty"`
}

// Version control details of an IaC finding
type IaCFindingVcsDetails struct {
	CodeOwners []VcsCodeOwner `json:"codeOwners"`
}

// IaC finding object
type IaCFinding struct {
	Branch                  *VersionControlResourceRef `json:"branch"`
	CloudPlatform           string                     `json:"cloudPlatform"`
	EndLine                 *int                       `json:"endLine"`
	EntityType              string                     `json:"entityType"`
	ExpectedContent         string                     `json:"expectedContent"`
	FilePath                string                     `json:"filePath"`
	FileURL                 string                     `json:"fileURL"`
	FirstSeenAt             string                     `json:"firstSeenAt"`
	FoundContent            string                     `json:"foundContent"`
	Id                      string                     `json:"id"`
	LastSeenAt              string                     `json:"lastSeenAt"`
	MatchContent            string                     `json:"matchContent"`
	Name                    string                     `json:"name"`
	NativeType              string                     `json:"nativeType"`
	Platform                string                     `json:"platform"`
	RemediationInstructions string                     `json:"remediationInstructions"`
	Repository              *VersionControlResourceRef `json:"repository"`
	ResourceGraphEntity     *IaCFindingResource        `json:"resourceGraphEntity"`
	Rule                    *IaCFindingRule            `json:"rule"`
	Severity                string                     `json:"severity"`
	StartLine               *int                       `json:"startLine"`
	Status                  string                     `json:"status"`
	UpdatedAt               string                     `json:"updatedAt"`
	VcsDetails              *IaCFindingVcsDetails      `json:"vcsDetails"`
	VcsPlatform             string                     `json:"vcsPlatform"`
}

// Relay-style node for the IaC finding
type IaCFindingConnection struct {
	Nodes      []IaCFinding `json:"nodes"`
	PageInfo   PageInfo     `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

// ListIaCFindingsResponse is returned by ListIaCFindings on success
type ListIaCFindingsResponse struct {
	IaCFindings IaCFindingConnection `json:"iacFindings"`
}

// Fields used to filter the IaC finding response
type ListIaCFindingsRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - filter by the file path of the finding.
	FilePath []string

	// Optional - return only findings on (or off) the default branch of a repository.
	IsDefaultBranch *bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter by IaC platform, e.g. TERRAFORM, CLOUD_FORMATION, KUBERNETES.
	Platform []string

	// Optional - filter by repository ID.
	RepositoryId []string

	// Optional - filter by the Wiz ID of the cloud configuration rule.
	RuleId []string

	// Optional - filter by severity.
	Severity []string

	// Optional - filter by finding status. Defaults to all statuses.
	Status []string

	// Optional - return only findings updated within the range.
	UpdatedAt *DateFilter
}

// ListIaCFindings returns a paginated list of the IaC findings
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListIaCFindings(
	ctx context.Context,
	client *Client,
	options *ListIaCFindingsRequestConfiguration,
) (*ListIaCFindingsResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryIaCFindingList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if len(options.RuleId) > 0 {
		filter["rule"] = map[string][]string{"id": options.RuleId}
	}
	if len(options.Severity) > 0 {
		filter["severity"] = map[string][]string{"equals": options.Severity}
	}
	if len(options.Status) > 0 {
		filter["status"] = map[string][]string{"equals": options.Status}
	} else {
		filter["status"] = map[string][]string{"equals": DefaultFindingStatuses}
	}
	if options.UpdatedAt != nil {
		filter["updatedAt"] = options.UpdatedAt.ToVar()
	}

	resource := map[string]interface{}{}
	if len(options.FilePath) > 0 {
		resource["filePathV2"] = map[string][]string{"equals": options.FilePath}
	}
	if options.IsDefaultBranch != nil {
		resource["isDefaultBranch"] = map[string]bool{"equals": *options.IsDefaultBranch}
	}
	if len(options.Platform) > 0 {
		resource["iacPlatform"] = map[string][]string{"equals": options.Platform}
	}
	if len(options.RepositoryId) > 0 {
		resource["repositoryId"] = map[string][]string{"equals": options.RepositoryId}
	}
	if len(resource) > 0 {
		filter["resource"] = resource
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
	var responseData ListIaCFindingsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetIaCFindingResponse is returned by GetIaCFinding on success
type GetIaCFindingResponse struct {
	IaCFinding *IaCFinding `json:"iacFinding"`
}

// GetIaCFinding returns a specific IaC finding that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetIaCFinding(
	ctx context.Context,
	client *Client,
	id string,
) (*GetIaCFindingResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryIaCFindingGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetIaCFindingResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
