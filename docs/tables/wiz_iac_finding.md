---
title: "Steampipe Table: wiz_iac_finding - Query Wiz IaC Findings using SQL"
description: "Allows users to query Wiz Infrastructure as Code (IaC) findings, providing details about misconfigurations detected in Terraform, CloudFormation, Kubernetes manifests and other IaC files."
---

# Table: wiz_iac_finding - Query Wiz IaC Findings using SQL

Wiz scans Infrastructure as Code (IaC) files in code repositories, such as Terraform, CloudFormation, Kubernetes manifests, Dockerfiles and CI pipelines, against its cloud configuration rules. Each finding records the file and lines that break a rule, the content that was found and expected, and the remediation instructions.

## Table Usage Guide

The `wiz_iac_finding` table provides insights into misconfigurations detected in IaC files before they are deployed. As a cloud security engineer, you can use it to track misconfigurations per repository, find the rules broken most often, and follow up on their remediation.

**Important Notes**
- The table returns findings in all statuses (`OPEN`, `IN_PROGRESS`, `RESOLVED`, `REJECTED`) unless you filter on `status`.
- The table can return a large dataset; which can increase the query execution time. It is recommended that queries to this table should include (usually in the `where` clause) at least one of these columns:
  - `file_path`
  - `is_default_branch`
  - `platform`
  - `repository_id`
  - `severity`
  - `status`
  - `updated_at`
- `is_default_branch` is a filter-only column: it supports only the `=` operator (e.g. `is_default_branch = true`), and is only populated when used that way in the `where` clause.

## Examples

### Basic info
Explore the open IaC findings and the files they were found in.

```sql+postgres
select
  name,
  severity,
  platform,
  repository_name,
  file_path,
  start_line
from
  wiz_iac_finding
where
  status = 'OPEN';
```

```sql+sqlite
select
  name,
  severity,
  platform,
  repository_name,
  file_path,
  start_line
from
  wiz_iac_finding
where
  status = 'OPEN';
```

### List high severity Terraform findings on the default branch
Focus on the misconfigurations that will be deployed.

```sql+postgres
select
  name,
  repository_name,
  file_path,
  remediation_instructions
from
  wiz_iac_finding
where
  platform = 'TERRAFORM'
  and severity in ('CRITICAL', 'HIGH')
  and is_default_branch = true
  and status = 'OPEN';
```

```sql+sqlite
select
  name,
  repository_name,
  file_path,
  remediation_instructions
from
  wiz_iac_finding
where
  platform = 'TERRAFORM'
  and severity in ('CRITICAL', 'HIGH')
  and is_default_branch = 1
  and status = 'OPEN';
```

### Count open findings per rule
Identify the rules broken most often across your IaC.

```sql+postgres
select
  rule ->> 'name' as rule_name,
  rule ->> 'severity' as rule_severity,
  count(*) as findings
from
  wiz_iac_finding
where
  status = 'OPEN'
group by
  rule_name,
  rule_severity
order by
  findings desc;
```

```sql+sqlite
select
  json_extract(rule, '$.name') as rule_name,
  json_extract(rule, '$.severity') as rule_severity,
  count(*) as findings
from
  wiz_iac_finding
where
  status = 'OPEN'
group by
  rule_name,
  rule_severity
order by
  findings desc;
```
