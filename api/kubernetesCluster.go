package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Kubernetes clusters are read from the Security Graph (graphSearch), which
// is the only API that lists every cluster Wiz knows about, connected or not.
// Each entity carries its properties as a JSON object, including the
// deployment coverage of the Wiz sensor, admission controller and connector.
const (
	queryKubernetesClusterList = `
query ListKubernetesClusters($first: Int, $after: String) {
  graphSearch(
    first: $first
    after: $after
    query: { type: [KUBERNETES_CLUSTER] }
  ) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      entities {
        id
        name
        providerUniqueId
        firstSeen
        lastSeen
        properties
        projects {
          id
          name
        }
      }
    }
  }
}
`
)

// KubernetesClusterProject is a project the cluster belongs to
type KubernetesClusterProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// KubernetesCluster is a Security Graph entity of type KUBERNETES_CLUSTER
type KubernetesCluster struct {
	Id               string                     `json:"id"`
	Name             string                     `json:"name"`
	ProviderUniqueId string                     `json:"providerUniqueId"`
	FirstSeen        string                     `json:"firstSeen"`
	LastSeen         string                     `json:"lastSeen"`
	Properties       map[string]interface{}     `json:"properties"`
	Projects         []KubernetesClusterProject `json:"projects"`
}

// KubernetesClusterSearchNode is one graphSearch result row
type KubernetesClusterSearchNode struct {
	Entities []*KubernetesCluster `json:"entities"`
}

// KubernetesClusterConnection is the paginated graphSearch result
type KubernetesClusterConnection struct {
	Nodes    []KubernetesClusterSearchNode `json:"nodes"`
	PageInfo PageInfo                      `json:"pageInfo"`
}

// ListKubernetesClustersResponse is returned by ListKubernetesClusters on success
type ListKubernetesClustersResponse struct {
	GraphSearch KubernetesClusterConnection `json:"graphSearch"`
}

// ListKubernetesClustersRequestConfiguration holds the request options
type ListKubernetesClustersRequestConfiguration struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The maximum number of results to return in a single call.
	Limit int
}

// ListKubernetesClusters returns a paginated list of Kubernetes clusters
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListKubernetesClusters(
	ctx context.Context,
	client *Client,
	options *ListKubernetesClustersRequestConfiguration,
) (*ListKubernetesClustersResponse, error) {
	req := graphql.NewRequest(queryKubernetesClusterList)

	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}
	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	var responseData ListKubernetesClustersResponse
	err := client.doRequest(req, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, err
}
