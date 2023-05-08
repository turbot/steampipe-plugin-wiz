# Table: wiz_security_category

The `wiz_security_category` table can be used to query information about all the security categories and subcategories.

A security category is used to to organize controls, issues and findings according to the preferred category system.

## Examples

### Basic info

```sql
select
  name,
  id,
  framework_id,
  description
from
  wiz_security_category;
```

### Get the count of categories per framework

```sql
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

### List all open issues related to vulnerability assessment

```sql
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

```sql
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
