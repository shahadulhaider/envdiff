# envdiff — .env File Diff & Validation CLI

## TL;DR

> **Quick Summary**: Build a Go CLI tool that compares .env files across environments, validates them against TOML schemas, and provides interactive sync, CI integration, and git history diffing. Fills a confirmed gap — no dedicated tool exists for cross-environment .env consistency.
> 
> **Deliverables**:
> - `envdiff` binary with 9 subcommands: `diff`, `check`, `compare`, `validate`, `init`, `sync`, `git`, `hook`, `ci`
> - TOML-based schema validation system
> - Bubbletea interactive sync TUI
> - GoReleaser + Homebrew distribution
> - Comprehensive test suite (stdlib `testing`, table-driven)
> 
> **Estimated Effort**: Large
> **Parallel Execution**: YES — 7 waves
> **Critical Path**: Scaffolding → Types → Parser → Diff Engine → CLI Commands → Advanced Features → Distribution

---

## Context

### Original Request
Build a Go CLI tool called `envdiff` for comparing, validating, and syncing .env files across environments. Include all proposed features: core diff/check/compare, schema validation, git integration, output formats, value masking, ignore patterns, init, interactive sync, pre-commit hooks, CI mode, secret detection, and grouped output.

### Interview Summary
**Key Discussions**:
- Validated idea against existing tools — confirmed gap (direnv, dotenv-linter, dotenvx don't do cross-file diffing)
- CLI framework: **Cobra** (user's explicit choice for 8+ subcommands, despite `flag` pattern in ports/glit projects)
- Test strategy: **Tests after implementation** (unit tests for parser, diff engine, schema)
- Schema format: **TOML** (BurntSushi/toml, Go-native)
- Interactive sync: **Bubbletea TUI** (consistent with user's Charmbracelet usage)

**Research Findings**:
- direnv: Has `EnvDiff` struct but for runtime shell env changes, not file comparison
- dotenv-linter (Rust): Lints single files for formatting, not cross-file diffing
- User's `ports` project: `cmd/ports/main.go` + `internal/` + Makefile + GoReleaser + Homebrew tap
- User's `glit-work` project: stdlib `testing` only (testify explicitly forbidden), table-driven tests

### Metis Review
**Identified Gaps** (addressed):
- **Cobra vs flag conflict**: User's glit-work forbids Cobra, but user explicitly chose Cobra for envdiff → respected user's choice
- **Bubbletea version**: ports=v1, glit-work=v2 → defaulted to v2 (newer, for new project)
- **Custom parser vs godotenv**: godotenv loses ordering/line numbers/duplicates → custom parser
- **Variable interpolation**: Not resolved → default to literal (simpler, more predictable for diff tool)
- **Exit code convention**: → 0=no diff, 1=diffs found, 2=error (matches `diff(1)`)
- **Testing constraints**: → stdlib `testing` only, no testify (matches user's glit-work pattern)

---

## Work Objectives

### Core Objective
Build a complete, production-ready Go CLI tool that makes .env file management reliable across environments — from local dev to CI/CD pipelines.

### Concrete Deliverables
- `envdiff` binary with 9 subcommands
- `internal/` packages: parser, diff, output, schema, secret, git, sync, hook
- TOML schema format specification
- GoReleaser config + Homebrew formula
- GitHub Actions release workflow
- Comprehensive test suite

### Definition of Done
- [ ] `go build ./cmd/envdiff` succeeds with zero warnings
- [ ] `go test ./... -v` passes all tests
- [ ] `go vet ./...` reports zero issues
- [ ] `./envdiff --version` outputs `envdiff dev`
- [ ] `./envdiff diff a.env b.env` shows correct diffs with exit code 1
- [ ] `./envdiff validate --schema schema.toml .env` validates correctly
- [ ] `./envdiff sync source.env target.env` opens Bubbletea TUI
- [ ] `goreleaser check` passes

### Must Have
- All 9 subcommands functional
- Exit codes: 0=clean, 1=diffs/failures, 2=errors
- `--format json|table|github` output modes
- `--mask` flag hides secret values
- `--ignore` flag skips keys by pattern
- Schema validation via TOML
- CI-friendly: non-interactive by default, correct exit codes

### Must NOT Have (Guardrails)
- **NO** colors in non-TUI output unless `--color` flag is explicit or terminal is detected. CI expects plain text.
- **NO** interactive prompts in any command except `sync`. Every other command must work in scripts/CI pipes.
- **NO** `utils/`, `helpers/`, `common/` packages. Every package has a clear domain name.
- **NO** over-abstracted parser (no parser factory, no strategy pattern). One package, one `Parse()` function.
- **NO** emoji in CLI output. Use `+`/`-`/`~` prefixes for diffs (matching user's `ports` style).
- **NO** `log` package. All output to `os.Stdout` (results) or `os.Stderr` (errors).
- **NO** config file for envdiff itself. All configuration via flags.
- **NO** variable interpolation — treat `${VAR}` as literal string. Diff shows what's in the file.
- **NO** testify, gomock, or any test framework beyond stdlib `testing`.
- **NO** godotenv dependency — custom parser to preserve ordering, line numbers, duplicates.
- **NO** `go-git` — use `os/exec` for git commands (matching user's glit pattern).

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: NO (new project)
- **Automated tests**: YES — tests after implementation
- **Framework**: stdlib `testing` only (`go test`)
- **Pattern**: Table-driven tests with `t.Run()`, `t.Helper()`, subtests

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Parser/Diff/Schema**: Use Bash — `go test ./internal/<pkg>/... -v`, assert PASS
- **CLI commands**: Use Bash — build binary, run with args, assert stdout + stderr + exit code
- **Interactive TUI**: Use interactive_bash (tmux) — launch envdiff sync, send keystrokes, capture output
- **Distribution**: Use Bash — `goreleaser check`, `make build`

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — foundation):
├── Task 1: Project scaffolding + Go module + Cobra root [quick]
└── Task 2: Core types + interfaces package [quick]

Wave 2 (After Wave 1 — independent packages, parallel):
├── Task 3: .env parser + tests [deep]
├── Task 4: Output formatters (table/JSON/GitHub) + tests [unspecified-high]
└── Task 5: Secret detection patterns + tests [unspecified-low]

Wave 3 (After Wave 2 — dependent packages, parallel):
├── Task 6: Diff engine + tests (depends: parser) [deep]
└── Task 7: Schema validation + tests (depends: parser) [deep]

Wave 4 (After Wave 3 — CLI commands, MAX PARALLEL):
├── Task 8:  `diff` command + global flags [unspecified-high]
├── Task 9:  `check` command [quick]
├── Task 10: `compare` command + grouped output [unspecified-high]
├── Task 11: `validate` command [unspecified-high]
├── Task 12: `init` command [quick]
└── Task 13: `git` command [unspecified-high]

Wave 5 (After Wave 4 — advanced features, parallel):
├── Task 14: CI mode + annotations [quick]
├── Task 15: Pre-commit hook install/uninstall [quick]
└── Task 16: Interactive sync Bubbletea TUI [visual-engineering]

Wave 6 (After Wave 5 — distribution):
└── Task 17: GoReleaser + Homebrew + GitHub Actions [quick]

Wave FINAL (After ALL tasks — verification, 4 parallel):
├── Task F1: Plan compliance audit (oracle)
├── Task F2: Code quality review (unspecified-high)
├── Task F3: Real manual QA (unspecified-high)
└── Task F4: Scope fidelity check (deep)

Critical Path: T1 → T3 → T6 → T8 → T14/T16 → T17 → F1-F4
Parallel Speedup: ~65% faster than sequential
Max Concurrent: 6 (Wave 4)
```

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| 1 | — | 2-17 | 1 |
| 2 | 1 | 3-17 | 1 |
| 3 | 2 | 6, 7, 8-13 | 2 |
| 4 | 2 | 8-13 | 2 |
| 5 | 2 | 8 | 2 |
| 6 | 3 | 8-16 | 3 |
| 7 | 3 | 11, 12 | 3 |
| 8 | 6, 4, 5 | 14 | 4 |
| 9 | 6 | 14 | 4 |
| 10 | 6, 4 | 14 | 4 |
| 11 | 7, 4 | 14 | 4 |
| 12 | 3, 7 | — | 4 |
| 13 | 6 | — | 4 |
| 14 | 8-11 | 17 | 5 |
| 15 | 8 | 17 | 5 |
| 16 | 6 | 17 | 5 |
| 17 | 14-16 | F1-F4 | 6 |

### Agent Dispatch Summary

- **Wave 1**: **2** — T1 → `quick`, T2 → `quick`
- **Wave 2**: **3** — T3 → `deep`, T4 → `unspecified-high`, T5 → `unspecified-low`
- **Wave 3**: **2** — T6 → `deep`, T7 → `deep`
- **Wave 4**: **6** — T8 → `unspecified-high`, T9 → `quick`, T10 → `unspecified-high`, T11 → `unspecified-high`, T12 → `quick`, T13 → `unspecified-high`
- **Wave 5**: **3** — T14 → `quick`, T15 → `quick`, T16 → `visual-engineering`
- **Wave 6**: **1** — T17 → `quick`
- **FINAL**: **4** — F1 → `oracle`, F2 → `unspecified-high`, F3 → `unspecified-high`, F4 → `deep`

---

## TODOs

> Implementation + Tests = ONE Task. Never separate.
> EVERY task MUST have: Recommended Agent Profile + Parallelization info + QA Scenarios.

- [ ] 1. Project Scaffolding + Cobra Root

  **What to do**:
  - Create project directory at `/Users/msh/code/pp/envdiff`
  - `go mod init github.com/shahadulhaider/envdiff` with Go 1.25
  - Create directory structure: `cmd/envdiff/`, `internal/{env,parser,diff,output,schema,secret,git,sync,hook}`, `testdata/`
  - Install Cobra: `go get github.com/spf13/cobra`
  - Create `cmd/envdiff/main.go` with Cobra root command, `--version` flag
  - Create `cmd/envdiff/root.go` with root command setup, persistent flags (`--format`, `--mask`, `--ignore`, `--no-values`, `--color`)
  - Create `Makefile` following ports pattern: `build install clean vet lint cross` targets with ldflags version injection
  - Create `.goreleaser.yaml` following ports pattern: v2 format, darwin+linux, amd64+arm64, homebrew tap at `shahadulhaider/homebrew-tap`
  - Create `.gitignore` (binary, dist/, .env files with actual secrets)
  - `git init` + initial commit

  **Must NOT do**:
  - Add any subcommands yet (only root + --version)
  - Add any business logic
  - Create a config file for envdiff itself
  - Use `log` package

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Boilerplate scaffolding, copying patterns from existing project
  - **Skills**: [`git-master`]
    - `git-master`: Git init + initial commit
  - **Skills Evaluated but Omitted**:
    - `playwright`: No browser interaction
    - `frontend-ui-ux`: No UI yet

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 2 once go.mod exists — but Task 2 needs go.mod, so sequential)
  - **Parallel Group**: Wave 1 (sequential with Task 2)
  - **Blocks**: Tasks 2-17
  - **Blocked By**: None (can start immediately)

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/cmd/ports/main.go` — Entry point pattern: flag parsing, version handling, exit codes
  - `/Users/msh/code/pp/ports/Makefile` — Makefile targets and ldflags pattern (copy exactly, change binary name)
  - `/Users/msh/code/pp/ports/.goreleaser.yaml` — GoReleaser v2 config (copy exactly, change binary/desc/homepage)
  - `/Users/msh/code/pp/ports/go.mod` — Module pattern and Go version

  **External References**:
  - Cobra docs: https://cobra.dev/ — Root command setup, persistent flags, version template

  **WHY Each Reference Matters**:
  - `main.go`: Shows the user's preferred entry point style — how they handle flags, version output format, exit codes
  - `Makefile`: User has specific ldflags and target naming — MUST match exactly
  - `.goreleaser.yaml`: Homebrew tap config at `shahadulhaider/homebrew-tap` — MUST match owner/name/token pattern
  - `go.mod`: Go version (`go 1.25`) — use same version

  **Acceptance Criteria**:

  - [ ] `go build ./cmd/envdiff` succeeds
  - [ ] `./envdiff --version` outputs `envdiff dev`
  - [ ] `make build` produces `envdiff` binary
  - [ ] `goreleaser check` passes
  - [ ] Directory structure exists: cmd/envdiff/, internal/{env,parser,diff,output,schema,secret,git,sync,hook}
  - [ ] `go vet ./...` reports zero issues

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Binary builds and runs
    Tool: Bash
    Preconditions: Go 1.25+ installed, project scaffolded
    Steps:
      1. Run `go build -o envdiff ./cmd/envdiff`
      2. Run `./envdiff --version`
      3. Assert stdout contains "envdiff dev"
      4. Assert exit code is 0
    Expected Result: Binary builds cleanly, version output matches "envdiff dev"
    Failure Indicators: Build errors, missing imports, wrong version format
    Evidence: .sisyphus/evidence/task-1-binary-builds.txt

  Scenario: Makefile targets work
    Tool: Bash
    Preconditions: Makefile exists
    Steps:
      1. Run `make clean`
      2. Run `make build`
      3. Assert `envdiff` binary exists
      4. Run `make vet`
      5. Assert exit code 0
    Expected Result: All Makefile targets execute without errors
    Failure Indicators: Missing targets, build failures
    Evidence: .sisyphus/evidence/task-1-makefile.txt

  Scenario: GoReleaser config valid
    Tool: Bash
    Preconditions: .goreleaser.yaml exists, goreleaser installed
    Steps:
      1. Run `goreleaser check`
      2. Assert exit code 0
    Expected Result: Config passes validation
    Failure Indicators: YAML errors, missing fields
    Evidence: .sisyphus/evidence/task-1-goreleaser.txt

  Scenario: Unknown command shows help (not crash)
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./envdiff nonexistent 2>&1`
      2. Assert stderr contains "unknown command"
      3. Assert exit code is non-zero
    Expected Result: Cobra shows helpful error, does not panic
    Failure Indicators: Panic, stack trace, silent failure
    Evidence: .sisyphus/evidence/task-1-unknown-cmd.txt
  ```

  **Commit**: YES
  - Message: `chore: scaffold envdiff project`
  - Files: `go.mod, go.sum, cmd/envdiff/main.go, cmd/envdiff/root.go, Makefile, .goreleaser.yaml, .gitignore, internal/*/`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 2. Core Types + Interfaces Package

  **What to do**:
  - Create `internal/env/types.go` with all shared types:
    - `EnvEntry` struct: `Key string`, `Value string`, `LineNum int`, `Comment string`, `IsExported bool` (has `export` prefix), `Raw string` (original line)
    - `EnvFile` struct: `Entries []EnvEntry`, `Path string`, `Comments []string` (standalone comments), `Duplicates []EnvEntry` (duplicate keys detected)
    - Methods on `EnvFile`: `Keys() []string`, `Get(key string) (EnvEntry, bool)`, `Len() int`
    - `DiffType` enum: `DiffAdded`, `DiffRemoved`, `DiffChanged`
    - `DiffEntry` struct: `Key string`, `Type DiffType`, `Left *EnvEntry` (nil if added), `Right *EnvEntry` (nil if removed)
    - `DiffResult` struct: `Left string` (file path), `Right string` (file path), `Entries []DiffEntry`, methods: `HasDiffs() bool`, `Added() []DiffEntry`, `Removed() []DiffEntry`, `Changed() []DiffEntry`
    - `FormatType` enum: `FormatTable`, `FormatJSON`, `FormatGitHub`
    - `Formatter` interface: `Format(result *DiffResult, w io.Writer) error`
    - `SchemaRule` struct: `Required bool`, `Type string`, `Pattern string`, `Default string`, `Enum []string`
    - `SchemaConfig` struct: `AllowExtra bool`, `Rules map[string]SchemaRule`
    - `ValidationError` struct: `Key string`, `Message string`, `Line int`
    - `ValidationResult` struct: `Errors []ValidationError`, `Warnings []ValidationError`, methods: `IsValid() bool`
  - All types in one file — no splitting across multiple files
  - Use `iota` for enums, `String()` method on each

  **Must NOT do**:
  - Put types in multiple packages (everything in `internal/env`)
  - Add any implementation logic (just types and interfaces)
  - Import external dependencies
  - Create utility functions

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Pure type definitions, no logic
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - All skills: No domain overlap with type definition

  **Parallelization**:
  - **Can Run In Parallel**: NO (needs go.mod from Task 1)
  - **Parallel Group**: Wave 1 (sequential after Task 1)
  - **Blocks**: Tasks 3-17
  - **Blocked By**: Task 1

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/internal/diff/diff.go` — Shows how user defines diff-related types (DiffEntry, formatting)
  - `/Users/msh/code/pp/ports/internal/scanner/proc.go` — Shows user's struct style and method conventions

  **WHY Each Reference Matters**:
  - `diff.go`: The user already has a diff concept in ports — match the naming style and struct conventions
  - `proc.go`: Shows how user organizes types + methods in one file

  **Acceptance Criteria**:

  - [ ] `go build ./internal/env/...` compiles
  - [ ] All types have `String()` methods where applicable
  - [ ] `Formatter` interface is defined with `Format(*DiffResult, io.Writer) error`
  - [ ] No external imports (stdlib only)
  - [ ] `go vet ./internal/env/...` passes

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Types package compiles cleanly
    Tool: Bash
    Preconditions: Task 1 complete, go.mod exists
    Steps:
      1. Run `go build ./internal/env/...`
      2. Assert exit code 0
      3. Run `go vet ./internal/env/...`
      4. Assert exit code 0
    Expected Result: Package compiles with zero errors and zero vet warnings
    Failure Indicators: Import cycles, undefined types, unused imports
    Evidence: .sisyphus/evidence/task-2-types-compile.txt

  Scenario: Types are usable from other packages
    Tool: Bash
    Preconditions: types.go exists
    Steps:
      1. Create a temporary test file that imports `internal/env` and instantiates each type
      2. Run `go build` on the test file
      3. Assert compilation succeeds
      4. Clean up temp file
    Expected Result: All exported types are accessible and instantiable
    Failure Indicators: Unexported fields that should be exported, circular deps
    Evidence: .sisyphus/evidence/task-2-types-usable.txt
  ```

  **Commit**: YES
  - Message: `feat: add core types and interfaces`
  - Files: `internal/env/types.go`
  - Pre-commit: `go build ./... && go vet ./...`

- [ ] 3. .env File Parser + Tests

  **What to do**:
  - Implement `internal/parser/parser.go`:
    - `Parse(r io.Reader) (*env.EnvFile, error)` — main entry point
    - `ParseFile(path string) (*env.EnvFile, error)` — convenience wrapper
    - Custom line-by-line parser (NOT godotenv) that preserves:
      - Key ordering (entries slice, not map)
      - Line numbers for every entry
      - Standalone comments (lines starting with `#`)
      - Inline comments (after value, outside quotes)
      - Duplicate key detection (add to `Duplicates` slice, keep last value in main entries)
      - Empty values (`KEY=` → value is `""`)
      - Quoted values (single and double quotes stripped)
      - `export` prefix (`export KEY=val` → `IsExported=true`, key is `KEY`)
      - BOM marker handling (strip UTF-8 BOM from first line)
      - Windows line endings (`\r\n` → normalize to `\n`)
    - Values are treated as LITERAL strings — no `${VAR}` interpolation
  - Implement `internal/parser/parser_test.go`:
    - Table-driven tests using `t.Run()` for each case
    - Test cases: empty input, comments-only file, simple `KEY=value`, `KEY="double quoted"`, `KEY='single quoted'`, `export KEY=val`, `KEY=` (empty value), `KEY` (no equals sign → error or skip), duplicate keys, inline comments (`KEY=val # comment`), whitespace around `=`, leading/trailing whitespace on lines, BOM marker, Windows `\r\n` endings, multiline values NOT supported (error or literal), lines with only whitespace
    - Use `strings.NewReader()` for input, not file I/O

  **Must NOT do**:
  - Import godotenv or any external parser
  - Implement variable interpolation (`${VAR}`)
  - Support heredoc syntax
  - Over-abstract (no parser interface, no parser options struct — just `Parse()`)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core component requiring careful edge case handling, extensive test coverage
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - All skills: Pure Go stdlib work, no external tooling needed

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 4, 5)
  - **Parallel Group**: Wave 2 (with Tasks 4, 5)
  - **Blocks**: Tasks 6, 7, 8-13
  - **Blocked By**: Task 2

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/internal/scanner/proc.go` — Shows user's Go function/struct style
  - `/Users/msh/code/pp/ports/internal/scanner/proc_darwin.go` — Platform-specific patterns (not needed here, but shows code organization)

  **API/Type References**:
  - `internal/env/types.go:EnvEntry` — The struct parser must populate (Key, Value, LineNum, Comment, IsExported, Raw)
  - `internal/env/types.go:EnvFile` — The struct Parse() returns (Entries, Path, Comments, Duplicates)

  **External References**:
  - https://hexdocs.pm/dotenvy/dotenv-file-format.html — Informal .env format reference
  - direnv's parser: `github.com/direnv/direnv/internal/cmd/env.go` — How direnv parses env (reference only, don't copy)

  **WHY Each Reference Matters**:
  - `proc.go`: Match the user's Go style — how they name functions, handle errors, organize methods
  - `types.go`: Parser must return exactly these types — fields must match
  - dotenv format ref: Defines the subset of .env syntax we support

  **Acceptance Criteria**:

  - [ ] `go test ./internal/parser/... -v` — all tests PASS
  - [ ] Parser handles all edge cases listed above
  - [ ] No external dependencies (stdlib only)
  - [ ] Key ordering is preserved (test: parse file with keys Z, A, M → entries are in Z, A, M order)
  - [ ] Duplicate keys produce entries in `Duplicates` slice AND keep last value in main entries

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Parser handles standard .env file
    Tool: Bash
    Preconditions: Parser implemented
    Steps:
      1. Create testdata/standard.env with content:
         # Database config
         DB_HOST=localhost
         DB_PORT=5432
         DB_NAME="myapp_dev"
         export API_KEY='secret123'
         DEBUG=
      2. Run `go test ./internal/parser/... -v -run TestParseStandard`
      3. Assert: 5 entries parsed, correct keys/values, ordering preserved
      4. Assert: DB_NAME value is `myapp_dev` (quotes stripped)
      5. Assert: API_KEY has IsExported=true
      6. Assert: DEBUG has empty string value
    Expected Result: All assertions pass
    Failure Indicators: Wrong entry count, quotes not stripped, ordering lost
    Evidence: .sisyphus/evidence/task-3-parse-standard.txt

  Scenario: Parser detects duplicate keys
    Tool: Bash
    Preconditions: Parser implemented
    Steps:
      1. Run test with input: "KEY=first\nKEY=second\n"
      2. Assert: main entries has KEY=second (last wins)
      3. Assert: Duplicates slice contains KEY=first
    Expected Result: Duplicate detected, last value wins, duplicate recorded
    Failure Indicators: Silent overwrite without recording, first value wins
    Evidence: .sisyphus/evidence/task-3-duplicates.txt

  Scenario: Parser handles BOM and Windows line endings
    Tool: Bash
    Preconditions: Parser implemented
    Steps:
      1. Run test with input containing UTF-8 BOM (0xEF,0xBB,0xBF) + "KEY=val\r\n"
      2. Assert: KEY parsed correctly (BOM stripped, \r stripped)
    Expected Result: BOM and \r\n handled transparently
    Failure Indicators: BOM appears in key name, \r in value
    Evidence: .sisyphus/evidence/task-3-bom-crlf.txt

  Scenario: Parser rejects malformed lines gracefully
    Tool: Bash
    Preconditions: Parser implemented
    Steps:
      1. Run test with input containing lines with no = sign (not comments)
      2. Assert: malformed lines are skipped (not crash)
      3. Assert: valid lines still parsed correctly
    Expected Result: Graceful handling of malformed input
    Failure Indicators: Panic, parse error on entire file, missing valid entries
    Evidence: .sisyphus/evidence/task-3-malformed.txt
  ```

  **Commit**: YES
  - Message: `feat: implement .env file parser`
  - Files: `internal/parser/parser.go, internal/parser/parser_test.go, testdata/standard.env`
  - Pre-commit: `go test ./internal/parser/... -v && go vet ./...`

---

- [ ] 4. Output Formatters (Table / JSON / GitHub) + Tests

  **What to do**:
  - Implement `internal/output/formatter.go` — factory function `NewFormatter(format env.FormatType) env.Formatter`
  - Implement `internal/output/table.go` — `TableFormatter` struct implementing `Formatter`:
    - `+` prefix for added keys (green if color enabled)
    - `-` prefix for removed keys (red if color enabled)
    - `~` prefix for changed keys (yellow if color enabled)
    - Aligned columns for key and value
    - Summary line: `N added, N removed, N changed`
    - When `--mask` is active, values show `****`
    - When `--no-values` is active, only keys shown
    - Use lipgloss for styled output when terminal detected
  - Implement `internal/output/json.go` — `JSONFormatter`:
    - Output valid JSON object: `{"left": "file1", "right": "file2", "added": [...], "removed": [...], "changed": [...]}`
    - Each entry: `{"key": "KEY", "left_value": "...", "right_value": "..."}`
    - Values masked with `"****"` when mask enabled
  - Implement `internal/output/github.go` — `GitHubFormatter`:
    - GitHub Actions annotation format: `::warning file={file},line={line}::{message}`
    - Added keys: `::notice` annotations
    - Removed keys: `::error` annotations
    - Changed keys: `::warning` annotations
  - Implement `internal/output/table_test.go`, `json_test.go`, `github_test.go`:
    - Test each formatter with known DiffResult, assert exact output string
    - Test masking behavior
    - Test JSON validity with `json.Valid()`

  **Must NOT do**:
  - Add color by default to non-TTY output. Colors ONLY when stdout is a terminal or `--color` flag is set.
  - Use emoji in any formatter output
  - Import charmbracelet/bubbletea (lipgloss only for styling)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Multiple implementations of same interface, need to handle color/TTY detection correctly
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: This is terminal output, not browser UI

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 3, 5)
  - **Parallel Group**: Wave 2 (with Tasks 3, 5)
  - **Blocks**: Tasks 8-13
  - **Blocked By**: Task 2

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/internal/diff/diff.go` — How user formats diff output (`+`/`-` prefixed lines)
  - `/Users/msh/code/pp/ports/internal/tui/styles.go` — Lipgloss styling patterns the user prefers

  **API/Type References**:
  - `internal/env/types.go:Formatter` — Interface to implement: `Format(*DiffResult, io.Writer) error`
  - `internal/env/types.go:DiffResult` — Input struct with `Added()`, `Removed()`, `Changed()` methods
  - `internal/env/types.go:FormatType` — Enum: `FormatTable`, `FormatJSON`, `FormatGitHub`

  **External References**:
  - GitHub Actions annotations: https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-a-warning-message
  - Lipgloss: https://github.com/charmbracelet/lipgloss — Terminal styling

  **WHY Each Reference Matters**:
  - `diff.go`: Shows user's preferred diff output style — `+`/`-` prefixes, no emoji
  - `styles.go`: Shows which lipgloss features and color choices the user prefers
  - GitHub docs: Exact annotation format required for CI integration

  **Acceptance Criteria**:

  - [ ] `go test ./internal/output/... -v` — all tests PASS
  - [ ] JSON output passes `json.Valid()` check
  - [ ] Table output uses `+`/`-`/`~` prefixes
  - [ ] GitHub output uses `::error`/`::warning`/`::notice` format
  - [ ] Masking replaces values with `****`

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: JSON output is valid JSON
    Tool: Bash
    Preconditions: JSON formatter implemented
    Steps:
      1. Run `go test ./internal/output/... -v -run TestJSONFormat`
      2. In test: create DiffResult with 1 added, 1 removed, 1 changed
      3. Format to buffer, run json.Valid() on output
      4. Assert valid JSON with correct structure
    Expected Result: Output parses as valid JSON with expected fields
    Failure Indicators: Invalid JSON, missing fields, wrong structure
    Evidence: .sisyphus/evidence/task-4-json-valid.txt

  Scenario: Table output uses correct prefixes
    Tool: Bash
    Preconditions: Table formatter implemented
    Steps:
      1. Run `go test ./internal/output/... -v -run TestTableFormat`
      2. Assert output contains lines starting with `+ ` (added), `- ` (removed), `~ ` (changed)
      3. Assert summary line present: "N added, N removed, N changed"
    Expected Result: Diff prefixes match convention, summary accurate
    Failure Indicators: Wrong prefixes, missing summary, emoji in output
    Evidence: .sisyphus/evidence/task-4-table-output.txt

  Scenario: Masking hides values
    Tool: Bash
    Preconditions: Formatters with masking implemented
    Steps:
      1. Run test with DiffResult containing value "secret123"
      2. Format with masking enabled
      3. Assert output contains "****" and NOT "secret123"
    Expected Result: Secret values replaced with mask
    Failure Indicators: Original value visible in output
    Evidence: .sisyphus/evidence/task-4-masking.txt
  ```

  **Commit**: YES
  - Message: `feat: add output formatters (table/json/github)`
  - Files: `internal/output/formatter.go, table.go, json.go, github.go, *_test.go`
  - Pre-commit: `go test ./internal/output/... -v && go vet ./...`

- [ ] 5. Secret Detection Patterns + Tests

  **What to do**:
  - Implement `internal/secret/secret.go`:
    - `IsSecret(key, value string) bool` — returns true if value looks like a secret
    - `DetectSecrets(entries []env.EnvEntry) []env.EnvEntry` — returns entries whose values are likely secrets
    - Detection patterns:
      - AWS access keys: `AKIA[0-9A-Z]{16}`
      - AWS secret keys: 40-char base64 strings
      - API tokens/keys: high-entropy strings (>4.5 Shannon entropy, >20 chars)
      - Database URLs with embedded passwords: `://user:pass@`
      - Private keys: `-----BEGIN (RSA |EC |)PRIVATE KEY-----`
      - Generic patterns: keys containing `SECRET`, `PASSWORD`, `TOKEN`, `API_KEY`, `PRIVATE` (case-insensitive)
    - `MaskValue(value string) string` — returns `****` (fixed length, don't reveal original length)
  - Implement `internal/secret/secret_test.go`:
    - Table-driven tests for each pattern
    - Test false positives: `PORT=5432` should NOT be detected as secret
    - Test key-name heuristic: `DB_PASSWORD=anything` IS a secret regardless of value

  **Must NOT do**:
  - Over-engineer entropy calculation (Shannon entropy is sufficient)
  - Create a configurable secrets engine (hardcoded patterns are fine)
  - Use external secret detection libraries

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Pattern matching with regex, moderate complexity
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 3, 4)
  - **Parallel Group**: Wave 2 (with Tasks 3, 4)
  - **Blocks**: Task 8
  - **Blocked By**: Task 2

  **References** (CRITICAL):

  **API/Type References**:
  - `internal/env/types.go:EnvEntry` — Input type: Key, Value fields to scan

  **External References**:
  - truffleHog patterns: https://github.com/trufflesecurity/trufflehog — Secret detection patterns reference
  - detect-secrets: https://github.com/Yelp/detect-secrets — Python tool for reference patterns

  **WHY Each Reference Matters**:
  - truffleHog: Industry-standard secret patterns — use their regex patterns as inspiration
  - detect-secrets: Shows key-name heuristics (matching by variable name, not just value)

  **Acceptance Criteria**:

  - [ ] `go test ./internal/secret/... -v` — all tests PASS
  - [ ] AWS key pattern detected correctly
  - [ ] High-entropy strings detected
  - [ ] `PORT=5432` NOT flagged as secret
  - [ ] `DB_PASSWORD=anything` flagged as secret (key-name heuristic)
  - [ ] `MaskValue()` returns fixed-length `****`

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Detects AWS access key
    Tool: Bash
    Preconditions: Secret detection implemented
    Steps:
      1. Run test with value "AKIAIOSFODNN7EXAMPLE"
      2. Assert IsSecret() returns true
    Expected Result: AWS key pattern matched
    Failure Indicators: False negative on known AWS key format
    Evidence: .sisyphus/evidence/task-5-aws-key.txt

  Scenario: Does not flag normal values
    Tool: Bash
    Preconditions: Secret detection implemented
    Steps:
      1. Run test with entries: PORT=5432, HOST=localhost, DEBUG=true, APP_NAME=myapp
      2. Assert none are flagged as secrets
    Expected Result: Zero false positives on normal config values
    Failure Indicators: Normal values incorrectly flagged
    Evidence: .sisyphus/evidence/task-5-false-positives.txt
  ```

  **Commit**: YES
  - Message: `feat: add secret detection patterns`
  - Files: `internal/secret/secret.go, internal/secret/secret_test.go`
  - Pre-commit: `go test ./internal/secret/... -v && go vet ./...`

---

- [ ] 6. Diff Engine + Tests

  **What to do**:
  - Implement `internal/diff/diff.go`:
    - `Diff(left, right *env.EnvFile) *env.DiffResult` — compare two parsed env files
    - Logic:
      - Keys in right but not left → `DiffAdded`
      - Keys in left but not right → `DiffRemoved`
      - Keys in both but values differ → `DiffChanged` (with both values)
      - Keys in both with same value → not included in result
    - Output entries ordered: removed first, then changed, then added (matches `diff` convention)
    - `MultiDiff(files []*env.EnvFile) *MultiDiffResult` — compare N files, return matrix
    - `MultiDiffResult` struct: `Files []string`, `Keys []string` (union of all keys), `Matrix map[string]map[string]*string` (key → filename → value, nil if missing)
  - Implement `internal/diff/diff_test.go`:
    - Table-driven tests:
      - Identical files → empty DiffResult, `HasDiffs()` returns false
      - All added → only `DiffAdded` entries
      - All removed → only `DiffRemoved` entries
      - Mixed changes → correct categorization
      - Empty value vs missing key → treated as different (KEY= is present, missing KEY is absent)
      - Multi-diff with 3+ files → correct matrix

  **Must NOT do**:
  - Implement value-level diffing (word-by-word diff of values) — just old/new values
  - Add filtering logic here (--ignore is CLI concern, not diff engine)
  - Import output package (diff engine produces data, formatters consume it)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core algorithm, needs careful edge case handling, extensive tests
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 7)
  - **Parallel Group**: Wave 3 (with Task 7)
  - **Blocks**: Tasks 8-16
  - **Blocked By**: Task 3

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/internal/diff/diff.go` — User's existing diff implementation (for ports, not .env — but shows naming/style)
  - `github.com/direnv/direnv/internal/cmd/env_diff.go:BuildEnvDiff` — How direnv computes env diffs (Prev/Next maps)

  **API/Type References**:
  - `internal/env/types.go:EnvFile` — Input: `Entries []EnvEntry`, `Get(key) (EnvEntry, bool)`
  - `internal/env/types.go:DiffResult` — Output: `Entries []DiffEntry`, `HasDiffs()`, `Added()`, `Removed()`, `Changed()`
  - `internal/env/types.go:DiffEntry` — Per-key diff: `Key, Type, Left, Right`

  **WHY Each Reference Matters**:
  - `ports/diff/diff.go`: Match the user's diff naming style
  - direnv's `BuildEnvDiff`: Algorithm reference — iterate both maps, classify as prev-only/next-only/changed
  - Types: MUST return exactly these types — the output package and CLI commands depend on them

  **Acceptance Criteria**:

  - [ ] `go test ./internal/diff/... -v` — all tests PASS
  - [ ] Identical files produce `HasDiffs() == false`
  - [ ] Added/removed/changed correctly categorized
  - [ ] `KEY=` (empty) treated as present (different from key not existing)
  - [ ] MultiDiff produces correct matrix for 3+ files

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Diff finds added, removed, and changed keys
    Tool: Bash
    Preconditions: Diff engine and parser implemented
    Steps:
      1. Run test with left={A=1, B=2, C=3} right={B=2, C=99, D=4}
      2. Assert: Removed=[A], Changed=[C (3→99)], Added=[D]
      3. Assert: HasDiffs() == true
    Expected Result: All three diff types correctly detected
    Failure Indicators: Missing entries, wrong categorization, B incorrectly flagged
    Evidence: .sisyphus/evidence/task-6-diff-mixed.txt

  Scenario: Identical files produce no diff
    Tool: Bash
    Preconditions: Diff engine implemented
    Steps:
      1. Run test with left={A=1, B=2} right={A=1, B=2}
      2. Assert: HasDiffs() == false, len(Entries) == 0
    Expected Result: No false positives on identical inputs
    Failure Indicators: Non-empty diff result
    Evidence: .sisyphus/evidence/task-6-identical.txt

  Scenario: Empty value vs missing key are different
    Tool: Bash
    Preconditions: Diff engine implemented
    Steps:
      1. Run test with left={KEY=} right={} (KEY missing entirely)
      2. Assert: KEY is in Removed list (it existed in left with empty value, absent in right)
    Expected Result: Empty string value is NOT the same as key not existing
    Failure Indicators: Empty value treated as missing
    Evidence: .sisyphus/evidence/task-6-empty-vs-missing.txt
  ```

  **Commit**: YES
  - Message: `feat: implement diff engine`
  - Files: `internal/diff/diff.go, internal/diff/diff_test.go`
  - Pre-commit: `go test ./internal/diff/... -v && go vet ./...`

---

- [ ] 7. TOML Schema Validation + Tests

  **What to do**:
  - Implement `internal/schema/schema.go`:
    - `LoadSchema(path string) (*env.SchemaConfig, error)` — parse TOML schema file
    - `Validate(envFile *env.EnvFile, schema *env.SchemaConfig) *env.ValidationResult` — validate env file against schema
    - TOML schema format:
      ```toml
      allow_extra = true  # allow keys not in schema (default: true)

      [vars.DB_HOST]
      required = true
      type = "string"

      [vars.DB_PORT]
      required = true
      type = "number"
      default = "5432"

      [vars.API_KEY]
      required = true
      type = "string"
      pattern = "^[A-Za-z0-9]{32,}$"

      [vars.LOG_LEVEL]
      required = false
      type = "enum"
      enum = ["debug", "info", "warn", "error"]

      [vars.ENABLE_CACHE]
      type = "bool"

      [vars.DATABASE_URL]
      required = true
      type = "url"
      ```
    - Type validation:
      - `string`: any non-empty value
      - `number`/`int`: parseable as integer
      - `bool`: one of `true`, `false`, `1`, `0`, `yes`, `no`
      - `url`: parseable by `url.Parse()` with scheme
      - `email`: contains `@` with domain
      - `enum`: value matches one of `enum` list
    - Pattern validation: value matches regex in `pattern` field
    - Required validation: key must exist and have non-empty value
  - Implement `internal/schema/schema_test.go`:
    - Test schema loading from TOML string
    - Test each type validator
    - Test required/optional
    - Test `allow_extra=false` rejecting unknown keys
    - Test pattern matching
    - Test enum validation

  **Must NOT do**:
  - Support nested schemas or groups
  - Add schema inheritance or includes
  - Validate value semantics beyond type (e.g., don't check if URL is reachable)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Schema parsing + validation logic + comprehensive tests
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 6)
  - **Parallel Group**: Wave 3 (with Task 6)
  - **Blocks**: Tasks 11, 12
  - **Blocked By**: Task 3

  **References** (CRITICAL):

  **API/Type References**:
  - `internal/env/types.go:SchemaConfig` — Struct: `AllowExtra bool`, `Rules map[string]SchemaRule`
  - `internal/env/types.go:SchemaRule` — Struct: `Required, Type, Pattern, Default, Enum`
  - `internal/env/types.go:ValidationResult` — Output: `Errors, Warnings, IsValid()`

  **External References**:
  - BurntSushi/toml docs: https://github.com/BurntSushi/toml — TOML parsing API for Go

  **WHY Each Reference Matters**:
  - Types: Schema package must populate exactly these types
  - BurntSushi/toml: The TOML library chosen by user — use its unmarshaling API

  **Acceptance Criteria**:

  - [ ] `go test ./internal/schema/... -v` — all tests PASS
  - [ ] Valid .env against schema → `IsValid() == true`
  - [ ] Missing required key → `IsValid() == false` with specific error message
  - [ ] Wrong type (number field with "abc") → validation error
  - [ ] `allow_extra=false` + unknown key → validation error
  - [ ] Pattern mismatch → validation error with pattern shown

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Valid env passes schema validation
    Tool: Bash
    Preconditions: Schema validation implemented
    Steps:
      1. Run test with schema requiring DB_HOST(string,required) + DB_PORT(number,required)
      2. Validate env with DB_HOST=localhost, DB_PORT=5432
      3. Assert IsValid() == true, zero errors
    Expected Result: Valid config passes validation
    Failure Indicators: False validation error
    Evidence: .sisyphus/evidence/task-7-valid-schema.txt

  Scenario: Missing required key fails validation
    Tool: Bash
    Preconditions: Schema validation implemented
    Steps:
      1. Run test with schema requiring API_KEY(string,required)
      2. Validate env WITHOUT API_KEY
      3. Assert IsValid() == false
      4. Assert error message mentions "API_KEY" and "required"
    Expected Result: Clear error identifying missing key
    Failure Indicators: Silent pass, generic error without key name
    Evidence: .sisyphus/evidence/task-7-missing-required.txt

  Scenario: Type mismatch detected
    Tool: Bash
    Preconditions: Schema validation implemented
    Steps:
      1. Run test with schema defining PORT as type="number"
      2. Validate env with PORT=not_a_number
      3. Assert validation error mentioning type mismatch
    Expected Result: Type error with key name and expected type
    Failure Indicators: No error on invalid type
    Evidence: .sisyphus/evidence/task-7-type-mismatch.txt
  ```

  **Commit**: YES
  - Message: `feat: add TOML schema validation`
  - Files: `internal/schema/schema.go, internal/schema/schema_test.go`
  - Pre-commit: `go test ./internal/schema/... -v && go vet ./...`

- [ ] 8. `diff` Command + Global Flags

  **What to do**:
  - Implement `cmd/envdiff/diff.go` — Cobra `diff` subcommand:
    - `envdiff diff <file1> <file2>` — compare two .env files
    - Wire: parser → diff engine → formatter → stdout
    - Global flags (on root, inherited by all commands):
      - `--format string` (table|json|github, default: table)
      - `--mask` (bool, mask secret values in output)
      - `--no-values` (bool, show keys only, no values)
      - `--ignore string` (glob pattern for keys to skip, e.g., `"DEBUG_*"`)
      - `--color string` (auto|always|never, default: auto — detect TTY)
    - Implement ignore filtering: after diff, before formatting, remove entries whose keys match `--ignore` glob
    - Implement masking: if `--mask`, call `secret.MaskValue()` on all values in output; if `--no-values`, remove values entirely
    - Exit codes: 0 = no differences, 1 = differences found, 2 = error (file not found, parse error)
    - Errors to stderr, diff output to stdout

  **Must NOT do**:
  - Add interactive prompts
  - Print colors when stdout is not a TTY (unless `--color always`)
  - Use `log` package for error output (use `fmt.Fprintf(os.Stderr, ...)`)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Wiring multiple packages together, flag handling, exit code logic
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 9-13)
  - **Parallel Group**: Wave 4 (with Tasks 9-13)
  - **Blocks**: Task 14
  - **Blocked By**: Tasks 6, 4, 5

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/cmd/ports/main.go:29-31` — How user handles diff mode with exit codes
  - `/Users/msh/code/pp/ports/internal/diff/diff.go` — User's diff output style

  **API/Type References**:
  - `internal/parser/parser.go:ParseFile(path)` — Parse input files
  - `internal/diff/diff.go:Diff(left, right)` — Compute diff
  - `internal/output/formatter.go:NewFormatter(format)` — Get formatter
  - `internal/secret/secret.go:MaskValue(value)` — Mask secret values
  - `internal/env/types.go:FormatType` — Format enum for --format flag

  **External References**:
  - Cobra command docs: https://cobra.dev/#concepts — Command + flag setup
  - `filepath.Match()` — Go stdlib glob matching for --ignore patterns

  **WHY Each Reference Matters**:
  - `main.go:29-31`: Shows how user wires diff mode to exit code — `os.Exit(diff.RunDiffMode(portFlag))`
  - Parser/diff/output: These are the packages being wired — must call correct functions with correct types

  **Acceptance Criteria**:

  - [ ] `./envdiff diff a.env b.env` shows diff with `+`/`-`/`~` prefixes
  - [ ] `./envdiff diff a.env b.env --format json` outputs valid JSON
  - [ ] `./envdiff diff a.env b.env --mask` hides values
  - [ ] `./envdiff diff a.env b.env --ignore "DEBUG_*"` skips matching keys
  - [ ] Exit code 0 when files are identical
  - [ ] Exit code 1 when differences exist
  - [ ] Exit code 2 when file not found

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: diff command shows differences with correct exit code
    Tool: Bash
    Preconditions: Binary built, testdata files with known differences
    Steps:
      1. Create testdata/left.env: "A=1\nB=2\nC=3"
      2. Create testdata/right.env: "B=2\nC=99\nD=4"
      3. Run `./envdiff diff testdata/left.env testdata/right.env`
      4. Assert stdout contains "- A" (removed), "~ C" (changed), "+ D" (added)
      5. Assert exit code is 1
    Expected Result: All diffs shown, exit code 1
    Failure Indicators: Wrong exit code, missing entries, wrong prefixes
    Evidence: .sisyphus/evidence/task-8-diff-basic.txt

  Scenario: diff with --format json outputs valid JSON
    Tool: Bash
    Preconditions: Binary built, testdata files
    Steps:
      1. Run `./envdiff diff testdata/left.env testdata/right.env --format json`
      2. Pipe output to `python3 -m json.tool`
      3. Assert exit code 0 from python (valid JSON)
      4. Assert JSON contains "added", "removed", "changed" fields
    Expected Result: Valid, parseable JSON output
    Failure Indicators: JSON parse error, missing fields
    Evidence: .sisyphus/evidence/task-8-diff-json.txt

  Scenario: diff with nonexistent file exits with code 2
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./envdiff diff nonexistent.env also-missing.env 2>&1`
      2. Assert stderr contains error message about file not found
      3. Assert exit code is 2
    Expected Result: Error exit code 2, helpful error message
    Failure Indicators: Panic, exit code 1, no error message
    Evidence: .sisyphus/evidence/task-8-diff-missing-file.txt

  Scenario: identical files exit with code 0
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create testdata/same.env: "A=1\nB=2"
      2. Run `./envdiff diff testdata/same.env testdata/same.env`
      3. Assert exit code is 0
      4. Assert stdout is empty or shows "no differences" message
    Expected Result: Clean exit, no false positives
    Failure Indicators: Exit code 1, phantom diffs
    Evidence: .sisyphus/evidence/task-8-diff-identical.txt
  ```

  **Commit**: YES
  - Message: `feat: add diff command with global flags`
  - Files: `cmd/envdiff/diff.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 9. `check` Command

  **What to do**:
  - Implement `cmd/envdiff/check.go` — Cobra `check` subcommand:
    - `envdiff check` — auto-detect `.env` and `.env.example` in current directory, diff them
    - `envdiff check --source .env.local --example .env.example` — explicit paths
    - Auto-detection logic: look for `.env.example` or `.env.sample` or `.env.template`, compare against `.env`
    - Output: which keys are in `.env.example` but missing from `.env` (you need these!)
    - Output: which keys are in `.env` but not in `.env.example` (maybe add these to example?)
    - Reuse diff engine + formatter from Task 8
    - Same exit codes: 0 = in sync, 1 = out of sync, 2 = error

  **Must NOT do**:
  - Modify any files (check is read-only)
  - Show values from `.env` (these might be secrets) — auto-mask `.env` values, only show `.env.example` values

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Thin wrapper around existing diff engine with auto-detection logic
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 8, 10-13)
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 14
  - **Blocked By**: Task 6

  **References** (CRITICAL):

  **API/Type References**:
  - `cmd/envdiff/diff.go` — Reuse the diff wiring pattern from Task 8
  - `internal/diff/diff.go:Diff()` — Same diff engine

  **Acceptance Criteria**:

  - [ ] `./envdiff check` in a directory with `.env` and `.env.example` shows diff
  - [ ] Auto-detects `.env.example`, `.env.sample`, `.env.template`
  - [ ] `.env` values are auto-masked (never shown)
  - [ ] Exit code 0 when in sync, 1 when out of sync

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: check detects missing keys
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create temp dir with .env.example (A, B, C keys) and .env (A, B keys — missing C)
      2. Run `./envdiff check` in that directory
      3. Assert output mentions C as missing from .env
      4. Assert exit code 1
    Expected Result: Missing key identified
    Failure Indicators: Missing key not detected, wrong exit code
    Evidence: .sisyphus/evidence/task-9-check-missing.txt

  Scenario: check with no .env.example shows helpful error
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create temp dir with only .env (no example file)
      2. Run `./envdiff check 2>&1` in that directory
      3. Assert stderr mentions no example file found
      4. Assert exit code 2
    Expected Result: Clear error message, not a crash
    Failure Indicators: Panic, misleading error
    Evidence: .sisyphus/evidence/task-9-check-no-example.txt
  ```

  **Commit**: YES
  - Message: `feat: add check command`
  - Files: `cmd/envdiff/check.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 10. `compare` Command + Grouped Output

  **What to do**:
  - Implement `cmd/envdiff/compare.go` — Cobra `compare` subcommand:
    - `envdiff compare .env .env.staging .env.production` — matrix view
    - Uses `diff.MultiDiff()` to compute matrix
    - Table output: rows = keys (union of all files), columns = file names, cells = values (or `<missing>`)
    - JSON output: `{"keys": [...], "files": [...], "matrix": {...}}`
    - Grouped output by prefix: when displaying, group keys by common prefix separated by `_`
      - Example: `DB_HOST`, `DB_PORT`, `DB_NAME` grouped under `DB` header
      - `REDIS_URL`, `REDIS_PORT` grouped under `REDIS` header
      - Only group when 2+ keys share a prefix
    - `--ignore` flag works here too (filter keys before display)
    - `--mask` masks all values
    - Highlight cells where value differs from other environments (use lipgloss if TTY)

  **Must NOT do**:
  - Sort keys alphabetically (preserve original file ordering from first file)
  - Add sparklines or charts
  - Implement a scrollable TUI view (just print the table)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Matrix computation, grouped output logic, table alignment
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 8, 9, 11-13)
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 14
  - **Blocked By**: Tasks 6, 4

  **References** (CRITICAL):

  **API/Type References**:
  - `internal/diff/diff.go:MultiDiff()` — Matrix diff computation
  - `internal/output/table.go` — Table formatter (may need extending for matrix layout)

  **Acceptance Criteria**:

  - [ ] `./envdiff compare a.env b.env c.env` shows matrix table
  - [ ] Missing keys show `<missing>` in correct cells
  - [ ] `--format json` produces valid JSON matrix
  - [ ] Keys grouped by prefix when 2+ share prefix

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Compare 3 files shows matrix
    Tool: Bash
    Preconditions: Binary built, 3 test env files
    Steps:
      1. Create dev.env (A=1, B=2, C=3), staging.env (A=1, B=99), prod.env (A=1, C=3, D=4)
      2. Run `./envdiff compare dev.env staging.env prod.env`
      3. Assert output is a table with 3 columns
      4. Assert B row shows "2" for dev, "99" for staging, "<missing>" for prod
      5. Assert D row shows "<missing>" for dev, "<missing>" for staging, "4" for prod
    Expected Result: Correct matrix with all keys and all files represented
    Failure Indicators: Missing rows, wrong columns, crash on missing keys
    Evidence: .sisyphus/evidence/task-10-compare-matrix.txt

  Scenario: Grouped output clusters by prefix
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create env with DB_HOST, DB_PORT, DB_NAME, REDIS_URL, REDIS_PORT, APP_NAME
      2. Run compare on this file
      3. Assert DB_* keys appear together under "DB" group
      4. Assert REDIS_* keys appear together under "REDIS" group
    Expected Result: Keys grouped by prefix
    Failure Indicators: Flat list, wrong grouping
    Evidence: .sisyphus/evidence/task-10-grouped.txt
  ```

  **Commit**: YES
  - Message: `feat: add compare command with grouped output`
  - Files: `cmd/envdiff/compare.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 11. `validate` Command

  **What to do**:
  - Implement `cmd/envdiff/validate.go` — Cobra `validate` subcommand:
    - `envdiff validate --schema .env.schema.toml .env` — validate .env against schema
    - `envdiff validate --schema .env.schema.toml .env .env.staging .env.prod` — validate multiple files
    - Wire: parser → schema.LoadSchema → schema.Validate → format output
    - Output: list of validation errors with key name, expected type, actual value (masked if secret)
    - Pass output through formatter (table/json/github)
    - Exit codes: 0 = all valid, 1 = validation failures, 2 = error (schema not found, parse error)

  **Must NOT do**:
  - Modify the .env files
  - Auto-fix validation errors (just report)
  - Show secret values in error messages (mask them)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Wiring schema + parser + output, multi-file validation loop
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 8-10, 12-13)
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 14
  - **Blocked By**: Tasks 7, 4

  **References** (CRITICAL):

  **API/Type References**:
  - `internal/schema/schema.go:LoadSchema()`, `Validate()` — Schema validation functions
  - `internal/env/types.go:ValidationResult` — Output type with `IsValid()`, `Errors`, `Warnings`

  **Acceptance Criteria**:

  - [ ] `./envdiff validate --schema s.toml .env` exits 0 for valid file
  - [ ] Missing required key → exit 1 with error listing the key
  - [ ] Type mismatch → exit 1 with error showing expected type
  - [ ] Multiple files validated: each file's results shown separately

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Valid env passes validation
    Tool: Bash
    Preconditions: Binary built, schema and valid env created
    Steps:
      1. Create schema.toml requiring HOST(string) and PORT(number)
      2. Create valid.env with HOST=localhost, PORT=5432
      3. Run `./envdiff validate --schema schema.toml valid.env`
      4. Assert exit code 0
    Expected Result: Clean validation pass
    Failure Indicators: False validation errors
    Evidence: .sisyphus/evidence/task-11-validate-pass.txt

  Scenario: Invalid env fails with useful errors
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create schema.toml requiring API_KEY(string,required,pattern="^[A-Z0-9]{32}$")
      2. Create invalid.env with API_KEY=tooshort
      3. Run `./envdiff validate --schema schema.toml invalid.env 2>&1`
      4. Assert exit code 1
      5. Assert output mentions "API_KEY", "pattern", and the expected pattern
    Expected Result: Specific, actionable validation error
    Failure Indicators: Generic error, missing key name, missing pattern
    Evidence: .sisyphus/evidence/task-11-validate-fail.txt
  ```

  **Commit**: YES
  - Message: `feat: add validate command`
  - Files: `cmd/envdiff/validate.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 12. `init` Command

  **What to do**:
  - Implement `cmd/envdiff/init.go` — Cobra `init` subcommand:
    - `envdiff init` — generate `.env.example` from `.env` in current directory
      - Reads `.env`, strips all values, keeps keys + comments
      - Output: `.env.example` with `KEY=` (empty values) preserving original key order and comments
      - If `.env.example` already exists, print warning and exit (don't overwrite unless `--force`)
    - `envdiff init --schema` — generate `.env.schema.toml` from `.env`
      - Reads `.env`, infers types from values (number if integer, bool if true/false/yes/no, url if has scheme, string otherwise)
      - All keys marked as `required = true` by default
      - Output: `.env.schema.toml` with inferred schema
    - Both operations preserve comments from original `.env`

  **Must NOT do**:
  - Overwrite existing files without `--force` flag
  - Copy secret values into `.env.example` (this is the whole point — strip values)
  - Make complex type inference (simple heuristic is fine)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: File reading + templated writing, straightforward logic
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 8-11, 13)
  - **Parallel Group**: Wave 4
  - **Blocks**: None
  - **Blocked By**: Tasks 3, 7

  **References** (CRITICAL):

  **API/Type References**:
  - `internal/parser/parser.go:ParseFile()` — Read .env
  - `internal/env/types.go:EnvFile`, `EnvEntry` — Entry data with comments, ordering
  - `internal/schema/schema.go` — Schema types for TOML generation

  **Acceptance Criteria**:

  - [ ] `./envdiff init` creates `.env.example` with keys but no values
  - [ ] Comments from `.env` preserved in `.env.example`
  - [ ] `./envdiff init --schema` creates `.env.schema.toml` with inferred types
  - [ ] Existing files NOT overwritten without `--force`

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: init generates .env.example
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Create temp dir with .env containing "# DB config\nDB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret123"
      2. Run `./envdiff init` in that directory
      3. Assert .env.example exists
      4. Assert .env.example contains "DB_HOST=", "DB_PORT=", "API_KEY=" (no values)
      5. Assert .env.example contains "# DB config" (comment preserved)
    Expected Result: Clean example file with keys only
    Failure Indicators: Values leaked into example, comments lost
    Evidence: .sisyphus/evidence/task-12-init-example.txt

  Scenario: init refuses to overwrite existing file
    Tool: Bash
    Preconditions: Binary built, .env.example already exists
    Steps:
      1. Create temp dir with .env AND .env.example
      2. Run `./envdiff init 2>&1`
      3. Assert exit code non-zero
      4. Assert stderr mentions file already exists
    Expected Result: Safe refusal, no data loss
    Failure Indicators: Silent overwrite, no warning
    Evidence: .sisyphus/evidence/task-12-init-nooverwrite.txt
  ```

  **Commit**: YES
  - Message: `feat: add init command`
  - Files: `cmd/envdiff/init.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 13. `git` Command

  **What to do**:
  - Implement `internal/git/git.go`:
    - `ShowFileAtRef(ref, path string) (string, error)` — runs `git show <ref>:<path>`, returns file content
    - `IsGitRepo() bool` — check if current directory is a git repo
    - Uses `os/exec` to shell out to `git` (NOT go-git library)
  - Implement `cmd/envdiff/git.go` — Cobra `git` subcommand:
    - `envdiff git .env` — diff current `.env` against last committed version (`HEAD:.env`)
    - `envdiff git --ref main .env` — diff current `.env` against version on `main` branch
    - `envdiff git --ref abc123 --ref def456 .env` — diff `.env` between two commits
    - Wire: git.ShowFileAtRef → parser.Parse → diff.Diff → formatter
    - If file doesn't exist at ref, treat as empty env (all keys added/removed)
    - All global flags work: `--format`, `--mask`, `--ignore`

  **Must NOT do**:
  - Use `go-git` library (shell out to `git` with `os/exec`)
  - Handle merge commits specially (just read file at given ref)
  - Implement git log/blame features

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Git plumbing, edge cases with refs, file-not-found handling
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 8-12)
  - **Parallel Group**: Wave 4
  - **Blocks**: None
  - **Blocked By**: Task 6

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/glit-work/` — User's git tool project, uses os/exec for git commands

  **Acceptance Criteria**:

  - [ ] `./envdiff git .env` shows diff vs HEAD (in a git repo)
  - [ ] `./envdiff git --ref main .env` shows diff vs branch
  - [ ] Non-git-repo → exit code 2 with helpful error
  - [ ] File not in git → treated as empty (all keys shown as added)

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: git diff against HEAD
    Tool: Bash
    Preconditions: Binary built, in a git repo with committed .env
    Steps:
      1. Create temp git repo, commit .env with A=1, B=2
      2. Modify .env to A=1, B=99, C=3
      3. Run `./envdiff git .env`
      4. Assert: B shown as changed (2→99), C shown as added
      5. Assert exit code 1
    Expected Result: Correct diff against committed version
    Failure Indicators: Wrong values, git command failure
    Evidence: .sisyphus/evidence/task-13-git-head.txt

  Scenario: git in non-repo shows error
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./envdiff git .env` in a directory that is NOT a git repo
      2. Assert stderr contains error about not being a git repo
      3. Assert exit code 2
    Expected Result: Clear error, not panic
    Failure Indicators: Panic, wrong exit code
    Evidence: .sisyphus/evidence/task-13-git-non-repo.txt
  ```

  **Commit**: YES
  - Message: `feat: add git command`
  - Files: `cmd/envdiff/git.go, internal/git/git.go, internal/git/git_test.go`
  - Pre-commit: `go build ./cmd/envdiff && go test ./internal/git/... -v && go vet ./...`

- [ ] 14. CI Mode + Annotations

  **What to do**:
  - Add `--ci` flag to root command (persistent bool flag):
    - When `--ci` is set: force `--color never`, suppress all interactive behavior, use strict exit codes
    - Auto-detect CI environment: check `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL` env vars
    - When GitHub Actions detected + `--format` not set: auto-switch to `github` format
  - Implement `cmd/envdiff/ci.go` — Cobra `ci` subcommand:
    - `envdiff ci --require .env.example` — run `check` in CI mode, fail pipeline if out of sync
    - `envdiff ci --require .env.schema.toml` — run `validate` in CI mode
    - `envdiff ci --require .env.example --require .env.schema.toml` — run both checks
    - Output annotations in GitHub Actions format when detected
    - Summary line at end: `envdiff: 3 errors, 1 warning`

  **Must NOT do**:
  - Add interactive prompts in CI mode
  - Print colored output in CI mode (unless GitHub Actions, which supports colors)
  - Assume specific CI environment — detect generically

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Thin wiring of existing check/validate with CI-specific flags
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 15, 16)
  - **Parallel Group**: Wave 5
  - **Blocks**: Task 17
  - **Blocked By**: Tasks 8-11

  **References** (CRITICAL):

  **API/Type References**:
  - `cmd/envdiff/check.go` — Check command to reuse
  - `cmd/envdiff/validate.go` — Validate command to reuse
  - `internal/output/github.go` — GitHub annotations formatter

  **External References**:
  - GitHub Actions workflow commands: https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions

  **Acceptance Criteria**:

  - [ ] `CI=true ./envdiff diff a.env b.env` — no colors, correct exit codes
  - [ ] `GITHUB_ACTIONS=true ./envdiff diff a.env b.env` — GitHub annotation format
  - [ ] `./envdiff ci --require .env.example` — runs check in CI mode
  - [ ] `--ci` flag disables all interactive behavior

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: CI mode detects GitHub Actions
    Tool: Bash
    Preconditions: Binary built, test env files with differences
    Steps:
      1. Run `GITHUB_ACTIONS=true ./envdiff diff left.env right.env`
      2. Assert output contains `::warning` or `::error` annotations
      3. Assert no ANSI color codes in output
    Expected Result: GitHub Actions annotations, no colors
    Failure Indicators: ANSI codes present, wrong annotation format
    Evidence: .sisyphus/evidence/task-14-ci-github.txt

  Scenario: ci --require runs validation
    Tool: Bash
    Preconditions: Binary built, .env.example with missing keys
    Steps:
      1. Run `./envdiff ci --require .env.example`
      2. Assert exit code 1 (missing keys)
      3. Assert output contains summary line with error count
    Expected Result: Failed check with summary
    Failure Indicators: Exit code 0 on failure, no summary
    Evidence: .sisyphus/evidence/task-14-ci-require.txt
  ```

  **Commit**: YES
  - Message: `feat: add CI mode with annotations`
  - Files: `cmd/envdiff/ci.go`
  - Pre-commit: `go build ./cmd/envdiff && go vet ./...`

