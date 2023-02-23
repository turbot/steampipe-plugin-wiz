# Table: wiz_user_role

The `wiz_user_role` table can be used to query information about user roles which can be defined for portal users as well as for federated SAML-based sessions in order to control the information they see and the actions they can perform.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  is_project_scoped,
  scopes
from
  wiz_user_role;
```

### List roles scoped to a specific project

```sql
select
  name,
  id,
  description,
  is_project_scoped,
  scopes
from
  wiz_user_role
where
  is_project_scoped;
```

### Count users per role

```sql
select
  r.name,
  count(u.id) as user_count
from
  wiz_user_role as r
  left join wiz_user as u on r.id = u.role ->> 'id'
group by
  r.name;
```

### List users assigned with Global Admin role

```sql
select
  u.name,
  u.email,
  r.name as role_name
from
  wiz_user as u
  join wiz_user_role as r on u.role ->> 'id' = r.id and r.id = 'GLOBAL_ADMIN';
```

### List admin roles

```sql
select
  name,
  id,
  jsonb_pretty(scopes)
  description,
  is_project_scoped
from
  wiz_user_role
where
  id like '%_ADMIN';
```
