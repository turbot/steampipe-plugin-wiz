---
title: "Steampipe Table: wiz_subscription - Query Wiz Subscriptions using SQL"
description: "Allows users to query Wiz Subscriptions, specifically providing detailed insights into each subscription's attributes and status."
---

# Table: wiz_subscription - Query Wiz Subscriptions using SQL

Wiz is a cloud security platform that discovers all assets in cloud and container environments, prioritizes risks based on potential impact, and continuously fixes critical security issues. It provides a comprehensive view of the risks in your cloud environment across all cloud resources, including virtual machines, databases, web applications, and more. Wiz helps you stay informed about the security status of your cloud resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `wiz_subscription` table provides insights into subscriptions within Wiz. As a security engineer, explore subscription-specific details through this table, including subscription attributes, status, and associated metadata. Utilize it to uncover information about subscriptions, such as those with active or inactive status, the attributes of each subscription, and the verification of subscription details.

## Examples

### Basic info
Explore the status and last scanned timestamp of your cloud subscriptions across different providers. This can help you monitor and manage the health and security of your cloud infrastructure.

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
Discover the segments that are currently connected to your AWS cloud accounts. This is particularly useful for understanding which accounts are active and when they were last scanned, aiding in maintaining efficient account management and security.

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
Uncover the details of cloud accounts that are only partially connected. This is useful to identify potential issues with your cloud accounts, such as incomplete setup or connection problems, which may impact your ability to fully utilize cloud services.

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
Discover the cloud accounts that have not been scanned in the last 24 hours. This query is useful in identifying potential security risks by pinpointing accounts that may have been overlooked during routine checks.

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
Discover the segments that are associated with cloud accounts not linked to any project. This can be particularly useful for organizations looking to streamline their cloud resources or identify unused accounts.

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