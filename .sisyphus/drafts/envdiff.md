# Draft: envdiff CLI Tool

## Requirements (confirmed)
- Go CLI tool for comparing/diffing .env files across environments
- All features from brainstorm: Core MVP + v1 + v2
- Project name: envdiff

## Technical Decisions (from ports project patterns)
- Module namespace: `github.com/shahadulhaider/envdiff`
- Project structure: `cmd/envdiff/main.go` + `internal/` packages
- Charmbracelet ecosystem for any TUI/interactive elements (lipgloss for styling)
- GoReleaser + Homebrew tap for distribution
- Makefile with standard targets (build, install, clean, vet, lint, cross)
- Version injection via ldflags

## Feature Set (ALL requested)

### Core
- `envdiff .env .env.production` — diff two files
- `envdiff check` — auto-compare .env vs .env.example
- `envdiff compare .env .env.staging .env.production` — multi-env matrix
- Exit codes: 0 = sync, 1 = drift

### v1
- `--mask` / `--no-values` — hide secret values in output
- `envdiff validate --schema .env.schema` — schema validation
- `envdiff git .env` — diff against git history
- `--format json|table|github` — output formats
- `--ignore "PREFIX_*"` — ignore patterns

### v2
- `envdiff init` — generate .env.example from .env
- `envdiff sync` — interactive add missing keys
- `envdiff hook install` — pre-commit hook
- `envdiff ci --require .env.example` — CI mode with annotations
- Secret detection in .env.example
- Grouped output by prefix

## Decisions Made
- CLI framework: **Cobra** (subcommand support, auto help/completions)
- Test strategy: **Tests after implementation** (unit tests for parser, diff engine, schema)
- Schema format: **TOML** (BurntSushi/toml, Go-native)
- Interactive sync: **Bubbletea TUI** (consistent with ports project style)
- Project location: /Users/msh/code/pp/envdiff

## Scope Boundaries
- INCLUDE: All features (core + v1 + v2)
- EXCLUDE: Web UI, cloud sync, secrets vault integration
