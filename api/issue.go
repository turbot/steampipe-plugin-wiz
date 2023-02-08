package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryIssueList = `
query ListIssues($first: Int, $after: String, $filter: IssueFilters) {
  issues(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      description
      status
      severity
      createdAt
      updatedAt
      resolvedAt
      dueAt
      rejectionExpiredAt
      statusChangedAt
      resolutionReason
      control {
        id
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
        action {
          id
        }
        integration {
          id
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
  }
}
`
	queryIssueGet = `
query GetIssue($id: ID!) {
  issue(id: $id) {
    id
    description
    status
    severity
    createdAt
    updatedAt
    resolvedAt
    dueAt
    rejectionExpiredAt
    statusChangedAt
    resolutionReason
    control {
      id
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
      action {
        id
      }
      integration {
        id
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
}
`
)

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

// Control information
type IssueQueryControl struct {
	Id string `json:"id"`
}

// Integration information
type IssueQueryIntegration struct {
	Id string `json:"id"`
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
	Control            IssueQueryControl     `json:"control"`
	CreatedAt          string                `json:"createdAt"`
	Description        string                `json:"description"`
	DueAt              string                `json:"dueAt"`
	Id                 string                `json:"id"`
	Notes              []IssueNote           `json:"notes"`
	Projects           []IssueQueryProject   `json:"projects"`
	RejectionExpiredAt string                `json:"rejectionExpiredAt"`
	ResolutionReason   string                `json:"resolutionReason"`
	ResolvedAt         string                `json:"resolvedAt"`
	ServiceTickets     []IssueServiceTickets `json:"serviceTickets"`
	Severity           string                `json:"severity"`
	Status             string                `json:"status"`
	StatusChangedAt    string                `json:"statusChangedAt"`
	UpdatedAt          string                `json:"updatedAt"`
}

// Issue filter object
type IssueFilter struct {
	// The control severity.
	//
	// Possible values are: INFORMATIONAL, LOW, MEDIUM, HIGH, CRITICAL.
	Severity string

	// The issue status.
	//
	// Possible values are: OPEN, IN_PROGRESS, RESOLVED, REJECTED.
	Status string

	// Filter issues by resolution reason.
	//
	// Possible values are: OBJECT_DELETED, ISSUE_FIXED, CONTROL_CHANGED, CONTROL_DISABLED, CONTROL_DELETED, FALSE_POSITIVE, EXCEPTION, WONT_FIX.
	ResolutionReason string
}

// Relay-style node for the issue
type IssueConnection struct {
	Nodes      []Issue  `json:"nodes"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount"`
}

// ListIssuesResponse is returned by ListIssues on success
type ListIssuesResponse struct {
	Issues IssueConnection `json:"issues"`
}

// Fields used to filter the issue response
type ListIssuesRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - filters object
	Filter *IssueFilter

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
	req := graphql.NewRequest(queryIssueList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.Filter != nil {
		if options.Filter.Severity != "" {
			filter["severity"] = options.Filter.Severity
		}
		if options.Filter.Status != "" {
			filter["status"] = options.Filter.Status
		}
		if options.Filter.ResolutionReason != "" {
			filter["resolutionReason"] = options.Filter.ResolutionReason
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
) (*GetIssueResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryIssueGet)

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
