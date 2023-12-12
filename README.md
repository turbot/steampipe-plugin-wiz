![image](https://hub.steampipe.io/images/plugins/turbot/wiz-social-graphic.png)

# Wiz Plugin for Steampipe

Use SQL to query security controls, findings, vulnerabilities, and more from your Wiz subscription.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/wiz)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/wiz/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-wiz/issues)

## Quick start

### Install

Download and install the latest Wiz plugin:

```shell
steampipe plugin install wiz
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/wiz#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/wiz#configuration).

Configure your subscription details in `~/.steampipe/config/wiz.spc`:

```hcl
connection "wiz" {
  plugin = "wiz"
  # Authentication information
  client_id     = "8rp38Z6yb2cOSTeaMpPIpepAt99eg3ry"
  client_secret = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IndJUnZwVWpBTU93WHQ5ZG5CXzRrVCJ9"
  url           = "https://api.us1.app.wiz.io/graphql"
}
```

Or through environment variables:

```sh
export WIZ_AUTH_CLIENT_ID=8rp38Z6yb2cOSTeaMpPIpepAt99eg3ry
export WIZ_AUTH_CLIENT_SECRET=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IndJUnZwVWpBTU93WHQ5ZG5CXzRrVCJ9
export WIZ_URL=https://api.us1.app.wiz.io/graphql
```

Run steampipe:

```shell
steampipe query
```

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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/index) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/index) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/index) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-wiz.git
cd steampipe-plugin-wiz
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/wiz.spc
```

Try it!

```shell
steampipe query
> .inspect wiz
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Wiz Plugin](https://github.com/turbot/steampipe-plugin-wiz/labels/help%20wanted)
