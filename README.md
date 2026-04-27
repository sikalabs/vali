# vali

**vali** validates that required environment variables and files exist before running your app or deployment.

## Installation

```bash
go install github.com/sikalabs/vali@latest
```

## Config file

Create a `vali.json` in your project root:

```json
{
  "meta": { "version": 1 },
  "env": ["DATABASE_URL", "API_KEY", "SECRET_TOKEN"],
  "files": ["/etc/app/config.yaml", "/run/secrets/tls.key"]
}
```

vali looks for config files in this order:

1. `vali.local.json`
2. `vali.json`
3. `/vali.json`

## Usage

```
vali validate [--config <path>]
```

### Validate using auto-discovered config

```bash
vali validate
```

### Validate using a specific config file

```bash
vali validate --config /path/to/vali.json
vali validate -c /path/to/vali.json
```

## Output

On success:

```
OK
```

On failure:

```
Validation failed
  missing env:  DATABASE_URL
  missing env:  SECRET_TOKEN
  missing file: /run/secrets/tls.key
```

Exit code is `0` on success and `1` on failure, so vali works naturally in scripts and CI pipelines:

```bash
vali validate && ./start-server
```
