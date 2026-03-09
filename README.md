# envdiff

Compare, validate, and sync .env files across environments.

## Features

- Diff two .env files with added/removed/changed detection
- Compare N files side-by-side in a matrix view
- Check .env against .env.example for missing keys
- Validate .env against TOML schema (types, required, patterns, enums)
- Generate .env.example from existing .env
- Diff .env against its last git commit
- Interactive sync TUI to selectively apply changes
- Pre-commit hook to catch .env drift
- CI mode with GitHub Actions annotations
- JSON, table, and GitHub annotation output formats
- Secret masking with `--mask`
- Key filtering with `--ignore` glob patterns

## Requirements

- macOS or Linux
- Go 1.25 or newer (for building from source)

## Installation

### Homebrew (macOS)

```bash
brew tap shahadulhaider/tap
brew install envdiff
```

### go install

```bash
go install github.com/shahadulhaider/envdiff/cmd/envdiff@latest
```

### Download binary

Pre-built binaries for macOS and Linux are available on the [GitHub Releases](https://github.com/shahadulhaider/envdiff/releases) page.

### Build from source

```bash
git clone https://github.com/shahadulhaider/envdiff
cd envdiff
make build
./envdiff
```

## Usage

```bash
# Diff two .env files
envdiff diff .env.dev .env.prod

# JSON output
envdiff diff .env.dev .env.prod --format json

# Mask secret values
envdiff diff .env.dev .env.prod --mask

# Check .env against .env.example
envdiff check

# Compare multiple files side-by-side
envdiff compare .env.dev .env.staging .env.prod

# Validate against schema
envdiff validate --schema .env.schema.toml .env

# Generate .env.example
envdiff init

# Diff against last git commit
envdiff git .env

# Interactive sync
envdiff sync source.env target.env

# Install pre-commit hook
envdiff hook install

# CI mode with GitHub Actions annotations
envdiff ci --require .env.example
```

## Output Examples

### Table (default)

```
- A=1
~ C=3 -> 99
+ D=4

1 added, 1 removed, 1 changed
```

### Compare matrix

```
KEY         dev.env      staging.env  prod.env
-------------------------------------------------
A           1            1            1
B           2            99           <missing>
C           3            <missing>    3
D           <missing>    <missing>    4
```

## Commands

| Command | Description |
|---------|-------------|
| `diff` | Compare two .env files |
| `check` | Check .env against .env.example |
| `compare` | Compare N files side-by-side |
| `validate` | Validate .env against TOML schema |
| `init` | Generate .env.example from .env |
| `git` | Diff .env against last git commit |
| `sync` | Interactive TUI to apply changes |
| `hook` | Install/uninstall pre-commit hook |
| `ci` | CI mode with annotations |

## Global Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `table`, `json`, `github` (default: `table`) |
| `--mask` | Hide secret values in output |
| `--ignore` | Glob pattern for keys to skip |
| `--no-values` | Show keys only, no values |
| `--color` | Color mode: `auto`, `always`, `never` |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | No differences / validation passed |
| `1` | Differences found / validation failed |
| `2` | Error (file not found, parse error) |

## Schema Format

envdiff uses TOML for schema validation:

```toml
allow_extra = true

[vars.DB_HOST]
required = true
type = "string"

[vars.DB_PORT]
required = true
type = "number"
default = "5432"

[vars.LOG_LEVEL]
required = false
type = "enum"
enum = ["debug", "info", "warn", "error"]

[vars.API_KEY]
required = true
type = "string"
pattern = "^[A-Za-z0-9]{32,}$"
```

Supported types: `string`, `number`, `bool`, `url`, `email`, `enum`

## License

GNU General Public License v3.0 — see [LICENSE](LICENSE) for details.
