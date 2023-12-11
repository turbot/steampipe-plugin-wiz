---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/wiz.svg"
brand_color: "#2D51DA"
display_name: "Wiz"
short_name: "wiz"
description: "Steampipe plugin to query security controls, findings, vulnerabilities, and more from your Wiz subscription."
og_description: "Query Wiz with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/wiz-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Wiz + Steampipe

[Wiz](https://www.wiz.io) provides direct visibility, risk prioritization, and remediation guidance for development teams to address risks in their own infrastructure and applications so they can ship faster and more securely.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

List all critical issues:

```sql
select
  id,
  status,
  severity,
  created_at
from
  wiz_issue
where
  severity = 'CRITICAL';
```

```
+--------------------------------------+----------+----------+---------------------------+
| id                                   | status   | severity | created_at                |
+--------------------------------------+----------+----------+---------------------------+
| fff8bfc2-c2f2-42ef-bfbc-2f4321ba85fd | OPEN     | CRITICAL | 2022-10-06T18:37:35+05:30 |
| fff9b66f-bf5e-1234-b567-8afdded9a0b0 | RESOLVED | CRITICAL | 2022-11-02T21:25:08+05:30 |
| fff1a2f3-4b56-78ac-bf90-12a34da5f67d | OPEN     | CRITICAL | 2022-09-28T23:40:49+05:30 |
+--------------------------------------+----------+----------+---------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/wiz/tables)**

## Quick start

### Install

Download and install the latest Wiz plugin:

```bash
steampipe plugin install wiz
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                          |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Credentials | Wiz requires your application's **Client ID** and **Client Secret** to authenticate all requests. You can find this value on the `Settings > Service Accounts` page. Also, provide the [GraphQL endpoint](https://docs.wiz.io/wiz-docs/docs/using-the-wiz-api#the-graphql-endpoint). |
| Permissions | Assign `read:all` scope to your service account.                                                                                                                                                                                                                                     |
| Radius      | Each connection represents a single Wiz installation.                                                                                                                                                                                                                                |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/wiz.spc`)<br />2. Credentials specified in environment variables, e.g., `WIZ_AUTH_CLIENT_ID`, `WIZ_AUTH_CLIENT_SECRET` and `WIZ_URL`.                                                                 |

### Configuration

Installing the latest wiz plugin will create a config file (`~/.steampipe/config/wiz.spc`) with a single connection named `wiz`:

```hcl
connection "wiz" {
  plugin = "wiz"

  # `client_id` (required) - Application's Client ID.
  # You can find this value on https://app.wiz.io/settings/service-accounts page.
  # This can also be set via the `WIZ_AUTH_CLIENT_ID` environment variable.
  # client_id = "8rp38Z6yb2cOSTeaMpPIpepAt99eg3ry"

  # `client_secret` (required) - Application's Client Secret.
  # You can find this value on https://app.wiz.io/settings/service-accounts page.
  # This can also be set via the `WIZ_AUTH_CLIENT_SECRET` environment variable.
  # client_secret = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IndJUnZwVWpBTU93WHQ5ZG5CXzRrVCJ9"

  # `url` (required) - Wiz API endpoint. This varies for each Wiz deployment.
  # See https://docs.wiz.io/wiz-docs/docs/using-the-wiz-api#the-graphql-endpoint.
  # This can also be set via the `WIZ_URL` environment variable.
  # url = "https://api.us1.app.wiz.io/graphql"
}
```

Alternatively, you can also use the standard Wiz environment variables to obtain credentials **only if other arguments (`client_id`, `client_secret` and `url`) are not specified** in the connection:

```sh
export WIZ_AUTH_CLIENT_ID=8rp38Z6yb2cOSTeaMpPIpepAt99eg3ry
export WIZ_AUTH_CLIENT_SECRET=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IndJUnZwVWpBTU93WHQ5ZG5CXzRrVCJ9
export WIZ_URL=https://api.us1.app.wiz.io/graphql
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-wiz
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
