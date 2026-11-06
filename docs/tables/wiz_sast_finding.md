---
title: "Steampipe Table: wiz_sast_finding - Query Wiz SAST Findings using SQL"
description: "Allows users to query Wiz SAST (Static Application Security Testing) findings, providing details about first-party code weaknesses detected in code repositories."
---

# Table: wiz_sast_finding - Query Wiz SAST Findings using SQL

Wiz Code scans source code repositories for first-party weaknesses (SAST), such as injection flaws or insecure cryptography. Each finding records the repository, branch, file and line where the weakness was found, the rule and CWE weaknesses that apply, and an optional AI analysis of whether the weakness is exploitable. Findings reported by integrated third-party scanners are included too, with the scanner shown in `origin`.

## Table Usage Guide

The `wiz_sast_finding` table provides insights into code weaknesses found by Wiz Code. As an application security engineer, you can use it to track open findings per repository, prioritize findings that the AI analysis marks as exploitable, and measure how quickly findings are resolved.

**Important Notes**
- `service_tickets` and `has_service_ticket` list the service tickets (e.g. Jira) linked to the finding through its related issues (up to 50 related issues per finding) and posture issues. The tickets are only fetched when one of these columns is selected, which makes the query slower and caps the page size at 500; they cannot be used as filters by the API.
- The table returns findings in all statuses (`OPEN`, `RESOLVED`, `REJECTED`) unless you filter on `status`. The Wiz API treats `IN_PROGRESS` as `OPEN` for findings, so the table never returns it, and a filter on `IN_PROGRESS` alone returns no rows. The Wiz API itself only returns open findings when no status is given.
- The table can return a large dataset; which can increase the query execution time. It is recommended that queries to this table should include (usually in the `where` clause) at least one of these columns:
  - `ai_verdict`
  - `created_at`
  - `file_path`
  - `is_ai_powered`
  - `is_default_branch`
  - `origin`
  - `repository_branch_id`
  - `repository_id`
  - `severity`
  - `status`
  - `updated_at`
- `is_default_branch` is a filter-only column: it supports only the `=` operator (e.g. `is_default_branch = true`), and is only populated when used that way in the `where` clause.

## Examples

### Basic info
Explore the open SAST findings, the repository and file they were found in, and their severity.

```sql+postgres
select
  name,
  severity,
  repository_name,
  file_path,
  start_line,
  created_at
from
  wiz_sast_finding
where
  status = 'OPEN';
```

```sql+sqlite
select
  name,
  severity,
  repository_name,
  file_path,
  start_line,
  created_at
from
  wiz_sast_finding
where
  status = 'OPEN';
```

### Count open findings per repository on the default branch
Identify the repositories with the most open code weaknesses on their default branch.

```sql+postgres
select
  repository_name,
  count(*) filter (where severity in ('CRITICAL', 'HIGH')) as critical_high,
  count(*) as total
from
  wiz_sast_finding
where
  status = 'OPEN'
  and is_default_branch = true
group by
  repository_name
order by
  critical_high desc,
  total desc;
```

```sql+sqlite
select
  repository_name,
  sum(case when severity in ('CRITICAL', 'HIGH') then 1 else 0 end) as critical_high,
  count(*) as total
from
  wiz_sast_finding
where
  status = 'OPEN'
  and is_default_branch = 1
group by
  repository_name
order by
  critical_high desc,
  total desc;
```

### List findings that the Wiz AI analysis marks as exploitable
Focus remediation on the weaknesses most likely to be exploitable.

```sql+postgres
select
  name,
  repository_name,
  file_path,
  ai_analysis ->> 'confidence' as confidence,
  ai_analysis ->> 'explanation' as explanation
from
  wiz_sast_finding
where
  ai_verdict = 'EXPLOITABLE'
  and status = 'OPEN';
```

```sql+sqlite
select
  name,
  repository_name,
  file_path,
  json_extract(ai_analysis, '$.confidence') as confidence,
  json_extract(ai_analysis, '$.explanation') as explanation
from
  wiz_sast_finding
where
  ai_verdict = 'EXPLOITABLE'
  and status = 'OPEN';
```

### List the CWE weaknesses of open critical findings
Understand which weakness classes are most common among critical findings.

```sql+postgres
select
  w ->> 'id' as cwe,
  w ->> 'name' as weakness,
  count(*) as findings
from
  wiz_sast_finding,
  jsonb_array_elements(weaknesses) as w
where
  severity = 'CRITICAL'
  and status = 'OPEN'
group by
  cwe,
  weakness
order by
  findings desc;
```

```sql+sqlite
select
  json_extract(w.value, '$.id') as cwe,
  json_extract(w.value, '$.name') as weakness,
  count(*) as findings
from
  wiz_sast_finding,
  json_each(weaknesses) as w
where
  severity = 'CRITICAL'
  and status = 'OPEN'
group by
  cwe,
  weakness
order by
  findings desc;
```

### List findings resolved in the last 30 days
Track remediation progress over time.

```sql+postgres
select
  name,
  repository_name,
  resolution_reason,
  updated_at
from
  wiz_sast_finding
where
  status = 'RESOLVED'
  and updated_at > now() - interval '30 days';
```

```sql+sqlite
select
  name,
  repository_name,
  resolution_reason,
  updated_at
from
  wiz_sast_finding
where
  status = 'RESOLVED'
  and updated_at > datetime('now', '-30 days');
```

### List open high severity findings without a service ticket
Find code weaknesses that nobody is tracking yet.

```sql+postgres
select
  name,
  repository_name,
  file_path
from
  wiz_sast_finding
where
  severity in ('CRITICAL', 'HIGH')
  and status = 'OPEN'
  and not has_service_ticket;
```

```sql+sqlite
select
  name,
  repository_name,
  file_path
from
  wiz_sast_finding
where
  severity in ('CRITICAL', 'HIGH')
  and status = 'OPEN'
  and has_service_ticket = 0;
```

### List the service tickets of open critical findings
See which tickets track each finding.

```sql+postgres
select
  f.name,
  f.repository_name,
  t ->> 'external_id' as ticket,
  t ->> 'integration' as integration
from
  wiz_sast_finding as f,
  jsonb_array_elements(f.service_tickets) as t
where
  f.severity = 'CRITICAL'
  and f.status = 'OPEN';
```

```sql+sqlite
select
  f.name,
  f.repository_name,
  json_extract(t.value, '$.external_id') as ticket,
  json_extract(t.value, '$.integration') as integration
from
  wiz_sast_finding as f,
  json_each(f.service_tickets) as t
where
  f.severity = 'CRITICAL'
  and f.status = 'OPEN';
```
