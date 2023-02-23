# Table: wiz_issue

The `wiz_issue` table can be used to query information about all issues created in Wiz.

A **Control** consists of a pre-defined **Security Graph** query and a severity level — if a control's query returns any results, an issue is generated for every result.

**Note**: The table can return a large dataset; which can increase the query execution time. It is recommended that queries to this table should include (usually in the `where` clause) at least one of these columns:

- `control_id`
- `created_at`
- `framework_category_id`
- `resolution_reason`
- `severity`
- `status`

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

### List high severity open issues

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

### Get the project information that the issue is related to

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
