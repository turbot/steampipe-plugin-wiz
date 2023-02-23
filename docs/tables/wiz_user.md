# Table: wiz_user

The `wiz_user` table can be used to query information about all users authenticated to Wiz.

## Examples

### Basic info

```sql
select
  name,
  email,
  identity_provider_type,
  role ->> 'name' as role,
  is_suspended,
  last_login_at,
  created_at
from
  wiz_user;
```

### List suspended users

```sql
select
  name,
  email,
  identity_provider_type,
  role ->> 'name' as role,
  last_login_at,
  created_at
from
  wiz_user
where
  is_suspended;
```

### List inactive users

```sql
select
  name,
  email,
  identity_provider_type,
  role ->> 'name' as role,
  is_suspended,
  created_at
from
  wiz_user
where
  last_login_at is null;
```

### List all administrators

```sql
select
  name,
  email,
  identity_provider_type,
  role ->> 'name' as role,
  created_at
from
  wiz_user
where
  role ->> 'id' = 'GLOBAL_ADMIN'
  and not (role ->> 'isProjectScoped')::boolean;
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

### List all SAML users

```sql
select
  name,
  email,
  identity_provider_type,
  role ->> 'name' as role,
  is_suspended,
  last_login_at,
  created_at
from
  wiz_user
where
  identity_provider_type = 'SAML';
```
