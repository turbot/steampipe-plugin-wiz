---
title: "Steampipe Table: wiz_user - Query Wiz User using SQL"
description: "Allows users to query Wiz Users, providing insights into user details and roles within the Wiz environment."
---

# Table: wiz_user - Query Wiz User using SQL

Wiz is a cloud security platform that discovers all assets in your cloud environment to detect critical risks and security holes. It provides a holistic view of each resource, including metadata, configurations, network paths, and potential vulnerabilities. Wiz User represents the user accounts in the Wiz platform, each having specific roles and permissions.

## Table Usage Guide

The `wiz_user` table provides insights into user accounts within the Wiz platform. As a security analyst, you can explore user-specific details through this table, including roles, permissions, and associated metadata. Utilize it to uncover information about user accounts, such as those with admin permissions, the roles assigned to each user, and the verification of user statuses.

## Examples

### Basic info
Gain insights into user details such as their role and status (active or suspended), as well as their last login and account creation dates. This can be useful for understanding user behavior and managing user access.

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
Discover the segments that consist of suspended users to assess potential security risks or to manage user access, ensuring a safer and more controlled environment.

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
Discover the segments that include users who have not logged in yet. This can be useful in identifying inactive users for potential outreach or account clean-up efforts.

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
Discover the segments that have global administrators and their corresponding details, without including project-scoped roles. This helps in identifying and managing users with overarching control and access within your system.

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
Explore which users are assigned to a specific project. This can help in understanding the distribution of team members across projects, facilitating efficient resource management and task allocation.

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
Discover the segments that are using SAML as their identity provider. This can be particularly useful to understand user distribution across different identity providers and assess if any particular group needs attention.

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