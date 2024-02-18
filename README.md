# CZERTAINLY HashiCorp Vault Connector



## Docker container

HashiCorp Vault `Connector` is provided as a Docker container. Use the `docker.io/3keycompany/czertainly-hashicorp-vaul-connector:tagname` to pull the required image from the repository. It can be configured using the following environment variables:

| Variable      | Description                       | Required                                           | Default value |
|---------------|-----------------------------------|----------------------------------------------------|---------------|
| `SERVER_PORT` | Port where the service is exposed | ![](https://img.shields.io/badge/-NO-red.svg)      | `8080`        |
| `DB_HOST`     | Database host                     | ![](https://img.shields.io/badge/-YES-success.svg) | `N/A`         |
| `DB_PORT`     | Database port                     | ![](https://img.shields.io/badge/-NO-red.svg)      | `5432`        |
| `DB_NAME`     | Database name                     | ![](https://img.shields.io/badge/-YES-success.svg) | `N/A`         |
| `DB_USERNAME` | Database user                     | ![](https://img.shields.io/badge/-YES-success.svg) | `N/A`         |
| `DB_PASSWORD` | Database password                 | ![](https://img.shields.io/badge/-YES-success.svg) | `N/A`         |
| `LOG_LEVEL`   | Logging level for the service     | ![](https://img.shields.io/badge/-NO-red.svg)      | `INFO`        |