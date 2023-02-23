package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	queryCloudConfigRuleList = `
query ListConfigRules($first: Int, $after: String, $filter: CloudConfigurationRuleFilters) {
  cloudConfigurationRules(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      shortId
      enabled
      severity
      description
      supportsNRT
      cloudProvider
      serviceType
      builtin
      functionAsControl
      hasAutoRemediation
      remediationInstructions
      createdBy {
        id
      }
      createdAt
      updatedAt
      targetNativeTypes
      scopeAccounts {
        id
      }
      control {
        id
      }
    }
  }
}
`
	queryCloudConfigRuleGet = `
query GetCloudConfigRule($id: ID!) {
  cloudConfigurationRule(id: $id) {
    id
    name
    shortId
    enabled
    severity
    description
    supportsNRT
    cloudProvider
    serviceType
    builtin
    functionAsControl
    hasAutoRemediation
    remediationInstructions
    createdBy {
      id
    }
    createdAt
    updatedAt
    targetNativeTypes
    scopeAccounts {
      id
    }
    control {
      id
    } 
  }
}
`
)

// Cloud configuration rule object
type CloudConfigRule struct {
	BuiltIn                 bool                               `json:"builtin"`
	CloudProvider           string                             `json:"cloudProvider"`
	Control                 []CloudConfigRuleQueryControl      `json:"control"`
	CreatedAt               string                             `json:"createdAt"`
	CreatedBy               CloudConfigRuleQueryUser           `json:"createdBy"`
	Description             string                             `json:"description"`
	Enabled                 bool                               `json:"enabled"`
	FunctionAsControl       bool                               `json:"functionAsControl"`
	HasAutoRemediation      bool                               `json:"hasAutoRemediation"`
	Id                      string                             `json:"id"`
	Name                    string                             `json:"name"`
	RemediationInstructions string                             `json:"remediationInstructions"`
	ScopeAccounts           []CloudConfigRuleQuerySubscription `json:"scopeAccounts"`
	ServiceType             string                             `json:"serviceType"`
	Severity                string                             `json:"severity"`
	ShortId                 string                             `json:"shortId"`
	SupportsNrt             bool                               `json:"supportsNRT"`
	TargetNativeTypes       []string                           `json:"targetNativeTypes"`
	UpdatedAt               string                             `json:"updatedAt"`
}

// Control information
type CloudConfigRuleQueryControl struct {
	Id string `json:"id"`
}

// Cloud account information
type CloudConfigRuleQuerySubscription struct {
	Id string `json:"id"`
}

// User information
type CloudConfigRuleQueryUser struct {
	Id string `json:"id"`
}

// Relay-style node for the cloud configuration rule
type CloudConfigRuleConnection struct {
	Nodes      []CloudConfigRule `json:"nodes"`
	PageInfo   PageInfo          `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}

// ListCloudConfigRulesResponse is returned by ListCloudConfigRules on success
type ListCloudConfigRulesResponse struct {
	CloudConfigRules CloudConfigRuleConnection `json:"cloudConfigurationRules"`
}

// Fields used to filter the cloud configuration rule response
type ListCloudConfigRulesRequestConfiguration struct {
	// Optional - filter CSPM rules by to cloud provider.
	//
	// Possible values are: AWS, GCP, OCI, Alibaba, Azure, Kubernetes, OpenShift, vSphere.
	CloudProvider string

	// Optional - filter CSPM Rule by enabled status.
	Enabled *bool

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// Optional - if true, all the rule that has auto remediation enabled will be returned.
	HasAutoRemediation *bool

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter CSPM rules related by service type.
	//
	// Possible values: AKS, AWS, Alibaba, Azure, EKS, GCP, GKE, Kubernetes, OCI, OKE, vSphere.
	ServiceType string

	// Optional - filter CSPM Rule by severity.
	//
	// Possible values are: CRITICAL, HIGH, INFORMATIONAL, LOW, MEDIUM.
	Severity string

	// Optional - if true, rule with support of "near real time" updates will be returned.
	SupportsNRT *bool
}

// ListCloudConfigRules returns a paginated list of the cloud configuration rules
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListCloudConfigRules(
	ctx context.Context,
	client *Client,
	options *ListCloudConfigRulesRequestConfiguration,
) (*ListCloudConfigRulesResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryCloudConfigRuleList)

	// Check for optional filters
	filter := map[string]interface{}{}
	if options.CloudProvider != "" {
		filter["cloudProvider"] = options.CloudProvider
	}
	if options.Enabled != nil {
		filter["enabled"] = *options.Enabled
	}
	if options.HasAutoRemediation != nil {
		filter["hasAutoRemediation"] = *options.HasAutoRemediation
	}
	if options.ServiceType != "" {
		filter["serviceType"] = options.ServiceType
	}
	if options.Severity != "" {
		filter["severity"] = options.Severity
	}
	if options.SupportsNRT != nil {
		filter["supportsNRT"] = *options.SupportsNRT
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
	var responseData ListCloudConfigRulesResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetCloudConfigRuleResponse is returned by GetCloudConfigRule on success
type GetCloudConfigRuleResponse struct {
	CloudConfigRule CloudConfigRule `json:"cloudConfigurationRule"`
}

// GetCloudConfigRule returns a specific CSPM rule that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetCloudConfigRule(
	ctx context.Context,
	client *Client,
	id string,
) (*GetCloudConfigRuleResponse, error) {
	// Make a request
	req := graphql.NewRequest(queryCloudConfigRuleGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetCloudConfigRuleResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
