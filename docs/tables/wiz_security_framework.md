# Table: wiz_security_framework

The `wiz_security_framework` table can be used to query information about all the security frameworks.

A security framework is used to to organize controls, issues and findings according to the preferred category system.

## Examples

### Basic info

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
