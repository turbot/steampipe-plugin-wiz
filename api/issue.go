package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryIssueList = `
query ListIssues($first: Int, $after: String, $filter: IssueFilters) {
  issuesV2(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...IssueDetails
    }
  }
}
` + fragmentIssue

	queryIssueListWithExtras = `
query ListIssuesWithExtras($first: Int, $after: String, $filter: IssueFilters) {
  issuesV2(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      ...IssueDetails
      ...IssueExtras
    }
  }
}
` + fragmentIssue + fragmentIssueExtras

	queryIssueGet = `
query GetIssue($id: ID!) {
  issue(id: $id) {
    ...IssueDetails
  }
}
` + fragmentIssue

	queryIssueGetWithExtras = `
query GetIssueWithExtras($id: ID!) {
  issue(id: $id) {
    ...IssueDetails
    ...IssueExtras
  }
}
` + fragmentIssue + fragmentIssueExtras
)

// Fields selected for every Issue row
const fragmentIssue = `
fragment IssueDetails on Issue {
  id
  type
  description
  status
  severity
  createdAt
  updatedAt
  resolvedAt
  dueAt
  rejectionExpiredAt
  statusChangedAt
  firstEventAt
  lastEventAt
  resolutionReason
  resolutionNote
  openReason
  url
  hasCodeRemediation
  hasCodeToCloudMapping
  validatedAsExploitable
  environments
  entity {
    id
    name
  }
  entitySnapshot {
    id
    type
    nativeType
    name
    cloudPlatform
    externalId
    providerId
    region
    status
    subscriptionId
    subscriptionExternalId
    subscriptionName
    kubernetesClusterId
    kubernetesClusterName
    kubernetesNamespaceName
    cloudProviderURL
    tags
  }
  sourceRules {
    __typename
    id
    name
  }
  notes {
    id
    createdAt
    updatedAt
    text
    user {
      id
      name
    }
    serviceAccount {
      id
    }
  }
  serviceTickets {
    id
    name
    url
    externalId
    integration {
      id
      name
    }
    project {
      id
      name
    }
  }
  projects {
    id
    name
  }
}
`

// Application services and risk factors of an issue; only selected when requested
const fragmentIssueExtras = `
fragment IssueExtras on Issue {
  applicationServices {
    id
    displayName
  }
  relatedApplicationServices {
    id
    displayName
  }
  riskFactors {
    id
    severity
    policy {
      id
      name
    }
  }
}
`

// The graph object
type GraphEntity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Issue note object
type IssueNote struct {
	CreatedAt      string                   `json:"createdAt,omitempty"`
	Id             string                   `json:"id,omitempty"`
	ServiceAccount IssueQueryServiceAccount `json:"serviceAccount,omitempty"`
	Text           string                   `json:"text,omitempty"`
	UpdatedAt      string                   `json:"updatedAt,omitempty"`
	User           IssueQueryUser           `json:"user,omitempty"`
}

// Automation action information
type IssueQueryAutomationAction struct {
	Id string `json:"id"`
}

// Source rule (control, cloud configuration rule or cloud event rule) that generated the issue
type IssueSourceRule struct {
	Typename string `json:"__typename"`
	Id       string `json:"id"`
	Name     string `json:"name"`
}

// Integration information
type IssueQueryIntegration struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// Snapshot of the entity an issue is about, taken when the issue was last updated
type IssueEntitySnapshot struct {
	CloudPlatform           string      `json:"cloudPlatform,omitempty"`
	CloudProviderURL        string      `json:"cloudProviderURL,omitempty"`
	ExternalId              string      `json:"externalId,omitempty"`
	Id                      string      `json:"id"`
	KubernetesClusterId     string      `json:"kubernetesClusterId,omitempty"`
	KubernetesClusterName   string      `json:"kubernetesClusterName,omitempty"`
	KubernetesNamespaceName string      `json:"kubernetesNamespaceName,omitempty"`
	Name                    string      `json:"name"`
	NativeType              string      `json:"nativeType,omitempty"`
	ProviderId              string      `json:"providerId,omitempty"`
	Region                  string      `json:"region,omitempty"`
	Status                  string      `json:"status,omitempty"`
	SubscriptionExternalId  string      `json:"subscriptionExternalId,omitempty"`
	SubscriptionId          string      `json:"subscriptionId,omitempty"`
	SubscriptionName        string      `json:"subscriptionName,omitempty"`
	Tags                    interface{} `json:"tags,omitempty"`
	Type                    string      `json:"type"`
}

// Application service linked to an issue
type IssueApplicationService struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
}

// Risk factor of an issue
type IssueRiskFactor struct {
	Id     string `json:"id"`
	Policy struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"policy"`
	Severity string `json:"severity"`
}

// Project information
type IssueQueryProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Service account information
type IssueQueryServiceAccount struct {
	Id string `json:"id"`
}

// User information
type IssueQueryUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Service ticket object
type IssueServiceTickets struct {
	Action      IssueQueryAutomationAction `json:"action,omitempty"`
	ExternalId  string                     `json:"externalId,omitempty"`
	Id          string                     `json:"id,omitempty"`
	Integration IssueQueryIntegration      `json:"integration,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Project     IssueQueryProject          `json:"project,omitempty"`
	Url         string                     `json:"url,omitempty"`
}

