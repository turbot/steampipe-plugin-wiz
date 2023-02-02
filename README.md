![image](https://hub.steampipe.io/images/plugins/turbot/wiz-social-graphic.png)

# Wiz Plugin for Steampipe

Use SQL to query security controls, cloud configuration rules, security frameworks, and more from your Wiz subscription.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/wiz)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/wiz/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-wiz/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install wiz
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/wiz#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/wiz#configuration).

Run a query:

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
