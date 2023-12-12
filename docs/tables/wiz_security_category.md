---
title: "Steampipe Table: wiz_security_category - Query Wiz Security Categories using SQL"
description: "Allows users to query Wiz Security Categories, specifically providing insights into security risk categories and potential vulnerabilities."
---

# Table: wiz_security_category - Query Wiz Security Categories using SQL

A Wiz Security Category is a classification within Wiz's cloud security platform that helps in identifying and managing potential security risks. It provides a way to categorize and prioritize security vulnerabilities, facilitating effective risk management. Wiz Security Categories help users stay informed about the security health of their cloud resources and take appropriate actions when potential vulnerabilities are discovered.

## Table Usage Guide

The `wiz_security_category` table provides insights into security risk categories within the Wiz cloud security platform. As a Security Engineer, explore category-specific details through this table, including the category's risk level, the associated vulnerabilities, and related metadata. Utilize it to uncover information about security risks, such as those with high-risk levels, the vulnerabilities associated with each category, and the verification of security protocols.

## Examples

### Basic info
Explore the different security categories within the Wiz framework to gain insights into their names, IDs, and descriptions. This could be useful for understanding the various security categories and their respective details, which can aid in improving overall security management.

```sql+postgres
select
  name,
  id,
  framework_id,
  description
from
  wiz_security_category;
```

```sql+sqlite
select
  name,
  id,
  framework_id,
  description
from
  wiz_security_category;
```

### Get the count of categories per framework
Discover the segments that have the highest number of categories within each framework. This can help prioritize which frameworks to focus on for security enhancements or audits.

```sql+postgres
select
  f.name,
  count(c.id) as category_count
from
  wiz_security_category as c
  join wiz_security_framework as f on f.id = c.framework_id
group by
  f.name;
```

```sql+sqlite
select
  f.name,
  count(c.id) as category_count
from
  wiz_security_category as c
  join wiz_security_framework as f on f.id = c.framework_id
group by
  f.name;
```

### List all open issues related to data security
Discover the segments that have open issues related to data security, categorized by severity. This information can help prioritize security efforts based on the severity of the issues.

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

### List all open issues related to vulnerability assessment
Gain insights into the number of open issues for each severity level related to vulnerability assessment. This can be used to prioritize security efforts based on the severity of the open issues.

```sql+postgres
select
  severity,
  count(id)
from
  wiz_issue
where
  status = 'OPEN'
  and framework_category_id = 'wct-id-3'
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
  and framework_category_id = 'wct-id-3'
group by
  severity;
```

### List all open issues related to cloud entitlements
Explore the severity level of open issues associated with cloud entitlements. This can be useful in prioritizing responses and allocating resources efficiently.

```sql+postgres
select
  severity,
  count(id)
from
  wiz_issue
where
  status = 'OPEN'
  and framework_category_id = 'wct-id-6'
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
  and framework_category_id = 'wct-id-6'
group by
  severity;
```