---

- [ ] 15. Pre-commit Hook Install/Uninstall

  **What to do**:
  - Implement `internal/hook/hook.go`:
    - `Install(repoRoot string) error` — write pre-commit hook script to `.git/hooks/pre-commit`
    - Hook script content: runs `envdiff check` (if `.env.example` exists) and `envdiff validate --schema .env.schema.toml` (if schema exists)
    - If `.git/hooks/pre-commit` already exists: append envdiff check to it (don't overwrite)
    - `Uninstall(repoRoot string) error` — remove envdiff section from pre-commit hook
    - Add comment markers to identify envdiff's section: `# BEGIN envdiff` / `# END envdiff`
  - Implement `cmd/envdiff/hook.go` — Cobra `hook` subcommand with sub-subcommands:
    - `envdiff hook install` — install pre-commit hook
    - `envdiff hook uninstall` — remove pre-commit hook
    - `envdiff hook status` — show if hook is installed

  **Must NOT do**:
  - Overwrite existing pre-commit hooks (append with markers)
  - Require a specific hook manager (husky, pre-commit framework)
  - Install hooks outside `.git/hooks/` directory

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: File writing with markers, simple logic
  - **Skills**: [`git-master`]
    - `git-master`: Git hooks directory structure knowledge

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 14, 16)
  - **Parallel Group**: Wave 5
  - **Blocks**: Task 17
  - **Blocked By**: Task 8

  **References** (CRITICAL):

  **External References**:
  - Git hooks docs: https://git-scm.com/docs/githooks#_pre_commit

  **Acceptance Criteria**:

  - [ ] `./envdiff hook install` creates executable `.git/hooks/pre-commit` with envdiff check
  - [ ] `./envdiff hook uninstall` removes envdiff section but preserves other hooks
  - [ ] `./envdiff hook status` reports installed/not installed
  - [ ] Existing hooks not destroyed on install

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: hook install creates pre-commit script
    Tool: Bash
    Preconditions: Binary built, in a git repo
    Steps:
      1. Run `./envdiff hook install`
      2. Assert .git/hooks/pre-commit exists
      3. Assert file is executable (chmod +x)
      4. Assert file contains "envdiff check"
      5. Assert file contains "# BEGIN envdiff" and "# END envdiff" markers
    Expected Result: Executable hook with envdiff section
    Failure Indicators: Non-executable, missing envdiff command, missing markers
    Evidence: .sisyphus/evidence/task-15-hook-install.txt

  Scenario: hook uninstall removes only envdiff section
    Tool: Bash
    Preconditions: Hook installed, other hook content present
    Steps:
      1. Add custom content before "# BEGIN envdiff" in pre-commit
      2. Run `./envdiff hook uninstall`
      3. Assert custom content still present
      4. Assert "# BEGIN envdiff" and "# END envdiff" removed
    Expected Result: Clean removal, other hooks preserved
    Failure Indicators: Entire file deleted, custom content lost
    Evidence: .sisyphus/evidence/task-15-hook-uninstall.txt
  ```

  **Commit**: YES
  - Message: `feat: add pre-commit hook management`
  - Files: `cmd/envdiff/hook.go, internal/hook/hook.go, internal/hook/hook_test.go`
  - Pre-commit: `go build ./cmd/envdiff && go test ./internal/hook/... -v && go vet ./...`

---

- [ ] 16. Interactive Sync Bubbletea TUI

  **What to do**:
  - Implement `internal/sync/model.go` — Bubbletea model for interactive sync:
    - Input: source `*env.EnvFile`, target `*env.EnvFile`, diff result
    - Display: list of diff entries as selectable items
      - Missing keys (in source, not in target): checkbox to add to target
      - Changed keys: show old→new, checkbox to update
      - Extra keys (in target, not in source): checkbox to remove
    - Navigation: j/k or arrows to move, space to toggle selection, enter to apply, q to cancel
    - On apply: write updated target file preserving original comments and ordering
    - Add selected keys at the end of file (or in appropriate group if grouped)
  - Implement `internal/sync/writer.go`:
    - `WriteEnvFile(path string, file *env.EnvFile) error` — write env file preserving formatting
    - Preserve: comment lines, blank lines, key ordering, quote style
    - New keys added at end of file
  - Implement `cmd/envdiff/sync.go` — Cobra `sync` subcommand:
    - `envdiff sync .env.example .env` — sync missing keys from example to local
    - `envdiff sync source.env target.env` — sync any two files
    - If stdout is not a TTY: print error "sync requires interactive terminal" and exit 2
  - Use Bubbletea v2 (`charm.land/bubbletea/v2`) and lipgloss v2 (`charm.land/lipgloss/v2`)
  - Style: match user's ports project TUI aesthetic (clean, minimal, keyboard-driven)

  **Must NOT do**:
  - Run in non-TTY mode (error out)
  - Auto-apply changes without user confirmation
  - Implement three-way merge (just source→target sync)
  - Use bubbletea v1 (use v2 — newer API)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: Bubbletea TUI with interactive elements, styling, keyboard navigation
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: TUI design, visual polish, interaction patterns

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 14, 15)
  - **Parallel Group**: Wave 5
  - **Blocks**: Task 17
  - **Blocked By**: Task 6

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/internal/tui/model.go` — User's Bubbletea model pattern (Init, Update, View)
  - `/Users/msh/code/pp/ports/internal/tui/keys.go` — Key binding definitions
  - `/Users/msh/code/pp/ports/internal/tui/styles.go` — Lipgloss style definitions

  **API/Type References**:
  - `internal/env/types.go:EnvFile, DiffResult` — Input types
  - `internal/diff/diff.go:Diff()` — Compute diff for display

  **External References**:
  - Bubbletea v2 examples: https://github.com/charmbracelet/bubbletea/tree/main/examples
  - Bubbles (list component): https://github.com/charmbracelet/bubbles — For selectable list

  **WHY Each Reference Matters**:
  - `model.go`: Shows EXACTLY how user structures Bubbletea models — must match this style
  - `keys.go`: Shows user's key binding pattern — j/k/arrows/q/enter conventions
  - `styles.go`: Shows user's lipgloss color and style preferences
  - Bubbletea v2: API differences from v1 (tea.KeyPressMsg vs tea.KeyMsg)

  **Acceptance Criteria**:

  - [ ] `./envdiff sync source.env target.env` opens TUI with selectable diff items
  - [ ] Space toggles selection, Enter applies, q cancels
  - [ ] Applied changes written to target file preserving formatting
  - [ ] Non-TTY mode: prints error and exits 2
  - [ ] `go test ./internal/sync/... -v` — writer tests pass

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Sync TUI launches and shows diff
    Tool: interactive_bash (tmux)
    Preconditions: Binary built, source.env and target.env created with known diffs
    Steps:
      1. Create source.env with A=1, B=2, C=3
      2. Create target.env with A=1, B=99
      3. Launch `./envdiff sync source.env target.env` in tmux
      4. Wait 2 seconds for TUI to render
      5. Capture tmux pane content
      6. Assert content shows C as missing (available to add)
      7. Assert content shows B as changed (2 vs 99)
      8. Send 'q' to exit
    Expected Result: TUI renders with correct diff items
    Failure Indicators: Blank screen, crash, wrong items
    Evidence: .sisyphus/evidence/task-16-sync-tui.txt

  Scenario: Sync in non-TTY prints error
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `echo "" | ./envdiff sync source.env target.env 2>&1`
      2. Assert stderr contains "requires interactive terminal"
      3. Assert exit code 2
    Expected Result: Clear error for non-interactive context
    Failure Indicators: TUI tries to launch, panic, wrong exit code
    Evidence: .sisyphus/evidence/task-16-sync-non-tty.txt

  Scenario: Writer preserves file formatting
    Tool: Bash
    Preconditions: Writer tests implemented
    Steps:
      1. Run `go test ./internal/sync/... -v -run TestWritePreservesFormatting`
      2. Assert: comments preserved, blank lines preserved, key ordering maintained
    Expected Result: Written file matches expected format
    Failure Indicators: Comments lost, ordering changed, extra blank lines
    Evidence: .sisyphus/evidence/task-16-writer-format.txt
  ```

  **Commit**: YES
  - Message: `feat: add interactive sync TUI`
  - Files: `cmd/envdiff/sync.go, internal/sync/model.go, internal/sync/keys.go, internal/sync/styles.go, internal/sync/writer.go, internal/sync/writer_test.go`
  - Pre-commit: `go build ./cmd/envdiff && go test ./internal/sync/... -v && go vet ./...`

---

- [ ] 17. GoReleaser + Homebrew + GitHub Actions

  **What to do**:
  - Update `.goreleaser.yaml` (created in Task 1) with final config:
    - Verify `main: ./cmd/envdiff`, `binary: envdiff`
    - Verify darwin+linux, amd64+arm64
    - Homebrew tap: `shahadulhaider/homebrew-tap`, formula with test `system "#{bin}/envdiff", "--version"`
  - Create `.github/workflows/release.yml`:
    - Trigger on tag push (`v*`)
    - Steps: checkout, setup Go, run tests (`go test ./...`), run GoReleaser
    - Uses `goreleaser/goreleaser-action`
  - Create `.github/workflows/ci.yml`:
    - Trigger on push/PR to main
    - Steps: checkout, setup Go, `go vet ./...`, `staticcheck ./...`, `go test ./... -v`, `go build ./cmd/envdiff`
    - Matrix: go 1.25, ubuntu-latest + macos-latest
  - Verify `make build && make vet && make lint` all pass
  - Verify `goreleaser check` passes

  **Must NOT do**:
  - Add Windows builds (darwin + linux only, matching ports)
  - Add Docker image builds
  - Push releases (just set up the automation)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Config files, copying patterns from existing project
  - **Skills**: [`git-master`]
    - `git-master`: GitHub Actions workflow setup

  **Parallelization**:
  - **Can Run In Parallel**: NO (needs all features complete)
  - **Parallel Group**: Wave 6 (solo)
  - **Blocks**: F1-F4
  - **Blocked By**: Tasks 14-16

  **References** (CRITICAL):

  **Pattern References**:
  - `/Users/msh/code/pp/ports/.goreleaser.yaml` — COPY THIS PATTERN exactly (change binary name, description, homepage)
  - `/Users/msh/code/pp/ports/.github/` — Check if CI workflows exist to copy pattern

  **WHY Each Reference Matters**:
  - GoReleaser config: User's proven config with correct Homebrew tap setup at `shahadulhaider/homebrew-tap`
  - The TAP_GITHUB_TOKEN env var pattern must be preserved

  **Acceptance Criteria**:

  - [ ] `goreleaser check` passes
  - [ ] `.github/workflows/release.yml` exists with tag trigger
  - [ ] `.github/workflows/ci.yml` exists with push/PR trigger
  - [ ] `make build && make vet` pass
  - [ ] All tests still pass: `go test ./... -v`

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: GoReleaser config is valid
    Tool: Bash
    Preconditions: GoReleaser installed
    Steps:
      1. Run `goreleaser check`
      2. Assert exit code 0
    Expected Result: Config passes validation
    Failure Indicators: YAML errors, missing fields
    Evidence: .sisyphus/evidence/task-17-goreleaser.txt

  Scenario: Full build + test pipeline passes
    Tool: Bash
    Preconditions: All tasks complete
    Steps:
      1. Run `make clean && make build`
      2. Run `go test ./... -v`
      3. Run `go vet ./...`
      4. Run `./envdiff --version`
      5. Assert all exit codes 0
      6. Assert version output is "envdiff dev"
    Expected Result: Clean build, all tests pass, binary works
    Failure Indicators: Build errors, test failures, vet warnings
    Evidence: .sisyphus/evidence/task-17-full-pipeline.txt
  ```

  **Commit**: YES
  - Message: `chore: add GoReleaser + Homebrew + CI workflows`
  - Files: `.goreleaser.yaml, .github/workflows/release.yml, .github/workflows/ci.yml`
  - Pre-commit: `goreleaser check && go test ./... && go vet ./...`

