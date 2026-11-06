package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query. Secret values and snippets are intentionally not selected.
const (
	querySecretFindingList = `
query ListSecretFindings($first: Int, $after: String, $filter: SecretInstanceFilters) {
  secretInstances(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...SecretFindingDetails
    }
  }
}
` + fragmentSecretFinding

	querySecretFindingGet = `
query GetSecretFinding($id: ID!) {
  secretInstance(id: $id) {
    ...SecretFindingDetails
  }
}
` + fragmentSecretFinding

	querySecretFindingListWithTickets = `
query ListSecretFindingsWithTickets($first: Int, $after: String, $filter: SecretInstanceFilters) {
  secretInstances(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...SecretFindingDetails
      ...SecretFindingTickets
    }
  }
}
` + fragmentSecretFinding + fragmentSecretFindingTickets + fragmentServiceTicket

	querySecretFindingGetWithTickets = `
query GetSecretFindingWithTickets($id: ID!) {
  secretInstance(id: $id) {
    ...SecretFindingDetails
    ...SecretFindingTickets
  }
}
` + fragmentSecretFinding + fragmentSecretFindingTickets + fragmentServiceTicket
)

// Fields selected for every SecretInstance row
const fragmentSecretFinding = `
fragment SecretFindingDetails on SecretInstance {
  id
  name
  type
  confidence
  severity
  status
  resolutionReason
  validationStatus
  lastValidatedAt
  isEncrypted
  isManaged
  hasLeaked
  origin
  externalId
  path
  lineNumber
  description
  remediation
  firstSeenAt
  lastSeenAt
  lastModifiedAt
  lastUpdatedAt
  resolvedAt
  codeToCloudPipelineStage
  rule {
    id
    name
    type
    validityCheckSupported
    isAiPowered
  }
  projects {
    id
    name
    slug
  }
  vcsCodeOwners {
    id
    name
  }
  vcsDetails {
    initialCommitHash
    initialCommitURL
  }
  relatedIssueAnalytics {
    issueCount
    criticalSeverityCount
    highSeverityCount
    mediumSeverityCount
    lowSeverityCount
    informationalSeverityCount
  }
  resource {
    id
    name
    type
    externalId
    nativeType
    region
    status
    cloudPlatform
    cloudAccount {
      id
      externalId
      name
      cloudProvider
    }
    typedProperties {
      ... on SecretInstanceResourceRepositoryBranch {
        isDefaultBranch
        repository {
          id
          name
        }
      }
    }
  }
}
`

// Service tickets of the issues related to a finding; only selected when requested
const fragmentSecretFindingTickets = `
fragment SecretFindingTickets on SecretInstance {
  relatedIssues(first: 50) {
    nodes {
      serviceTickets {
        ...ServiceTicketDetails
      }
    }
  }
  postureIssues {
    serviceTickets {
      ...ServiceTicketDetails
    }
  }
}
`

// Secret detection rule reference
type SecretDetectionRuleRef struct {
	Id                     string `json:"id"`
	IsAiPowered            bool   `json:"isAiPowered"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	ValidityCheckSupported bool   `json:"validityCheckSupported"`
}

// Repository branch details of the resource holding a secret
type SecretFindingResourceTypedProperties struct {
	IsDefaultBranch *bool `json:"isDefaultBranch,omitempty"`
	Repository      *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"repository,omitempty"`
}

// Resource where a secret was found
type SecretFindingResource struct {
	CloudAccount *struct {
		CloudProvider string `json:"cloudProvider,omitempty"`
		ExternalId    string `json:"externalId,omitempty"`
		Id            string `json:"id"`
		Name          string `json:"name,omitempty"`
	} `json:"cloudAccount,omitempty"`
	CloudPlatform   string                                `json:"cloudPlatform,omitempty"`
	ExternalId      string                                `json:"externalId,omitempty"`
	Id              string                                `json:"id"`
	Name            string                                `json:"name"`
	NativeType      string                                `json:"nativeType,omitempty"`
	Region          string                                `json:"region,omitempty"`
	Status          string                                `json:"status,omitempty"`
	Type            string                                `json:"type"`
	TypedProperties *SecretFindingResourceTypedProperties `json:"typedProperties,omitempty"`
}

// Version control details of a secret finding
type SecretFindingVcsDetails struct {
	InitialCommitHash string `json:"initialCommitHash,omitempty"`
	InitialCommitURL  string `json:"initialCommitURL,omitempty"`
}

