## v1.0.1 [2025-02-05]

_Bug fixes_

- Fixed pagination in `wiz_user` table. ([#41](https://github.com/turbot/steampipe-plugin-wiz/pull/41))
- Removed the unsupported `profile_completion` column from the `wiz_project` table. ([#47](https://github.com/turbot/steampipe-plugin-wiz/pull/47))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. 
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. 

## v0.3.0 [2024-02-15]

_Enhancements_

- Improved the plugin error message when invalid credentials are set in the `wiz.spc` file. ([#23](https://github.com/turbot/steampipe-plugin-wiz/pull/23))

_Bug fixes_

- Fixed the `service_tickets` column in `wiz_issue` table by removing the `action` subfield from the `ServiceTickets` field in the GraphQL response since it was no longer available. ([#24](https://github.com/turbot/steampipe-plugin-wiz/pull/24) [#25](https://github.com/turbot/steampipe-plugin-wiz/pull/25)) (Thanks [@sycophantic](https://github.com/sycophantic) for the contribution!)

## v0.2.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#18](https://github.com/turbot/steampipe-plugin-wiz/pull/18))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#18](https://github.com/turbot/steampipe-plugin-wiz/pull/18))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-wiz/blob/main/docs/LICENSE). ([#18](https://github.com/turbot/steampipe-plugin-wiz/pull/18))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#17](https://github.com/turbot/steampipe-plugin-wiz/pull/17))

## v0.1.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#9](https://github.com/turbot/steampipe-plugin-wiz/pull/9))

## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#7](https://github.com/turbot/steampipe-plugin-wiz/pull/7))
- Recompiled plugin with Go version `1.21`. ([#7](https://github.com/turbot/steampipe-plugin-wiz/pull/7))

## v0.0.1 [2023-05-08]

_What's new?_

- New tables added
  - [wiz_cloud_config_rule](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_cloud_config_rule)
  - [wiz_cloud_configuration_finding](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_cloud_configuration_finding)
  - [wiz_control](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_control)
  - [wiz_issue](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_issue)
  - [wiz_project](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_project)
  - [wiz_security_category](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_security_category)
  - [wiz_security_framework](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_security_framework)
  - [wiz_service_account](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_service_account)
  - [wiz_subscription](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_subscription)
  - [wiz_user_role](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_user_role)
  - [wiz_user](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_user)
  - [wiz_vulnerability_finding](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_vulnerability_finding)
