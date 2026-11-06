package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	querySASTFindingList = `
query ListSASTFindings($first: Int, $after: String, $filter: SASTFindingFilters) {
  sastFindings(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...SASTFindingDetails
    }
  }
}
` + fragmentSASTFinding

	querySASTFindingGet = `
query GetSASTFinding($id: ID!) {
  sastFinding(id: $id) {
    ...SASTFindingDetails
  }
}
` + fragmentSASTFinding

	querySASTFindingListWithTickets = `
query ListSASTFindingsWithTickets($first: Int, $after: String, $filter: SASTFindingFilters) {
  sastFindings(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...SASTFindingDetails
      ...SASTFindingTickets
    }
  }
}
` + fragmentSASTFinding + fragmentSASTFindingTickets + fragmentServiceTicket

	querySASTFindingGetWithTickets = `
query GetSASTFindingWithTickets($id: ID!) {
  sastFinding(id: $id) {
    ...SASTFindingDetails
    ...SASTFindingTickets
  }
}
` + fragmentSASTFinding + fragmentSASTFindingTickets + fragmentServiceTicket
)

// Fields selected for every SASTFinding row
const fragmentSASTFinding = `
fragment SASTFindingDetails on SASTFinding {
  id
  name
  description
  severity
  originalSeverity
  impact
  likelihood
  status
  resolutionReason
  rejectionExpiredAt
  origin
  isAiPowered
  createdAt
  updatedAt
  firstDetectedAtSource
  filePath
  startLine
  startColumn
  endLine
  endColumn
  snippet
  remediationInstructions
  sourceUrl
  codeLibraryLanguage
  repository {
    id
    name
    platform
  }
  repositoryBranch {
    id
    name
  }
  rule {
    id
    name
    externalId
  }
  weaknesses {
    id
    name
  }
  projects {
    id
    name
    slug
  }
  vcsDetails {
    commitHash
    commitURL
    codeOwners {
      id
      name
    }
  }
  aiAnalysis {
    verdict
    likelihood
    confidence
    explanation
  }
  relatedIssueAnalytics {
    issueCount
    criticalSeverityCount
    highSeverityCount
    mediumSeverityCount
    lowSeverityCount
    informationalSeverityCount
  }
}
`

// Service tickets of the issues related to a finding; only selected when requested
const fragmentSASTFindingTickets = `
fragment SASTFindingTickets on SASTFinding {
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

// Version control resource (repository or branch) referenced by a finding
type VersionControlResourceRef struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform,omitempty"`
}

// SAST rule reference
type SASTRuleRef struct {
	ExternalId string `json:"externalId"`
	Id         string `json:"id"`
	Name       string `json:"name"`
}

// CWE weakness
type Weakness struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Version control details of a SAST finding
type SASTFindingVcsDetails struct {
	CodeOwners []GraphEntity `json:"codeOwners"`
	CommitHash string        `json:"commitHash,omitempty"`
	CommitURL  string        `json:"commitURL,omitempty"`
}

// AI analysis of a SAST finding
type SASTFindingAIAnalysis struct {
	Confidence  string `json:"confidence"`
	Explanation string `json:"explanation"`
	Likelihood  string `json:"likelihood"`
	Verdict     string `json:"verdict"`
}

// SAST finding object
type SASTFinding struct {
	FindingServiceTickets

	AiAnalysis              *SASTFindingAIAnalysis     `json:"aiAnalysis"`
	CodeLibraryLanguage     []string                   `json:"codeLibraryLanguage"`
	CreatedAt               string                     `json:"createdAt"`
	Description             string                     `json:"description"`
	EndColumn               *int                       `json:"endColumn"`
	EndLine                 *int                       `json:"endLine"`
	FilePath                string                     `json:"filePath"`
	FirstDetectedAtSource   string                     `json:"firstDetectedAtSource"`
	Id                      string                     `json:"id"`
	Impact                  string                     `json:"impact"`
	IsAiPowered             bool                       `json:"isAiPowered"`
	Likelihood              string                     `json:"likelihood"`
	Name                    string                     `json:"name"`
	Origin                  string                     `json:"origin"`
	OriginalSeverity        string                     `json:"originalSeverity"`
	Projects                []FindingProject           `json:"projects"`
	RejectionExpiredAt      string                     `json:"rejectionExpiredAt"`
	RelatedIssueAnalytics   *RelatedIssueAnalytics     `json:"relatedIssueAnalytics"`
	RemediationInstructions string                     `json:"remediationInstructions"`
	Repository              *VersionControlResourceRef `json:"repository"`
	RepositoryBranch        *VersionControlResourceRef `json:"repositoryBranch"`
	ResolutionReason        string                     `json:"resolutionReason"`
	Rule                    *SASTRuleRef               `json:"rule"`
	Severity                string                     `json:"severity"`
	Snippet                 string                     `json:"snippet"`
	SourceUrl               string                     `json:"sourceUrl"`
	StartColumn             *int                       `json:"startColumn"`
	StartLine               *int                       `json:"startLine"`
	Status                  string                     `json:"status"`
	UpdatedAt               string                     `json:"updatedAt"`
	VcsDetails              *SASTFindingVcsDetails     `json:"vcsDetails"`
	Weaknesses              []Weakness                 `json:"weaknesses"`
}

