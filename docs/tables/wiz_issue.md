# Table: wiz_issue

The `wiz_issue` table can be used to query information about all issues created in Wiz.

## Examples

### Basic info

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

### List high severity issues which are open

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

### Get the project information the issue is related with

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

### List all open high-severity issues with age more than 1 week

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
