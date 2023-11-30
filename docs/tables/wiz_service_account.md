---
title: "Steampipe Table: wiz_service_account - Query Wiz Service Accounts using SQL"
description: "Allows users to query Wiz Service Accounts, specifically the details of each service account, providing insights into account configurations and potential security risks."
---

# Table: wiz_service_account - Query Wiz Service Accounts using SQL

Wiz Service Accounts represent the accounts used by the Wiz platform to conduct security assessments and data collection in your environment. Each service account is associated with a specific cloud environment and has specific permissions and roles assigned to it. Understanding these service accounts can help in managing permissions and ensuring proper security controls.

## Table Usage Guide

The `wiz_service_account` table provides insights into service accounts within the Wiz platform. As a security engineer, explore service account-specific details through this table, including permissions, associated cloud environments, and account roles. Utilize it to understand the permission scope of each service account, identify any accounts with overly broad permissions, and verify the security posture of your environment.

## Examples

### Basic info
Explore which service accounts were created and when they were last updated to gain insights into your account's security protocols and potential vulnerabilities. This query is particularly useful for auditing purposes, allowing you to identify instances where outdated or unused service accounts may pose a security risk.

```sql
select
  name,
  client_id,
  type,
  created_at,
  last_rotated_at,
  authentication_source
from
  wiz_service_account;
```

### List service accounts scoped to a specific project
Explore which service accounts are designated to a particular project. This can help you understand the distribution of service accounts across your projects, aiding in efficient resource allocation and management.

```sql
select
  s.name,
  p.name as project,
  s.client_id,
  s.created_at,
  s.last_rotated_at,
  s.scopes
from
  wiz_service_account as s,
  jsonb_array_elements(s.assigned_projects) as ap
  join wiz_project as p on p.id = ap ->> 'id';
```

### List service accounts by name
Explore which service accounts are associated with a particular name to understand their creation and last rotation dates, as well as their scopes. This can be useful for auditing purposes, to ensure that accounts are being properly managed and updated.

```sql
select
  name,
  client_id,
  created_at,
  last_rotated_at,
  scopes
from
  wiz_service_account
where
  name = 'Steampipe';
```

### List service accounts that have not been rotated in the last 90 days
Identify service accounts that have been dormant for the past 90 days. This could be useful in maintaining the security of your system by ensuring regular rotation of service accounts.

```sql
select
  name,
  client_id,
  created_at,
  last_rotated_at,
  scopes
from
  wiz_service_account
where
  last_rotated_at < (current_timestamp - interval '90 days');
```

### List users scoped to a specific project
Discover the segments that include users assigned to a specific project. This can be particularly useful for project managers who need to quickly identify all users linked to their project, enabling efficient management and communication.

```sql
select
  u.name,
  u.email,
  u.identity_provider_type,
  u.role ->> 'name' as role,
  p.name as project,
  u.created_at
from
  wiz_user as u,
  jsonb_array_elements(u.effective_assigned_projects) as ep
  join wiz_project as p on ep ->> 'id' = p.id;
```