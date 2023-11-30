---
title: "Steampipe Table: wiz_control - Query Wiz Controls using SQL"
description: "Allows users to query Wiz Controls, specifically the status and details of security risks, providing insights into potential vulnerabilities and security threats."
---

# Table: wiz_control - Query Wiz Controls using SQL

Wiz is a cloud security platform that provides visibility and threat detection for cloud infrastructure. It offers comprehensive coverage across the full stack of multi-cloud environments, identifying security risks in the most critical areas of vulnerability, such as identity and access, network and firewall, data and storage, and cloud native services. Wiz Controls are the specific security risks that the platform identifies and monitors, providing detailed insights into potential vulnerabilities and security threats.

## Table Usage Guide

The `wiz_control` table provides insights into specific security risks within the Wiz cloud security platform. As a security analyst, explore control-specific details through this table, including risk severity, status, and associated metadata. Utilize it to uncover information about controls, such as those with high risk severity, the status of these risks, and the verification of mitigation efforts.

## Examples

### Basic info
Explore which security controls are currently active, their severity, and when they were created to better understand your system's security posture. This can help you identify potential vulnerabilities and prioritize remediation efforts.

```sql
select
  name,
  id,
  severity,
  enabled,
  type,
  created_at
from
  wiz_control;
```

### List disabled controls
Explore which controls are currently disabled in your system. This query is useful for identifying potential system vulnerabilities and areas that may require attention or reconfiguration.

```sql
select
  name,
  id,
  severity,
  enabled,
  type,
  created_at
from
  wiz_control
where
  not enabled;
```

### List controls with high severity
Discover the segments that have high severity controls activated. This is useful in prioritizing and managing risks within your system.

```sql
select
  name,
  id,
  severity,
  enabled,
  type,
  created_at
from
  wiz_control
where
  enabled
  and severity = 'HIGH';
```

### Get the count of open issues per control
Explore the number of unresolved issues for each control mechanism. This can help prioritize which controls require immediate attention and action.

```sql
select
  c.name,
  count(i.id) as issue_count
from
  wiz_issue as i
  join wiz_control as c on
    i.control_id = c.id
    and i.status = 'OPEN'
    and c.enabled
group by
  c.name;
```

### Get all issues created by a specific control
This example helps to identify all open issues that have been created by a specific control within your system. It's particularly useful for understanding the severity and status of these issues, as well as when they were created, thus aiding in prioritizing and managing system vulnerabilities.

```sql
select
  c.name as control_name,
  i.entity ->> 'name' as resource,
  i.severity as issue_severity,
  i.status as issue_status,
  i.created_at
from
  wiz_issue as i
  join wiz_control as c on (
    c.id = 'wc-id-613'
    and i.control_id = c.id
    and i.status = 'OPEN'
    and c.enabled
  );
```