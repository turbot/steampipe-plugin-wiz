package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Define the query
const (
	querySubscriptionList = `
query ListSubscriptions($first: Int, $after: String, $filter: CloudAccountFilters) {
  cloudAccounts(first: $first, after: $after, filterBy: $filter) {
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
    nodes {
      id
      name
      cloudProvider
      status
      resourceCount
      containerCount
      virtualMachineCount
      externalId
      lastScannedAt
      linkedProjects {
        id
      }
    }
  }
}
`
	querySubscriptionGet = `
query GetSubscription($id: ID!) {
  cloudAccount(id: $id) {
    id
    name
    cloudProvider
    status
    resourceCount
    containerCount
    virtualMachineCount
    externalId
    lastScannedAt
    linkedProjects {
      id
    }
  }
}
`
)

// Subscription object
type Subscription struct {
	CloudProvider       string                     `json:"cloudProvider"`
	ContainerCount      int                        `json:"containerCount"`
	ExternalId          string                     `json:"externalId"`
	Id                  string                     `json:"id"`
	LastScannedAt       string                     `json:"lastScannedAt"`
	LinkedProjects      []SubscriptionQueryProject `json:"linkedProjects"`
	Name                string                     `json:"name"`
	ResourceCount       int                        `json:"resourceCount"`
	Status              string                     `json:"status"`
	VirtualMachineCount int                        `json:"virtualMachineCount"`
}

// Assigned project information
type SubscriptionQueryProject struct {
	Id string `json:"id"`
}

// Relay-style node for the subscription
type SubscriptionConnection struct {
	Nodes      []Subscription `json:"nodes"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

// ListSubscriptionsResponse is returned by ListSubscriptions on success
type ListSubscriptionsResponse struct {
	Subscriptions SubscriptionConnection `json:"cloudAccounts"`
}

// Fields used to filter the subscription response
type ListSubscriptionsRequestConfiguration struct {
	// Optional - filter subscriptions of specific cloud provider.
	//
	// Possible values are: AWS, GCP, OCI, Alibaba, Azure, Kubernetes, OpenShift, vSphere.
	CloudProvider string

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 500.
	Limit int

	// Optional - filter subscriptions by it's status.
	//
	// Possible values: CONNECTED, DISABLED, DISCONNECTED, DISCOVERED, ERROR, INITIAL_SCANNING, PARTIALLY_CONNECTED.
	Status string
}

// ListSubscriptions returns a paginated list of the cloud account subscriptions
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListSubscriptions(
	ctx context.Context,
	client *Client,
	options *ListSubscriptionsRequestConfiguration,
) (*ListSubscriptionsResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySubscriptionList)

	// Check for optional filters
	filter := map[string]string{}
	if options.CloudProvider != "" {
		filter["cloudProvider"] = options.CloudProvider
	}
	if options.Status != "" {
		filter["status"] = options.Status
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
	var responseData ListSubscriptionsResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}

// GetSubscriptionResponse is returned by GetSubscription on success
type GetSubscriptionResponse struct {
	Subscription Subscription `json:"cloudAccount"`
}

// GetSubscription returns a specific subscription that matches the ID
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the resource
func GetSubscription(
	ctx context.Context,
	client *Client,
	id string,
) (*GetSubscriptionResponse, error) {
	// Make a request
	req := graphql.NewRequest(querySubscriptionGet)

	// Set the required variables
	req.Var("id", id)

	// execute api call
	var responseData GetSubscriptionResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
