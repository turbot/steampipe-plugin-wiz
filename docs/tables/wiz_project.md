# Table: wiz_project

The `wiz_project` table can be used to query information about the Wiz projects.

Projects helps to group your cloud resources according to their users and/or purposes.

## Examples

### Basic info

```sql
select
  name,
  business_unit,
  security_score as security_score_in_percentage,
  repository_count,
  cloud_account_count,
  kubernetes_cluster_count
from
  wiz_project;
```

### Get count of critical issues per project

```sql
with critical_issues as (
  select
    id,
    severity,
    p ->> 'id' as project
  from
    wiz_issue,
    jsonb_array_elements(projects) as p
  where
    severity = 'CRITICAL'
)
select
  p.name as project,
  count(c.id)
from
  wiz_project as p
  left join critical_issues as c on p.id = c.project
group by
  p.name;
```

### Get the owner details of each project

```sql
select
  p.name,
  p.slug,
  u.name as user_name,
  u.email as user_email
from
  wiz_project as p
  left join jsonb_array_elements(project_owners) as o on true
  left join wiz_user as u on u.id = o ->> 'id';
```

### List archived projects

```sql
select
  name,
  business_unit,
  security_score as security_score_in_percentage,
  repository_count,
  cloud_account_count,
  kubernetes_cluster_count
from
  wiz_project
where
  archived;
```
