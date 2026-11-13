package wiz

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-wiz/api"
)

//// TABLE DEFINITION

func tableWizKubernetesCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "wiz_kubernetes_cluster",
		Description: "Wiz Kubernetes Cluster",
		List: &plugin.ListConfig{
			Hydrate: listWizKubernetesClusters,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the cluster."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The Wiz identifier of the cluster."},
			{Name: "external_id", Type: proto.ColumnType_STRING, Description: "The identifier of the cluster in its cloud provider.", Transform: transform.FromP(clusterProperty, "externalId")},
			{Name: "provider_unique_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the cluster in its cloud provider (for example an ARN)."},
			{Name: "cloud_platform", Type: proto.ColumnType_STRING, Description: "The cloud platform hosting the cluster, such as AWS, GCP, Azure or OpenShift.", Transform: transform.FromP(clusterProperty, "cloudPlatform")},
			{Name: "native_type", Type: proto.ColumnType_STRING, Description: "The cluster type as named by the cloud provider.", Transform: transform.FromP(clusterProperty, "nativeType")},
			{Name: "kubernetes_flavor", Type: proto.ColumnType_STRING, Description: "The Kubernetes distribution, such as EKS, GKE or AKS.", Transform: transform.FromP(clusterProperty, "kubernetes_kubernetesFlavor")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region of the cluster.", Transform: transform.FromP(clusterProperty, "region")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the cluster, such as Active or Inactive.", Transform: transform.FromP(clusterProperty, "status")},
			{Name: "subscription_external_id", Type: proto.ColumnType_STRING, Description: "The external identifier of the cloud account (AWS account, GCP project, Azure subscription) the cluster runs in.", Transform: transform.FromP(clusterProperty, "subscriptionExternalId")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The Kubernetes version of the cluster.", Transform: transform.FromP(clusterProperty, "version")},
			{Name: "is_version_end_of_life", Type: proto.ColumnType_BOOL, Description: "True, if the Kubernetes version of the cluster is end of life.", Transform: transform.FromP(clusterProperty, "isVersionEndOfLife")},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The time the cluster was created.", Transform: transform.FromP(clusterProperty, "creationDate")},
			{Name: "sensor_deployment_status", Type: proto.ColumnType_STRING, Description: "The deployment status of the Wiz runtime sensor on the cluster, such as Installed or NotInstalled.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_sensor_deploymentStatus")},
			{Name: "sensor_installed", Type: proto.ColumnType_BOOL, Description: "True, if the Wiz runtime sensor is installed on the cluster.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_sensor_installed")},
			{Name: "sensor_version", Type: proto.ColumnType_STRING, Description: "The version of the Wiz runtime sensor running on the cluster.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_sensor_version")},
			{Name: "sensor_is_end_of_life", Type: proto.ColumnType_BOOL, Description: "True, if the Wiz runtime sensor version on the cluster is end of life.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_sensor_isEndOfLife")},
			{Name: "admission_controller_deployment_status", Type: proto.ColumnType_STRING, Description: "The deployment status of the Wiz admission controller on the cluster.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_admissionController_deploymentStatus")},
			{Name: "connector_deployment_status", Type: proto.ColumnType_STRING, Description: "The deployment status of the Wiz Kubernetes connector on the cluster.", Transform: transform.FromP(clusterProperty, "deploymentCoverage_connector_deploymentStatus")},
			{Name: "first_seen", Type: proto.ColumnType_TIMESTAMP, Description: "The first time Wiz saw the cluster."},
			{Name: "last_seen", Type: proto.ColumnType_TIMESTAMP, Description: "The last time Wiz saw the cluster."},
			{Name: "projects", Type: proto.ColumnType_JSON, Description: "The projects the cluster belongs to."},
			{Name: "properties", Type: proto.ColumnType_JSON, Description: "All properties of the cluster in the Wiz Security Graph."},
		},
	}
}

//// LIST FUNCTION

func listWizKubernetesClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("wiz_kubernetes_cluster.listWizKubernetesClusters", "connection_error", err)
		return nil, err
	}

	options := &api.ListKubernetesClustersRequestConfiguration{}

	// The Security Graph returns at most 500 entities per page
	pageLimit := 500
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	for {
		query, err := api.ListKubernetesClusters(ctx, conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("wiz_kubernetes_cluster.listWizKubernetesClusters", "query_error", err)
			return nil, err
		}

		for _, node := range query.GraphSearch.Nodes {
			for _, cluster := range node.Entities {
				if cluster == nil {
					continue
				}
				d.StreamListItem(ctx, cluster)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

		if !query.GraphSearch.PageInfo.HasNextPage {
			break
		}
		options.EndCursor = query.GraphSearch.PageInfo.EndCursor
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

// clusterProperty returns one key of the cluster's Security Graph properties
func clusterProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cluster, ok := d.HydrateItem.(*api.KubernetesCluster)
	if !ok || cluster.Properties == nil {
		return nil, nil
	}
	return cluster.Properties[d.Param.(string)], nil
}
