# Table: wiz_cloud_config_rule

The `wiz_cloud_config_rule` table can be used to query information about all the cloud configuration rule created.

A **Cloud Configuration Rule** is a configuration check that applies to a specific cloud resource type — if a resource does not pass a **Rule**, a **Configuration Finding** is generated and associated with the resource on the **Security Graph**.

## Examples

### Basic info

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule;
```

### List disabled rules

```sql
select
  name,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  not enabled;
```

### List built-in rules

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  built_in;
```

### List high-severity rules specific to AWS S3 bucket

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  cloud_provider = 'AWS'
  and severity = 'HIGH'
  and target_native_types ?| array['bucket'];
```

### List all findings of a rule specific to AWS S3 bucket in last 3 days

```sql
with list_s3_bucket_rules as (
  select
    name,
    id
  from
    wiz_cloud_config_rule
  where
    cloud_provider = 'AWS'
    and severity = 'HIGH'
    and target_native_types ?| array['bucket']
)
select
  r.name as rule_name,
  f.resource ->> 'Name' as resource_name,
  f.result as finding_status,
  f.analyzed_at
from
  wiz_cloud_configuration_finding as f
  join list_s3_bucket_rules as r on
    f.rule_id = r.id
    and f.severity = 'HIGH'
    and f.analyzed_at > (current_timestamp - interval '3 day');
```

### List rules with auto remediation enabled

```sql
select
  name,
  enabled,
  severity,
  cloud_provider,
  target_native_types,
  created_at
from
  wiz_cloud_config_rule
where
  has_auto_remediation;
```