// Relay-style node for the SAST finding
type SASTFindingConnection struct {
	Nodes      []SASTFinding `json:"nodes"`
	PageInfo   PageInfo      `json:"pageInfo"`
	TotalCount int           `json:"totalCount"`
}

// ListSASTFindingsResponse is returned by ListSASTFindings on success
type ListSASTFindingsResponse struct {
	SASTFindings SASTFindingConnection `json:"sastFindings"`
}

// Fields used to filter the SAST finding response
type ListSASTFindingsRequestConfiguration struct {
	// Optional - also select the service tickets of the related issues (slower).
	IncludeServiceTickets bool

	// Optional - filter by AI analysis verdict, e.g. EXPLOITABLE, TEST_CODE.
	AiVerdict []string

	// Optional - return only findings created within the range.
	CreatedAt *DateFilter

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - return only findings in the given file path.
	FilePath []string

	// Optional - return only findings from AI-powered (true) or rule-based (false) scanning.
	IsAiPowered *bool

	// Optional - return only findings on (or off) the default branch of a repository.
	IsDefaultBranch *bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter by the source that reported the finding, e.g. WIZ.
	Origin []string

	// Optional - filter by repository branch ID.
	RepositoryBranchId []string

	// Optional - filter by repository ID.
	RepositoryId []string

	// Optional - filter by the Wiz ID of the SAST rule.
	RuleId []string

	// Optional - filter by severity. Possible values are: CRITICAL, HIGH, MEDIUM, LOW, INFORMATIONAL.
	Severity []string

	// Optional - filter by finding status. Defaults to all statuses.
	Status []string

	// Optional - return only findings updated within the range.
	UpdatedAt *DateFilter
}

// ListSASTFindings returns a paginated list of the SAST findings
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSASTFindings(
	ctx context.Context,
	client *Client,
	options *ListSASTFindingsRequestConfiguration,
) (*ListSASTFindingsResponse, error) {
	// Make a request
	query := querySASTFindingList
	if options.IncludeServiceTickets {
		query = querySASTFindingListWithTickets
	}
	req := graphql.NewRequest(query)

	// Check for optional filters
	filter := map[string]interface{}{}
	if len(options.AiVerdict) > 0 {
		filter["aiVerdict"] = map[string][]string{"equals": options.AiVerdict}
	}
	if options.CreatedAt != nil {
		filter["createdAt"] = options.CreatedAt.ToVar()
	}
	if len(options.FilePath) > 0 {
		filter["filePath"] = map[string][]string{"equals": options.FilePath}
	}
	if options.IsAiPowered != nil {
		filter["isAiPowered"] = map[string]bool{"equals": *options.IsAiPowered}
	}
	if len(options.Origin) > 0 {
		filter["origin"] = map[string][]string{"equals": options.Origin}
	}
	if len(options.RuleId) > 0 {
		filter["rule"] = map[string][]string{"equals": options.RuleId}
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
	if options.IsDefaultBranch != nil {
		resource["isDefaultBranch"] = map[string]bool{"equals": *options.IsDefaultBranch}
	}
	if len(options.RepositoryBranchId) > 0 {
		resource["repositoryBranchId"] = map[string][]string{"equals": options.RepositoryBranchId}
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
	var responseData ListSASTFindingsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetSASTFindingResponse is returned by GetSASTFinding on success
type GetSASTFindingResponse struct {
	SASTFinding *SASTFinding `json:"sastFinding"`
}

// GetSASTFinding returns a specific SAST finding that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetSASTFinding(
	ctx context.Context,
	client *Client,
	id string,
	includeServiceTickets bool,
) (*GetSASTFindingResponse, error) {
	// Make a request
	query := querySASTFindingGet
	if includeServiceTickets {
		query = querySASTFindingGetWithTickets
	}
	req := graphql.NewRequest(query)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetSASTFindingResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
