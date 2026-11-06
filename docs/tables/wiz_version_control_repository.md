---
title: "Steampipe Table: wiz_version_control_repository - Query Wiz Version Control Repositories using SQL"
description: "Allows users to query the version control repositories connected to Wiz, providing details such as platform, visibility, default branch and project assignment."
---

# Table: wiz_version_control_repository - Query Wiz Version Control Repositories using SQL

Wiz connects to version control platforms such as GitHub, GitLab, Azure DevOps and Bitbucket to scan code repositories. This table lists the repositories Wiz knows about, along with their platform, visibility, default branch, owning organization and the Wiz projects they belong to.

## Table Usage Guide

The `wiz_version_control_repository` table provides an inventory of the code repositories connected to Wiz. As a security engineer, you can use it to check scan coverage, find public or archived repositories, and join repositories with code findings from the `wiz_sast_finding`, `wiz_secret_finding`, `wiz_iac_finding` and `wiz_vulnerability_finding` tables.

## Examples

### Basic info
List the repositories connected to Wiz.

```sql+postgres
select
  name,
  platform,
  visibility,
  default_branch,
  archived,
  url
from
  wiz_version_control_repository;
```

```sql+sqlite
select
  name,
  platform,
  visibility,
  default_branch,
  archived,
  url
from
  wiz_version_control_repository;
```

### List public repositories
Review the repositories that are publicly visible.

```sql+postgres
select
  name,
  organization_name,
  url
from
  wiz_version_control_repository
where
  public;
```

```sql+sqlite
select
  name,
  organization_name,
  url
from
  wiz_version_control_repository
where
  public = 1;
```

### List repositories not assigned to any project
Find repositories that are missing project ownership in Wiz.

```sql+postgres
select
  name,
  platform,
  url
from
  wiz_version_control_repository
where
  projects is null
  or jsonb_array_length(projects) = 0;
```

```sql+sqlite
select
  name,
  platform,
  url
from
  wiz_version_control_repository
where
  projects is null
  or json_array_length(projects) = 0;
```

### Count open SAST findings per repository
Join the repository inventory with the SAST findings.

```sql+postgres
select
  r.name,
  count(f.id) as open_sast_findings
from
  wiz_version_control_repository as r
  left join wiz_sast_finding as f on f.repository_id = r.id
  and f.status = 'OPEN'
where
  not r.archived
group by
  r.name
order by
  open_sast_findings desc;
```

```sql+sqlite
select
  r.name,
  count(f.id) as open_sast_findings
from
  wiz_version_control_repository as r
  left join wiz_sast_finding as f on f.repository_id = r.id
  and f.status = 'OPEN'
where
  r.archived = 0
group by
  r.name
order by
  open_sast_findings desc;
```