// Issue object
type Issue struct {
	ApplicationServices        []IssueApplicationService `json:"applicationServices"`
	CreatedAt                  string                    `json:"createdAt"`
	Description                string                    `json:"description"`
	DueAt                      string                    `json:"dueAt"`
	Entity                     GraphEntity               `json:"entity"`
	EntitySnapshot             *IssueEntitySnapshot      `json:"entitySnapshot"`
	Environments               []string                  `json:"environments"`
	FirstEventAt               string                    `json:"firstEventAt"`
	HasCodeRemediation         bool                      `json:"hasCodeRemediation"`
	HasCodeToCloudMapping      bool                      `json:"hasCodeToCloudMapping"`
	Id                         string                    `json:"id"`
	LastEventAt                string                    `json:"lastEventAt"`
	Notes                      []IssueNote               `json:"notes"`
	OpenReason                 string                    `json:"openReason"`
	Projects                   []IssueQueryProject       `json:"projects"`
	RejectionExpiredAt         string                    `json:"rejectionExpiredAt"`
	RelatedApplicationServices []IssueApplicationService `json:"relatedApplicationServices"`
	ResolutionNote             string                    `json:"resolutionNote"`
	ResolutionReason           string                    `json:"resolutionReason"`
	ResolvedAt                 string                    `json:"resolvedAt"`
	RiskFactors                []IssueRiskFactor         `json:"riskFactors"`
	ServiceTickets             []IssueServiceTickets     `json:"serviceTickets"`
	Severity                   string                    `json:"severity"`
	SourceRules                []IssueSourceRule         `json:"sourceRules"`
	Status                     string                    `json:"status"`
	StatusChangedAt            string                    `json:"statusChangedAt"`
	Type                       string                    `json:"type"`
	UpdatedAt                  string                    `json:"updatedAt"`
	Url                        string                    `json:"url"`
	ValidatedAsExploitable     bool                      `json:"validatedAsExploitable"`
}

// Issue filter object
type IssueFilter struct {
	// Filter Issues created before/after provided date.
	CreatedAt *DateFilter

	// Filter issues using any of securityFramework | securitySubCategory | securityCategory.
	FrameworkCategory string

	// Filter issues that have (or don't have) a service ticket.
	HasServiceTicket *bool

	// Filter issues by project ID.
	ProjectId []string

	// The control severity.
	//
	// Possible values are: INFORMATIONAL, LOW, MEDIUM, HIGH, CRITICAL.
	Severity []string

	// Filter issues created by a specific source rule (control) ID.
	SourceControl []string

	// The issue status.
	//
	// Possible values are: OPEN, IN_PROGRESS, RESOLVED, REJECTED.
	Status []string

	// The issue type.
	//
	// Possible values are: TOXIC_COMBINATION, THREAT_DETECTION, CLOUD_CONFIGURATION, ATTACK_SURFACE, RISK_TOXIC_COMBINATION.
	Type []string

	// Filter issues by resolution reason.
	//
	// Possible values are: OBJECT_DELETED, ISSUE_FIXED, CONTROL_CHANGED, CONTROL_DISABLED, CONTROL_DELETED, FALSE_POSITIVE, EXCEPTION, WONT_FIX.
	ResolutionReason []string
}

// Relay-style node for the issue
type IssueConnection struct {
	Nodes      []Issue  `json:"nodes"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount"`
}

// ListIssuesResponse is returned by ListIssues on success
type ListIssuesResponse struct {
	Issues IssueConnection `json:"issuesV2"`
}

// Fields used to filter the issue response
type ListIssuesRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - filters object
	Filter *IssueFilter

	// Optional - also select application services and risk factors (slower).
	IncludeExtras bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int
}

// ListIssues returns a paginated list of all the issues
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListIssues(
	ctx context.Context,
	client *Client,
	options *ListIssuesRequestConfiguration,
) (*ListIssuesResponse, error) {
	// Make a request
	query := queryIssueList
	if options.IncludeExtras {
		query = queryIssueListWithExtras
	}
	req := graphql.NewRequest(query)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.Filter != nil {
		if options.Filter.FrameworkCategory != "" {
			filter["frameworkCategory"] = options.Filter.FrameworkCategory
		}
		if options.Filter.HasServiceTicket != nil {
			filter["hasServiceTicket"] = *options.Filter.HasServiceTicket
		}
		if len(options.Filter.ProjectId) > 0 {
			filter["project"] = options.Filter.ProjectId
		}
		if len(options.Filter.Severity) > 0 {
			filter["severity"] = options.Filter.Severity
		}
		if len(options.Filter.SourceControl) > 0 {
			filter["sourceRule"] = map[string]interface{}{"id": options.Filter.SourceControl}
		}
		if len(options.Filter.Status) > 0 {
			filter["status"] = options.Filter.Status
		}
		if len(options.Filter.Type) > 0 {
			filter["type"] = options.Filter.Type
		}
		if len(options.Filter.ResolutionReason) > 0 {
			filter["resolutionReason"] = options.Filter.ResolutionReason
		}

		if options.Filter.CreatedAt != nil {
			filter["createdAt"] = options.Filter.CreatedAt.ToVar()
		}
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
	var responseData ListIssuesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetIssueResponse is returned by GetIssue on success
type GetIssueResponse struct {
	Issue Issue `json:"issue"`
}

// GetIssue returns a specific issue that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetIssue(
	ctx context.Context,
	client *Client,
	id string,
	includeExtras bool,
) (*GetIssueResponse, error) {
	// Make a request
	query := queryIssueGet
	if includeExtras {
		query = queryIssueGetWithExtras
	}
	req := graphql.NewRequest(query)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetIssueResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
