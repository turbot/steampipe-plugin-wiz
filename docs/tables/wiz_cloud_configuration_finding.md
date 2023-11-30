---
title: "Steampipe Table: wiz_cloud_configuration_finding - Query Wiz Cloud Configuration Findings using SQL"
description: "Allows users to query Wiz Cloud Configuration Findings, specifically providing insights into the cloud environment's security posture, including misconfigurations, compliance status, and potential vulnerabilities."
---

# Table: wiz_cloud_configuration_finding - Query Wiz Cloud Configuration Findings using SQL

Wiz Cloud Configuration Findings is a resource within Wiz that provides a comprehensive view of the security posture of your cloud environment. It identifies misconfigurations, compliance status, and potential vulnerabilities across various cloud resources. This resource helps you stay informed about the health and security of your cloud resources and take appropriate actions when security issues are detected.

## Table Usage Guide

The `wiz_cloud_configuration_finding` table provides insights into your cloud environment's security posture. As a Security Analyst, explore specific details through this table, including misconfigurations, compliance status, and potential vulnerabilities. Utilize it to uncover information about security issues, such as those related to misconfigurations, the compliance status of resources, and the identification of potential vulnerabilities.

## Examples

### Basic info
Explore the findings of your cloud configuration analysis. This can help you assess the severity and status of any potential issues, allowing for timely and effective resolution.

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding;
```

### List all failed resources with high severity
Gain insights into high-risk areas by identifying system resources with a high severity failure status. This can aid in prioritizing and addressing critical issues promptly for optimized system performance and security.

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and severity = 'HIGH';
```

### List failed resources which are not resolved
Identify unresolved, high-severity issues within your cloud configuration to prioritize and address potential vulnerabilities or misconfigurations.

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and status = 'OPEN'
  and severity = 'HIGH';
```

### List all findings detected in the last 3 days
Discover the segments that have been flagged with high severity issues in the last three days. This query helps in identifying unresolved, high-risk problems for prioritized troubleshooting.

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and status = 'OPEN'
  and severity = 'HIGH'
  and analyzed_at > (current_timestamp - interval '3 day');
```

### List failed resources with rule information
Identify instances where high severity resources have failed, including when the failure occurred and its current status. This information can help prioritize remediation efforts for resources that are not meeting compliance rules.

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding as f
  left join wiz_cloud_config_rule as r on
    f.rule ->> 'id' = r.id
    and f.result = 'FAIL'
    and f.status = 'OPEN'
    and f.severity = 'HIGH';
```