---

## Final Verification Wave (MANDATORY — after ALL implementation tasks)

> 4 review agents run in PARALLEL. ALL must APPROVE. Rejection → fix → re-run.

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, run command). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/. Compare deliverables against plan.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go vet ./...` + `staticcheck ./...` + `go test ./...`. Review all files for: `any` type assertions, empty catches, `fmt.Println` in prod code, commented-out code, unused imports. Check AI slop: excessive comments, over-abstraction, generic names (data/result/item/temp).
  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Tests [N pass/N fail] | Files [N clean/N issues] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high`
  Start from clean state. Execute EVERY QA scenario from EVERY task — follow exact steps, capture evidence. Test cross-command integration (diff → check → compare workflow). Test edge cases: empty .env, huge .env (1000+ keys), binary values, Unicode keys. Save to `.sisyphus/evidence/final-qa/`.
  Output: `Scenarios [N/N pass] | Integration [N/N] | Edge Cases [N tested] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  For each task: read "What to do", read actual code. Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT do" compliance. Detect cross-task contamination. Flag unaccounted changes.
  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

| Task | Commit Message | Key Files |
|------|---------------|-----------|
| 1 | `chore: scaffold envdiff project` | go.mod, Makefile, .goreleaser.yaml, cmd/envdiff/main.go |
| 2 | `feat: add core types and interfaces` | internal/env/types.go |
| 3 | `feat: implement .env file parser` | internal/parser/parser.go, parser_test.go |
| 4 | `feat: add output formatters (table/json/github)` | internal/output/*.go |
| 5 | `feat: add secret detection patterns` | internal/secret/secret.go, secret_test.go |
| 6 | `feat: implement diff engine` | internal/diff/diff.go, diff_test.go |
| 7 | `feat: add TOML schema validation` | internal/schema/schema.go, schema_test.go |
| 8 | `feat: add diff command with global flags` | cmd/envdiff/diff.go |
| 9 | `feat: add check command` | cmd/envdiff/check.go |
| 10 | `feat: add compare command with grouped output` | cmd/envdiff/compare.go |
| 11 | `feat: add validate command` | cmd/envdiff/validate.go |
| 12 | `feat: add init command` | cmd/envdiff/init.go |
| 13 | `feat: add git command` | cmd/envdiff/git.go, internal/git/git.go |
| 14 | `feat: add CI mode with annotations` | cmd/envdiff/ci.go |
| 15 | `feat: add pre-commit hook management` | cmd/envdiff/hook.go, internal/hook/hook.go |
| 16 | `feat: add interactive sync TUI` | cmd/envdiff/sync.go, internal/sync/*.go |
| 17 | `chore: add GoReleaser + Homebrew + CI` | .goreleaser.yaml, .github/workflows/release.yml |

---

## Success Criteria

### Verification Commands
```bash
go build ./cmd/envdiff                    # Expected: clean build, zero warnings
go test ./... -v                          # Expected: all tests PASS
go vet ./...                              # Expected: zero issues
staticcheck ./...                         # Expected: zero issues
./envdiff --version                       # Expected: "envdiff dev"
./envdiff diff testdata/a.env testdata/b.env  # Expected: diff output, exit 1
./envdiff check                           # Expected: .env vs .env.example comparison
./envdiff validate --schema s.toml .env   # Expected: validation output
./envdiff compare *.env                   # Expected: matrix view
echo $?                                   # Expected: 0 or 1 (never 2 for normal ops)
goreleaser check                          # Expected: config valid
```

### Final Checklist
- [ ] All 9 subcommands functional (diff, check, compare, validate, init, sync, git, hook, ci)
- [ ] All 3 output formats work (table, json, github)
- [ ] --mask hides values in all relevant commands
- [ ] --ignore skips keys by pattern in all relevant commands
- [ ] Exit codes are consistent (0/1/2)
- [ ] No interactive prompts except in `sync` command
- [ ] All tests pass
- [ ] GoReleaser config valid
- [ ] Binary runs on macOS (darwin/arm64)
