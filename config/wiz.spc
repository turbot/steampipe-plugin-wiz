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