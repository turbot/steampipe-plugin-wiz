# Table: wiz_control

The `wiz_control` table can be used to query information about all controls in Wiz.

A Control consists of a pre-defined Security Graph query and a severity level — if a Control's query returns any results, an issue is generated for every result.

## Examples

### Basic info

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
