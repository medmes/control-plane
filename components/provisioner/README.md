# Provisioner

## Overview

Runtime Provisioner is a Kyma Control Plane component responsible for provisioning, installing, and deprovisioning clusters. When provisioning a cluster, you have an option to provision a cluster with Kyma (Kyma Runtime), or without it. To provision a Kyma Runtime, provide its configuration as **kymaConfig**.

For more details, see the Runtime Provisioner [documentation](https://github.com/kyma-project/control-plane/tree/main/docs/provisioner).

## Prerequisites

Before you can run Runtime Provisioner, you have to configure access to the PostgreSQL database. For development purposes, you can run a PostgreSQL instance in the Docker container executing the following command:

```bash
$ docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=password postgres
```

Runtime Provisioner also needs access to the cluster from which it fetches Secrets.  

## Development

### GraphQL schema

After you introduce changes in the GraphQL schema, run the `gqlgen.sh` script.

### Database schema

For tests to run properly, update the database schema in `./assets/database/provisioner.sql`. Provide the new migration in the Schema Migrator component in `resources/kcp/charts/provisioner/migrations`.

### Run Provisioner

To run Runtime Provisioner, use the following command:
```bash
go run cmd/main.go
```

### Environment Variables

This table lists the environment variables, their descriptions, and default values:


| Parameter | Description | Default value |
|-----------|-------------|---------------|
| **APP_ADDRESS** | Runtime Provisioner's address with the port | `127.0.0.1:3000` |
| **APP_METRICS_ADDRESS** | Runtime Provisioner Metrics' address with the port | `127.0.0.1:9000` |
| **APP_API_ENDPOINT** | Endpoint for the GraphQL API | `/graphql` |
| **APP_PLAYGROUND_API_ENDPOINT** | Endpoint for the API playground | `/graphql` |
| **APP_DIRECTOR_URL** | Director URL | `https://compass-gateway-auth-oauth.kyma.local/director/graphql` |
| **APP_SKIP_DIRECTOR_CERT_VERIFICATION** | Flag to skip certificate verification for Director | `false` |
| **APP_OAUTH_CREDENTIALS_NAMESPACE** | Namespace where the Director credentials are stored | `kcp-system` |
| **APP_OAUTH_CREDENTIALS_SECRET_NAME** | Runtime Provisioner credentials | `kcp-provisioner-credentials` |
| **APP_DATABASE_USER** | Database username | `postgres` |
| **APP_DATABASE_PASSWORD** | Database user password | `password` |
| **APP_DATABASE_HOST** | Database host | `localhost` |
| **APP_DATABASE_PORT** | Database port | `5432` |
| **APP_DATABASE_NAME** | Database name | `provisioner` |
| **APP_DATABASE_SSLMODE** | SSL Mode for PostgrSQL. See [all the possible values](https://www.postgresql.org/docs/9.1/libpq-ssl.html)  | `disable`|
| **APP_DATABASE_SSLROOTCERT** | Location of the PostgreSQL CA cert (Optional) | **optional** |
| **APP_PROVISIONING_TIMEOUT_INSTALLATION** | Kyma installation timeout | `60m`|
| **APP_PROVISIONING_TIMEOUT_UPGRADE** | Kyma installation timeout | `60m`|
| **APP_PROVISIONING_TIMEOUT_AGENT_CONFIGURATION** | Runtime Agent configuration timeout | `15m`|
| **APP_PROVISIONING_TIMEOUT_AGENT_CONNECTION** | Runtime Agent connection timeout | `15m`|
| **APP_GARDENER_PROJECT** | Name of the Gardener project connected to the service account  | `gardenerProject`|
| **APP_GARDENER_KUBECONFIG_PATH** | Filepath for the Gardener kubeconfig  | `./dev/kubeconfig.yaml`|
| **APP_GARDENER_AUDIT_LOGS_POLICY_CONFIG_MAP** | Name of the Config Map containing the audit logs policy  | **optional** |
| **APP_GARDENER_AUDIT_LOGS_TENANT** | Tenant used for storing audit logs  | **optional** |
| **APP_ENQUEUE_IN_PROGRESS_OPERATIONS** | Specifies whether operations in the `InProgress` state should be enqueued on the application startup | `true`|
