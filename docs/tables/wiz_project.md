---
title: "Steampipe Table: wiz_project - Query Wiz Projects using SQL"
description: "Allows users to query Wiz Projects, providing detailed information about the projects, such as their name, ID, and creation time."
---

# Table: wiz_project - Query Wiz Projects using SQL

Wiz is a cloud security platform that identifies the most critical risks and enables quick remediation by using a new scanner-less approach. It scans the entire environment to build a graph-based inventory and then applies cloud-native analysis to prioritize risks. Wiz supports multi-cloud environments and provides a holistic view of risks across Azure, AWS, GCP, and Kubernetes.

## Table Usage Guide

The `wiz_project` table provides insights into Projects within the Wiz platform. As a Security Engineer, explore project-specific details through this table, including ID, name, and creation time. Utilize it to uncover information about projects, such as their unique identifiers and the time they were created, which can aid in managing and securing your cloud environment.

## Examples

### Basic info
Gain insights into your project's structure by understanding the distribution of business units, the number of repositories, cloud accounts, and Kubernetes clusters, along with their respective security scores. This is useful to assess the overall security posture and resource allocation within your project.

```sql
select
  name,
  business_unit,
  security_score as security_score_in_percentage,
  repository_count,
  cloud_account_count,
  kubernetes_cluster_count
from
  wiz_project;
```

### Get count of critical issues per project
This query helps identify the number of critical issues per project, providing a clear overview of project health and potential areas of concern. This can be beneficial in prioritizing resources and remediation efforts.

```sql
with critical_issues as (
  select
    id,
    severity,
    p ->> 'id' as project
  from
    wiz_issue,
    jsonb_array_elements(projects) as p
  where
    severity = 'CRITICAL'
)
select
  p.name as project,
  count(c.id)
from
  wiz_project as p
  left join critical_issues as c on p.id = c.project
group by
  p.name;
```

### Get the owner details of each project
Explore which projects are owned by which users to better understand project responsibility distribution. This can assist in identifying the point of contact for each project, facilitating smoother communication and project management.

```sql
select
  p.name,
  p.slug,
  u.name as user_name,
  u.email as user_email
from
  wiz_project as p
  left join jsonb_array_elements(project_owners) as o on true
  left join wiz_user as u on u.id = o ->> 'id';
```

### List archived projects
Explore which projects have been archived, allowing you to assess elements like the associated business unit, security score, and linked resources such as repositories, cloud accounts, and Kubernetes clusters. This can be useful in understanding the scope and impact of archived projects within your organization.

```sql
select
  name,
  business_unit,
  security_score as security_score_in_percentage,
  repository_count,
  cloud_account_count,
  kubernetes_cluster_count
from
  wiz_project
where
  archived;
```