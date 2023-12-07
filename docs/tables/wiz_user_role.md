---
title: "Steampipe Table: wiz_user_role - Query Wiz User Roles using SQL"
description: "Allows users to query User Roles in Wiz, specifically the permissions and characteristics of each role, providing insights into user access controls and potential security risks."
---

# Table: wiz_user_role - Query Wiz User Roles using SQL

The Wiz User Role is a resource in Wiz that defines the permissions and access rights of a user within the platform. It outlines the actions a user can perform, the resources they can access, and the level of access they have. This resource is crucial for managing user access, ensuring security, and maintaining compliance within the Wiz platform.

## Table Usage Guide

The `wiz_user_role` table provides insights into user roles within Wiz. As a security or IT professional, explore role-specific details through this table, including permissions, access levels, and associated metadata. Utilize it to uncover information about roles, such as those with elevated permissions, the relationships between users and roles, and the verification of access controls.

## Examples

### Basic info
Explore the roles within your user base, including their scope and descriptions, to gain insights into user permissions and responsibilities. This can help in managing access control and understanding the distribution of roles in your system.

```sql+postgres
select
  name,
  id,
  description,
  is_project_scoped,
  scopes
from
  wiz_user_role;
```

```sql+sqlite
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
Explore which user roles are specifically scoped to a project. This is useful in assessing the allocation of responsibilities and permissions within a project context.

```sql+postgres
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

```sql+sqlite
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
Explore which roles have the most users associated with them, allowing for an understanding of user distribution across different roles within the system. This can be beneficial for resource allocation and system management.

```sql+postgres
select
  r.name,
  count(u.id) as user_count
from
  wiz_user_role as r
  left join wiz_user as u on r.id = u.role ->> 'id'
group by
  r.name;
```

```sql+sqlite
select
  r.name,
  count(u.id) as user_count
from
  wiz_user_role as r
  left join wiz_user as u on r.id = json_extract(u.role, '$.id')
group by
  r.name;
```

### List users assigned with Global Admin role
Explore which users have been assigned the Global Admin role. This is useful for managing user permissions and ensuring only authorized individuals have access to sensitive data or administrative functions.

```sql+postgres
select
  u.name,
  u.email,
  r.name as role_name
from
  wiz_user as u
  join wiz_user_role as r on u.role ->> 'id' = r.id and r.id = 'GLOBAL_ADMIN';
```

```sql+sqlite
select
  u.name,
  u.email,
  r.name as role_name
from
  wiz_user as u
  join wiz_user_role as r on json_extract(u.role, '$.id') = r.id and r.id = 'GLOBAL_ADMIN';
```

### List admin roles
Explore which user roles have administrative privileges. This could be useful in auditing user permissions and ensuring appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  scopes as description,
  is_project_scoped
from
  wiz_user_role
where
  id like '%_ADMIN';
```