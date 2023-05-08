![image](https://hub.steampipe.io/images/plugins/turbot/wiz-social-graphic.png)

# Wiz Plugin for Steampipe

Use SQL to query security controls, findings, vulnerabilities, and more from your Wiz subscription.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/wiz)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/wiz/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
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

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-wiz/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Wiz Plugin](https://github.com/turbot/steampipe-plugin-wiz/labels/help%20wanted)
