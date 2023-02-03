# Table: wiz_cloud_configuration_finding

The `wiz_cloud_configuration_finding` table can be used to query information about all the findings.

A finding is generated when a configuration check, cloud configuration rule is applies to a specific cloud resource type.

**Note**: The table can limit a large dataset based on the number of rules and the number of cloud accounts where the rule is applied to; which can increase the query execution time. It is recommended that use the available filters to point out the specific set of data required. This table supports the below optional qualifiers that can be used:

- `result`
- `severity`
- `status`

## Examples

### Basic info

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding;
```

### List all failed resource with high severity

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and severity = 'HIGH';
```

### List failed resources which are not resolved

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and status = 'OPEN'
  and severity = 'HIGH';
```

### List all findings detected in last 3 days

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding
where
  result = 'FAIL'
  and status = 'OPEN'
  and severity = 'HIGH'
  and analyzed_at > (current_timestamp - interval '1 day');
```

### List failed resource with rule information

```sql
select
  title,
  result,
  severity,
  analyzed_at,
  status
from
  wiz_cloud_configuration_finding as f
  left join wiz_cloud_config_rule as r on
    f.rule ->> 'id' = r.id
    and f.result = 'FAIL'
    and f.status = 'OPEN'
    and f.severity = 'HIGH';
```