// Secret finding object
type SecretFinding struct {
	FindingServiceTickets

	CodeToCloudPipelineStage string                   `json:"codeToCloudPipelineStage"`
	Confidence               string                   `json:"confidence"`
	Description              string                   `json:"description"`
	ExternalId               string                   `json:"externalId"`
	FirstSeenAt              string                   `json:"firstSeenAt"`
	HasLeaked                bool                     `json:"hasLeaked"`
	Id                       string                   `json:"id"`
	IsEncrypted              bool                     `json:"isEncrypted"`
	IsManaged                bool                     `json:"isManaged"`
	LastModifiedAt           string                   `json:"lastModifiedAt"`
	LastSeenAt               string                   `json:"lastSeenAt"`
	LastUpdatedAt            string                   `json:"lastUpdatedAt"`
	LastValidatedAt          string                   `json:"lastValidatedAt"`
	LineNumber               *int                     `json:"lineNumber"`
	Name                     string                   `json:"name"`
	Origin                   string                   `json:"origin"`
	Path                     string                   `json:"path"`
	Projects                 []FindingProject         `json:"projects"`
	RelatedIssueAnalytics    *RelatedIssueAnalytics   `json:"relatedIssueAnalytics"`
	Remediation              string                   `json:"remediation"`
	ResolutionReason         string                   `json:"resolutionReason"`
	ResolvedAt               string                   `json:"resolvedAt"`
	Resource                 *SecretFindingResource   `json:"resource"`
	Rule                     *SecretDetectionRuleRef  `json:"rule"`
	Severity                 string                   `json:"severity"`
	Status                   string                   `json:"status"`
	Type                     string                   `json:"type"`
	ValidationStatus         string                   `json:"validationStatus"`
	VcsCodeOwners            []VcsCodeOwner           `json:"vcsCodeOwners"`
	VcsDetails               *SecretFindingVcsDetails `json:"vcsDetails"`
}

// Relay-style node for the secret finding
type SecretFindingConnection struct {
	Nodes      []SecretFinding `json:"nodes"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

// ListSecretFindingsResponse is returned by ListSecretFindings on success
type ListSecretFindingsResponse struct {
	SecretFindings SecretFindingConnection `json:"secretInstances"`
}

// Fields used to filter the secret finding response
type ListSecretFindingsRequestConfiguration struct {
	// Optional - also select the service tickets of the related issues (slower).
	IncludeServiceTickets bool

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - return only findings first seen within the range.
	FirstSeenAt *DateFilter

	// Optional - return only findings on (or off) the default branch of a repository.
	IsDefaultBranch *bool

	// Optional - return only findings last updated within the range.
	LastUpdatedAt *DateFilter

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter by the source that reported the finding, e.g. WIZ.
	Origin []string

	// Optional - filter by the file path of the secret.
	Path []string

	// Optional - filter by the Wiz ID of the code repository.
	RepositoryId []string

	// Optional - filter by the Wiz ID of the resource holding the secret.
	ResourceId []string

	// Optional - filter by the type of the resource holding the secret, e.g. REPOSITORY_BRANCH.
	ResourceType []string

	// Optional - filter by secret detection rule ID.
	RuleId []string

	// Optional - filter by severity.
	Severity []string

	// Optional - filter by finding status. Defaults to all statuses.
	Status []string

	// Optional - filter by secret type, e.g. CLOUD_KEY, PASSWORD.
	Type []string

	// Optional - filter by validation status. Possible values are: VALID, INVALID, UNKNOWN.
	ValidationStatus []string
}

// ListSecretFindings returns a paginated list of the secret findings
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSecretFindings(
	ctx context.Context,
	client *Client,
	options *ListSecretFindingsRequestConfiguration,
) (*ListSecretFindingsResponse, error) {
	// Make a request
	query := querySecretFindingList
	if options.IncludeServiceTickets {
		query = querySecretFindingListWithTickets
	}
	req := graphql.NewRequest(query)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.FirstSeenAt != nil {
		filter["firstSeenAt"] = options.FirstSeenAt.ToVar()
	}
	if options.IsDefaultBranch != nil {
		filter["isDefaultBranch"] = map[string]bool{"equals": *options.IsDefaultBranch}
	}
	if options.LastUpdatedAt != nil {
		filter["lastUpdatedAt"] = options.LastUpdatedAt.ToVar()
	}
	if len(options.Origin) > 0 {
		filter["origin"] = map[string][]string{"equals": options.Origin}
	}
	if len(options.Path) > 0 {
		filter["path"] = map[string][]string{"equals": options.Path}
	}
	if len(options.RepositoryId) > 0 {
		filter["codeRepository"] = map[string][]string{"equals": options.RepositoryId}
	}
	if len(options.ResourceId) > 0 {
		filter["resourceId"] = map[string][]string{"equals": options.ResourceId}
	}
	if len(options.ResourceType) > 0 {
		filter["resourceType"] = map[string][]string{"equals": options.ResourceType}
	}
	if len(options.RuleId) > 0 {
		filter["secretDetectionRule"] = map[string][]string{"equals": options.RuleId}
	}
	if len(options.Severity) > 0 {
		filter["severity"] = map[string][]string{"equals": options.Severity}
	}
	if len(options.Status) > 0 {
		filter["status"] = map[string][]string{"equals": options.Status}
	} else {
		filter["status"] = map[string][]string{"equals": DefaultFindingStatuses}
	}
	if len(options.Type) > 0 {
		filter["type"] = map[string][]string{"equals": options.Type}
	}
	if len(options.ValidationStatus) > 0 {
		filter["validationStatus"] = map[string][]string{"equals": options.ValidationStatus}
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
	var responseData ListSecretFindingsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetSecretFindingResponse is returned by GetSecretFinding on success
type GetSecretFindingResponse struct {
	SecretFinding *SecretFinding `json:"secretInstance"`
}

// GetSecretFinding returns a specific secret finding that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetSecretFinding(
	ctx context.Context,
	client *Client,
	id string,
	includeServiceTickets bool,
) (*GetSecretFindingResponse, error) {
	// Make a request
	query := querySecretFindingGet
	if includeServiceTickets {
		query = querySecretFindingGetWithTickets
	}
	req := graphql.NewRequest(query)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetSecretFindingResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
