# envoy-cli

A lightweight CLI for managing environment variable sets across multiple deployment targets.

---

## Installation

```bash
go install github.com/your-username/envoy-cli@latest
```

Or build from source:

```bash
git clone https://github.com/your-username/envoy-cli.git
cd envoy-cli && go build -o envoy-cli .
```

---

## Usage

```bash
# Initialize a new environment config
envoy-cli init

# Add a variable to a target environment
envoy-cli set --target production DATABASE_URL=postgres://...

# List all variables for a target
envoy-cli list --target production

# Apply variables to a deployment target
envoy-cli apply --target staging
```

Environment sets are stored in a local `.envoy.yaml` file, making them easy to version control and share across teams.

---

## Configuration

```yaml
targets:
  production:
    DATABASE_URL: postgres://prod-host/db
    APP_ENV: production
  staging:
    DATABASE_URL: postgres://staging-host/db
    APP_ENV: staging
```

---

## Requirements

- Go 1.21+

---

## License

MIT © [your-username](https://github.com/your-username)