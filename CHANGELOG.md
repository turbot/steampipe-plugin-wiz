## Unreleased

_What's new?_

- New tables added
  - [wiz_iac_finding](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_iac_finding)
  - [wiz_sast_finding](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_sast_finding)
  - [wiz_secret_finding](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_secret_finding)
  - [wiz_version_control_repository](https://hub.steampipe.io/plugins/turbot/wiz/tables/wiz_version_control_repository)
- Added code (software composition analysis) and risk columns to the `wiz_vulnerability_finding` table, including `asset_type`, `repository_id`, `repository_name`, `reachability`, `reachability_evidence`, `call_paths`, `transitivity`, `is_transitive`, `has_fix`, `wiz_severity`, `vendor_severity`, `nvd_severity`, `weighted_severity`, `epss_*`, `origin`, `projects`, `vcs_code_owners`, `related_issue_analytics`, `artifact_type`, `updated_at` and `status_updated_at`, with new optional quals `asset_type`, `code_library_language`, `has_fix`, `is_default_branch`, `is_malicious_package`, `origin`, `reachability`, `transitivity`, `updated_at` and `wiz_severity`.
- Added `service_tickets` and `has_service_ticket` columns to the `wiz_vulnerability_finding`, `wiz_sast_finding` and `wiz_secret_finding` tables. They list the service tickets linked through the finding's related and posture issues, and are only fetched when selected.
- Added `source_rules`, `type`, `url`, `first_event_at`, `last_event_at`, `resolution_note`, `open_reason`, `has_service_ticket`, `has_code_remediation`, `has_code_to_cloud_mapping`, `validated_as_exploitable`, `environments`, `entity_snapshot`, `entity_*`, `entity_tags`, `kubernetes_cluster_name`, `kubernetes_namespace_name`, `application_services`, `related_application_services` and `risk_factors` columns to the `wiz_issue` table, with new optional quals `type`, `has_service_ticket` and `project_id`, and `in` lists on `severity`, `status` and `control_id`. Application services and risk factors are only fetched when selected.
- Added `tags_v2` column to the `wiz_control` table.
- Added `builtin`, `is_tenant_group_scoped`, `user_count`, `created_at` and `updated_at` columns to the `wiz_user_role` table, with optional quals `builtin` and `is_project_scoped`.
- Added `status` and `updated_at` columns to the `wiz_user` table.
- Added `*_system_health_issue_count` columns to the `wiz_subscription` table.

_Bug fixes_

- Fixed a panic when credentials are set only through the `WIZ_AUTH_CLIENT_ID`, `WIZ_AUTH_CLIENT_SECRET` and `WIZ_URL` environment variables.
- Fixed the `client_secret` presence check, which read `client_id` instead.
- Fixed a panic on network errors while requesting an access token.
- Fixed the `created_at` filter on the `wiz_issue` table and the `analyzed_at` filter on the `wiz_cloud_configuration_finding` table, which kept only one bound when both were given. Time filters now keep sub-second precision, and `in` lists on timestamps are filtered by Postgres instead of the API.

_Breaking changes_

- The `wiz_vulnerability_finding` table now returns findings in all statuses (`OPEN`, `RESOLVED`, `REJECTED`) when `status` is not filtered. Before, the Wiz API returned only open findings, so counts without a `status` filter now include resolved and rejected findings. The new `wiz_sast_finding`, `wiz_secret_finding` and `wiz_iac_finding` tables behave the same way.
- Findings statuses are `OPEN`, `RESOLVED` and `REJECTED`. The Wiz API treats `IN_PROGRESS` as `OPEN` for findings, so these tables never return it, and a filter on `IN_PROGRESS` alone returns no rows. `wiz_vulnerability_finding.status` no longer accepts the old values `UNRESOLVED`, `PASSED` and `IGNORED`. The `vulnerability_external_id` qual now uses `vulnerabilityExternalIdV2`.
- `wiz_user.identity_provider_type` now returns `SSO` instead of `SAML` (the Wiz API renamed the value). Queries filtering on `'SAML'` must use `'SSO'`.
- Moved queries off other fields Wiz has deprecated:
  - `wiz_issue` now uses `issuesV2` and `sourceRules`. The `control_id` qual filters by source rule; `control_id` returns the matching source rule, or else the first source rule of type control.
  - `wiz_control.tags` is now built from `tagsV2` as `key` or `key:value` strings.
  - `wiz_project.project_owners` now comes from `projectOwnersV2` (owner identities), and `resource_tag_links` from `resourceFilterLinks`.
  - `wiz_user.role` is now the effective role.
  - `wiz_user_role` now uses `userRolesV2`.

_Deprecations_

- `wiz_subscription.status` and `wiz_user.is_analytics_enabled`, `ip_address` and `tenant_id` are still served but deprecated by Wiz with no replacement, and may stop returning data in a future Wiz API release.

## v1.3.0 [2025-10-13]

_Dependencies_

- Recompiled plugin with Go version `1.24`. ([#64](https://github.com/turbot/steampipe-plugin-wiz/pull/64))
- Recompiled plugin with [steampipe-plugin-sdk v5.13.1](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5131-2025-09-25) that addresses critical and high vulnerabilities in dependent packages. ([#68](https://github.com/turbot/steampipe-plugin-wiz/pull/68))

## v1.2.0 [2025-07-17]

_Bug fixes_

- Fixed the `control` column of `wiz_cloud_config_rule` table to correctly return data instead of an error. ([#61](https://github.com/turbot/steampipe-plugin-wiz/pull/61))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.7](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5117-2025-06-04) which resolves an issue where rate limiters were not being applied to hydrate functions correctly. ([#58](https://github.com/turbot/steampipe-plugin-wiz/pull/58))

## v1.1.1 [2025-04-18]

_Bug fixes_

- Fixed Linux AMD64 plugin build failures for `Postgres 14 FDW`, `Postgres 15 FDW`, and `SQLite Extension` by upgrading GitHub Actions runners from `ubuntu-20.04` to `ubuntu-22.04`.

## v1.1.0 [2025-04-17]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#50](https://github.com/turbot/steampipe-plugin-wiz/pull/50))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#50](https://github.com/turbot/steampipe-plugin-wiz/pull/50))

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
