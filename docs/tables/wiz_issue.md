---
title: "Steampipe Table: wiz_issue - Query Wiz Issues using SQL"
description: "Allows users to query Wiz Issues, providing detailed information about security risks and vulnerabilities identified in their environment."
---

# Table: wiz_issue - Query Wiz Issues using SQL

Wiz is a Cloud Security Posture Management (CSPM) tool that provides continuous security posture monitoring for cloud environments. It identifies security risks and vulnerabilities across a wide range of categories, including misconfigurations, policy violations, and threats. Wiz provides a holistic view of your security posture, enabling you to identify and remediate issues quickly and effectively.

## Table Usage Guide

The `wiz_issue` table provides insights into the security issues identified by Wiz in your cloud environment. As a security engineer, you can use this table to explore detailed information about each issue, including its severity, status, and associated resources. This can help you prioritize remediation efforts and improve your overall security posture.

**Important Notes**
- The table can return a large dataset; which can increase the query execution time. It is recommended that queries to this table should include (usually in the `where` clause) at least one of these columns:
  - `control_id`
  - `created_at`
  - `framework_category_id`
  - `has_service_ticket`
  - `project_id`
  - `resolution_reason`
  - `severity`
  - `status`
  - `type`
- `project_id` and `framework_category_id` are filter-only columns, and `control_id` matches one of the issue's source rules. Use them with the `=` operator and a single value. With an `in` list, Steampipe makes one API call per value, so an issue that matches several values is returned once per value.
- `application_services`, `related_application_services` and `risk_factors` are only fetched when selected, which makes the query slower.

## Examples

### Basic info
Explore which issues have been logged in your system, their severity and status, and when they were created. This can help you understand the range and depth of problems encountered, and the reasons provided for their resolution.

```sql+postgres
select
  id,
  status,
  severity,
  created_at,
  resolution_reason
from
  wiz_issue;
```

```sql+sqlite
select
  id,
  status,
  severity,
  created_at,
  resolution_reason
from
  wiz_issue;
```

### List critical issues
Pinpoint the specific instances where critical issues have arisen. This can assist in prioritizing problem-solving efforts and focusing on the most significant challenges first.

```sql+postgres
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'CRITICAL';
```

```sql+sqlite
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'CRITICAL';
```

### List high severity open issues
Discover the segments that contain high severity issues that are still open. This can be particularly useful in prioritizing and addressing critical issues promptly to minimize potential impacts.

```sql+postgres
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'HIGH'
  and status = 'OPEN';
```

```sql+sqlite
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'HIGH'
  and status = 'OPEN';
```

### List data security open issues using framework category ID
Explore open data security issues, categorized by their severity level, within a specific framework category. This helps in understanding the distribution of issues and prioritizing them for resolution.

```sql+postgres
select
  severity,
  count(id)
from
  wiz_issue
where
  status = 'OPEN'
  and framework_category_id = 'wct-id-422'
group by
  severity;
```

```sql+sqlite
select
  severity,
  count(id)
from
  wiz_issue
where
  status = 'OPEN'
  and framework_category_id = 'wct-id-422'
group by
  severity;
```

### List all open issues created in last 30 days
Explore which issues remain unresolved within the past month. This can help prioritize and manage ongoing tasks effectively.

```sql+postgres
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  status = 'OPEN'
  and created_at >= (current_timestamp - interval '30 days');
```

```sql+sqlite
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  status = 'OPEN'
  and created_at >= datetime('now','-30 days');
```

### Get the project information that the issue is related to
Explore which projects are associated with specific issues, including their status and severity, to understand the overall impact and urgency of each issue. This enables efficient project management and issue resolution.

```sql+postgres
select
  i.id,
  i.status,
  i.severity,
  i.created_at,
  p.name as project
from
  wiz_issue as i,
  jsonb_array_elements(i.projects) as pr
  left join wiz_project as p on p.id = pr ->> 'id';
```

```sql+sqlite
select
  i.id,
  i.status,
  i.severity,
  i.created_at,
  p.name as project
from
  wiz_issue as i,
  json_each(i.projects) as pr
  left join wiz_project as p on p.id = json_extract(pr.value, '$.id');
```

### List all high-severity issues open for more than a week
Explore which high-severity issues have remained unresolved for more than a week. This is useful in prioritizing and addressing critical problems that have been open for an extended period.

```sql+postgres
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'HIGH'
  and status = 'OPEN'
  and created_at < (current_timestamp - interval '7 days');
```

```sql+sqlite
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'HIGH'
  and status = 'OPEN'
  and created_at < datetime('now', '-7 days');
```

### Count open issues by service tag
Group open issues by the `Service` tag of the affected resource, to see which services carry the most risk.

```sql+postgres
select
  entity_tags ->> 'Service' as service,
  count(*) filter (where severity = 'CRITICAL') as critical,
  count(*) filter (where severity = 'HIGH') as high,
  count(*) as total
from
  wiz_issue
where
  status = 'OPEN'
group by
  service
order by
  critical desc,
  high desc;
```

```sql+sqlite
select
  json_extract(entity_tags, '$.Service') as service,
  sum(case when severity = 'CRITICAL' then 1 else 0 end) as critical,
  sum(case when severity = 'HIGH' then 1 else 0 end) as high,
  count(*) as total
from
  wiz_issue
where
  status = 'OPEN'
group by
  service
order by
  critical desc,
  high desc;
```

### List open critical issues without a service ticket
Find critical issues that nobody is tracking yet, with the affected resource.

```sql+postgres
select
  id,
  type,
  entity_type,
  entity_name,
  entity_subscription_name,
  url
from
  wiz_issue
where
  status = 'OPEN'
  and severity = 'CRITICAL'
  and has_service_ticket = false;
```

```sql+sqlite
select
  id,
  type,
  entity_type,
  entity_name,
  entity_subscription_name,
  url
from
  wiz_issue
where
  status = 'OPEN'
  and severity = 'CRITICAL'
  and has_service_ticket = 0;
```
