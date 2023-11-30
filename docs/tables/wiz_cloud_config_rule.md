---
title: "Steampipe Table: wiz_cloud_config_rule - Query Wiz Cloud Config Rules using SQL"
description: "Allows users to query Wiz Cloud Config Rules, specifically the rules associated with cloud resources, providing insights into compliance and configuration patterns."
---

# Table: wiz_cloud_config_rule - Query Wiz Cloud Config Rules using SQL

Wiz Cloud Config Rules is a resource within Wiz that allows you to monitor and manage the configuration rules associated with your cloud resources. It provides a centralized way to set up and manage rules for various cloud resources, including virtual machines, databases, web applications, and more. Wiz Cloud Config Rules helps you ensure compliance and maintain desired configurations for your cloud resources.

## Table Usage Guide

The `wiz_cloud_config_rule` table provides insights into configuration rules within Wiz Cloud Config Rules. As a cloud engineer, explore rule-specific details through this table, including rule identifiers, descriptions, compliance types, and associated metadata. Utilize it to uncover information about rules, such as those related to specific resources, the compliance status of resources, and the verification of rule compliance.

## Examples

### Basic info
Gain insights into the status and severity of your cloud configuration rules across different providers. This can help you assess the robustness of your cloud security and compliance measures.

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule;
```

### List disabled rules
Uncover the details of inactive rules within your cloud configuration. This query is useful for identifying which rules have been disabled, allowing you to assess potential vulnerabilities and ensure optimal security settings.

```sql
select
  name,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  not enabled;
```

### List built-in rules
Explore which built-in rules are currently enabled, their severity level, and the specific cloud provider they apply to. This can help in understanding the existing security and compliance measures in place, and when they were created.

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  built_in;
```

### List high-severity rules specific to AWS S3 bucket
Identify instances where high-severity rules are applied to your AWS S3 buckets. This can help prioritize security measures and ensure the most critical areas of your cloud storage are adequately protected.

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  cloud_provider = 'AWS'
  and severity = 'HIGH'
  and target_native_types ?| array['bucket'];
```

### List all findings of a rule specific to AWS S3 bucket in the last 3 days
Determine the areas in which high-severity AWS S3 bucket rules have been violated in the past three days. This can provide insight into potential security risks and areas for improvement in your cloud configuration.

```sql
with list_s3_bucket_rules as (
  select
    name,
    id
  from
    wiz_cloud_config_rule
  where
    cloud_provider = 'AWS'
    and severity = 'HIGH'
    and target_native_types ?| array['bucket']
)
select
  r.name as rule_name,
  f.resource ->> 'Name' as resource_name,
  f.result as finding_status,
  f.analyzed_at
from
  wiz_cloud_configuration_finding as f
  join list_s3_bucket_rules as r on
    f.rule_id = r.id
    and f.severity = 'HIGH'
    and f.analyzed_at > (current_timestamp - interval '3 day');
```

### List rules with auto-remediation enabled
Discover the segments that have auto-remediation enabled, which helps in identifying rules that automatically correct violations, enhancing security and compliance within the cloud environment.

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  has_auto_remediation;
```