---
title: Init Command
description: Initialize project configuration files (phi.yml / phi.json) using the Phi CLI.
category: CLI
tags: [cli, init, configuration, setup, phi-yml, phi-json]
---

# Init Command

The `phi init` command bootstraps a configuration file for your project. The configuration file dictates schema locations, output directories for generated Go client packages, migration directories, client package naming, embedded migration flags, and logging options.

---

## 1. Quick Usage

Generate a default YAML configuration file (`phi.yml`) in the current directory:

```bash
phi init
```

Specify a target output directory:

```bash
phi init ./my-project
```

---

## 2. Supported Formats

Phi supports both **YAML** (`.yml` / `.yaml`) and **JSON** (`.json`) configuration formats.

### YAML (Default)

```bash
phi init yml [directory]
```

Generates `phi.yml`:

```yaml
# Name of the generated Go client package (default: "phi")
client_name: phi

# Enable Go 1.16+ embedded migrations (//go:embed) inside client
embed_migrations: true

database:
  url_env: DATABASE_URL              # Environment variable for runtime client connections.
  direct_url_env: DATABASE_DIRECT_URL # Environment variable used for DDL migrations.

schema: ./schema.prisma              # Path to your Prisma schema.

output:
  client: ./phi                      # Output directory for generated Go client code.
  migrations: ./phi/migrations       # Directory for Goose-compatible .sql migration files.

log:
  - none                             # Logging flags: none, query, warn, error, info, all
```

### JSON

```bash
phi init json [directory]
```

Generates `phi.json`:

```json
{
  "client_name": "phi",
  "embed_migrations": true,
  "database": {
    "url_env": "DATABASE_URL",
    "direct_url_env": "DATABASE_DIRECT_URL"
  },
  "schema": "./schema.prisma",
  "output": {
    "client": "./phi",
    "migrations": "./phi/migrations"
  },
  "log": [
    "none"
  ]
}
```

---

## 3. Configuration Fields Breakdown

| Field | Description | Default |
| :--- | :--- | :--- |
| `client_name` | Go package name for the generated client code. | `"phi"` |
| `embed_migrations` | Controls whether Go 1.16+ `//go:embed` migration code is generated inside the client. | `true` |
| `database.url_env` | Environment variable holding runtime database connection URL. | `"DATABASE_URL"` |
| `database.direct_url_env` | Environment variable holding migration/direct DDL connection URL. | `"DATABASE_DIRECT_URL"` |
| `schema` | Path to your input `schema.prisma` file. | `"./schema.prisma"` |
| `output.client` | Directory path where generated client code will be saved. | `"./phi"` |
| `output.migrations` | Directory path where versioned `.sql` migration files are stored. | `"./phi/migrations"` |
| `log` | Active logging categories (`none`, `query`, `warn`, `error`, `info`, `all`). | `["none"]` |

---

## 4. Configuration Ambiguity Protection

To prevent conflicting settings across different files, Phi inspects your project directory for configuration files (`phi.yml`, `phi.yaml`, `phi.json`).

If **more than one** configuration file is detected in the same directory (e.g. both `phi.yml` and `phi.json`), Phi CLI immediately halts with a fatal error:

```text
FATAL: multiple configuration files found (phi.yml, phi.json). Please keep only one configuration file in the project directory.
```
