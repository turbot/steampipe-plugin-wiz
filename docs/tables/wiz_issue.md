---
title: "Steampipe Table: wiz_issue - Query Wiz Issues using SQL"
description: "Allows users to query Wiz Issues, providing detailed information about security risks and vulnerabilities identified in their environment."
---

# Table: wiz_issue - Query Wiz Issues using SQL

Wiz is a Cloud Security Posture Management (CSPM) tool that provides continuous security posture monitoring for cloud environments. It identifies security risks and vulnerabilities across a wide range of categories, including misconfigurations, policy violations, and threats. Wiz provides a holistic view of your security posture, enabling you to identify and remediate issues quickly and effectively.

## Table Usage Guide

The `wiz_issue` table provides insights into the security issues identified by Wiz in your cloud environment. As a security engineer, you can use this table to explore detailed information about each issue, including its severity, status, and associated resources. This can help you prioritize remediation efforts and improve your overall security posture.

## Examples

### Basic info
Explore which issues have been logged in your system, their severity and status, and when they were created. This can help you understand the range and depth of problems encountered, and the reasons provided for their resolution.

```sql
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

```sql
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

```sql
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

```sql
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

```sql
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

### Get the project information that the issue is related to
Explore which projects are associated with specific issues, including their status and severity, to understand the overall impact and urgency of each issue. This enables efficient project management and issue resolution.

```sql
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

### List all high-severity issues open for more than a week
Explore which high-severity issues have remained unresolved for more than a week. This is useful in prioritizing and addressing critical problems that have been open for an extended period.

```sql
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