# Table: wiz_service_account

The `wiz_service_account` table can be used to query information about the service accounts.

Service accounts are accounts created to serve as machine to machine interfaces that can authenticate with Wiz API.

## Examples

### Basic info

```sql
select
  name,
  client_id,
  type,
  created_at,
  last_rotated_at,
  authentication_source
from
  wiz_service_account;
```

### List service accounts scoped to a specific project

```sql
select
  s.name,
  p.name as project,
  s.client_id,
  s.created_at,
  s.last_rotated_at,
  s.scopes
from
  wiz_service_account as s,
  jsonb_array_elements(s.assigned_projects) as ap
  join wiz_project as p on p.id = ap ->> 'id';
```

### List service accounts by name

```sql
select
  name,
  client_id,
  created_at,
  last_rotated_at,
  scopes
from
  wiz_service_account
where
  name = 'Steampipe';
```

### List service accounts not rotated in last 90 days

```sql
select
  name,
  client_id,
  created_at,
  last_rotated_at,
  scopes
from
  wiz_service_account
where
  last_rotated_at < (current_timestamp - interval '90 days');
```

### List users scoped to a specific project

```sql
select
  u.name,
  u.email,
  u.identity_provider_type,
  u.role ->> 'name' as role,
  p.name as project,
  u.created_at
from
  wiz_user as u,
  jsonb_array_elements(u.effective_assigned_projects) as ep
  join wiz_project as p on ep ->> 'id' = p.id;
```
