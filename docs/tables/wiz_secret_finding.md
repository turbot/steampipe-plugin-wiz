---
title: "Steampipe Table: wiz_secret_finding - Query Wiz Secret Findings using SQL"
description: "Allows users to query Wiz secret findings, providing details about exposed secrets such as cloud keys, passwords and private keys found in code repositories and cloud resources."
---

# Table: wiz_secret_finding - Query Wiz Secret Findings using SQL

Wiz detects secrets, such as cloud keys, passwords, private keys and SaaS API keys, in code repositories, virtual machines, container images and other cloud resources. Each finding records where the secret was found, the detection rule, and whether Wiz validated the secret as live.

## Table Usage Guide

The `wiz_secret_finding` table provides insights into exposed secrets detected by Wiz. As a security engineer, you can use it to find live secrets that need rotating, track secrets committed to code repositories, and follow up on their remediation.

**Important Notes**
- `service_tickets` and `has_service_ticket` list the service tickets (e.g. Jira) linked to the finding through its related issues (up to 50 related issues per finding) and posture issues. The tickets are only fetched when one of these columns is selected, which makes the query slower and caps the page size at 500; they cannot be used as filters by the API.
- The table returns finding metadata only. Secret values and code snippets are never selected.
- The table returns findings in all statuses (`OPEN`, `IN_PROGRESS`, `RESOLVED`, `REJECTED`) unless you filter on `status`. The Wiz API itself only returns open findings when no status is given.
- The table can return a large dataset; which can increase the query execution time. It is recommended that queries to this table should include (usually in the `where` clause) at least one of these columns:
  - `first_seen_at`
  - `is_default_branch`
  - `last_updated_at`
  - `origin`
  - `path`
  - `repository_id`
  - `resource_id`
  - `resource_type`
  - `rule_id`
  - `severity`
  - `status`
  - `type`
  - `validation_status`

## Examples

### Basic info
Explore the open secret findings and where they were found.

```sql+postgres
select
  name,
  type,
  severity,
  resource_type,
  resource ->> 'name' as resource_name,
  path,
  first_seen_at
from
  wiz_secret_finding
where
  status = 'OPEN';
```

```sql+sqlite
select
  name,
  type,
  severity,
  resource_type,
  json_extract(resource, '$.name') as resource_name,
  path,
  first_seen_at
from
  wiz_secret_finding
where
  status = 'OPEN';
```

### List secrets validated as live
Identify the secrets that must be rotated first because Wiz confirmed they are valid.

```sql+postgres
select
  name,
  type,
  resource ->> 'name' as resource_name,
  path,
  last_validated_at
from
  wiz_secret_finding
where
  validation_status = 'VALID'
  and status = 'OPEN';
```

```sql+sqlite
select
  name,
  type,
  json_extract(resource, '$.name') as resource_name,
  path,
  last_validated_at
from
  wiz_secret_finding
where
  validation_status = 'VALID'
  and status = 'OPEN';
```

### Count open secrets per code repository
Find the repositories with the most secrets committed to code.

```sql+postgres
select
  repository_name,
  count(*) as secrets
from
  wiz_secret_finding
where
  resource_type = 'REPOSITORY_BRANCH'
  and status = 'OPEN'
group by
  repository_name
order by
  secrets desc;
```

```sql+sqlite
select
  repository_name,
  count(*) as secrets
from
  wiz_secret_finding
where
  resource_type = 'REPOSITORY_BRANCH'
  and status = 'OPEN'
group by
  repository_name
order by
  secrets desc;
```

### List cloud keys found on virtual machines
Detect cloud access keys stored on compute resources.

```sql+postgres
select
  name,
  resource ->> 'name' as vm_name,
  resource -> 'cloudAccount' ->> 'name' as cloud_account,
  path
from
  wiz_secret_finding
where
  type = 'CLOUD_KEY'
  and resource_type = 'VIRTUAL_MACHINE'
  and status = 'OPEN';
```

```sql+sqlite
select
  name,
  json_extract(resource, '$.name') as vm_name,
  json_extract(resource, '$.cloudAccount.name') as cloud_account,
  path
from
  wiz_secret_finding
where
  type = 'CLOUD_KEY'
  and resource_type = 'VIRTUAL_MACHINE'
  and status = 'OPEN';
```

### List live secrets without a service ticket
Find validated secrets that nobody is tracking yet.

```sql+postgres
select
  name,
  type,
  repository_name,
  path
from
  wiz_secret_finding
where
  validation_status = 'VALID'
  and status = 'OPEN'
  and not has_service_ticket;
```

```sql+sqlite
select
  name,
  type,
  repository_name,
  path
from
  wiz_secret_finding
where
  validation_status = 'VALID'
  and status = 'OPEN'
  and has_service_ticket = 0;
```

### List the service tickets of open secret findings
See which tickets track each secret.

```sql+postgres
select
  f.name,
  f.repository_name,
  t ->> 'external_id' as ticket,
  t ->> 'url' as ticket_url
from
  wiz_secret_finding as f,
  jsonb_array_elements(f.service_tickets) as t
where
  f.status = 'OPEN';
```

```sql+sqlite
select
  f.name,
  f.repository_name,
  json_extract(t.value, '$.external_id') as ticket,
  json_extract(t.value, '$.url') as ticket_url
from
  wiz_secret_finding as f,
  json_each(f.service_tickets) as t
where
  f.status = 'OPEN';
```
