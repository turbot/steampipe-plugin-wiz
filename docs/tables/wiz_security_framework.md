---
title: "Steampipe Table: wiz_security_framework - Query Wiz Security Frameworks using SQL"
description: "Allows users to query Wiz Security Frameworks, providing detailed insights into the security status and configurations of cloud environments."
---

# Table: wiz_security_framework - Query Wiz Security Frameworks using SQL

Wiz Security Framework is a tool within the Wiz platform that provides a comprehensive view of the security status and configurations of cloud environments. It allows users to monitor and manage security risks across their cloud infrastructure, including virtual machines, databases, web applications, and more. Wiz Security Framework helps you stay informed about the security health of your cloud resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `wiz_security_framework` table provides insights into the security status and configurations of cloud environments. As a security engineer, explore the details of your cloud infrastructure's security through this table, including risk levels, associated metadata, and more. Utilize it to uncover information about potential security risks, the status of security configurations, and to verify the effectiveness of current security measures.

## Examples

### Basic info
Explore which security frameworks are enabled and built-in, and gain insights into their descriptions. This can help in assessing the security setup and identifying areas for potential improvement.

```sql
select
  name,
  id,
  enabled,
  built_in,
  description
from
  wiz_security_framework;
```

### List disabled security frameworks  
Identify instances where certain security frameworks are disabled. This can be useful for assessing the areas in your system that may lack necessary protection.

```sql
select
  name,
  id,
  enabled,
  built_in,
  description
from
  wiz_security_framework
where
  not enabled;
```

### List built-in frameworks
Explore which built-in security frameworks are currently enabled. This can help you understand what default security measures are in place and assist in identifying potential areas for improvement.

```sql
select
  name,
  id,
  enabled,
  built_in,
  description
from
  wiz_security_framework
where
  built_in;
```

### Get the count of security categories per framework
Explore the number of security categories within each framework to better understand the security measures in place and to identify areas for potential improvement. This could be particularly useful for IT teams looking to enhance their organization's security posture.

```sql
select
  name,
  built_in,
  enabled,
  jsonb_array_length(categories) as categories
from
  wiz_security_framework;
```

### Get the count of controls per framework
Explore the distribution of controls across different security frameworks to understand which frameworks have the most controls. This can be useful for prioritizing which frameworks to implement based on their comprehensiveness.

```sql
select
  f.name,
  f.built_in,
  f.enabled,
  count(c.id) as control_count
from
  wiz_control as c
  join wiz_security_framework as f on c.framework_category_id = f.id
group by
  1, 2, 3;
```