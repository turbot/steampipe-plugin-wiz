# Table: wiz_subscription

The `wiz_subscription` table can be used to query information about all the cloud accounts connected to Wiz.

Cloud accounts are detected and provisioned automatically when you connect Wiz on the management group level (in Azure) or organization level (in GCP, AWS, and Alibaba Cloud), or Tenancy level (in OCI).

## Examples

### Basic info

```sql
select
  name,
  cloud_provider,
  status,
  last_scanned_at
from
  wiz_subscription;
```

### List all connected AWS cloud accounts

```sql
select
  name,
  cloud_provider,
  status,
  last_scanned_at
from
  wiz_subscription
where
  cloud_provider = 'AWS';
```

### List partially connected cloud accounts

```sql
select
  name,
  cloud_provider,
  status,
  last_scanned_at
from
  wiz_subscription
where
  status = 'PARTIALLY_CONNECTED';
```

### List cloud accounts not checked in last 24 hours

```sql
select
  name,
  cloud_provider,
  status,
  last_scanned_at
from
  wiz_subscription
where
  last_scanned_at < (current_timestamp - interval '1 day');
```

### List cloud accounts not linked to any project

```sql
select
  name,
  cloud_provider,
  status,
  last_scanned_at
from
  wiz_subscription
where
  linked_projects is null;
```
