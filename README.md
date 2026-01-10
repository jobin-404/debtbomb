# 🧨 DebtBomb

DebtBomb is a cross-language **technical-debt enforcement tool** that scans source code comments for time-limited “debt bombs” and fails CI when they expire.

It lets teams ship temporary hacks safely by attaching an expiry date to them.
When the date passes, the build fails — forcing the debt to be cleaned up instead of silently rotting forever.

![DebtBomb Logo](assets/logo.png)

---

## Why this exists

Every codebase has comments like:

```
TODO: remove later
FIXME: temporary workaround
```

They almost never get removed.

DebtBomb gives those comments a **deadline**.

Temporary code is allowed — but it must be **time-bounded, owned, and visible**.

---

## Installation

### Using Go

```bash
go install github.com/jobin-404/debtbomb/cmd/debtbomb@latest
```

If the `debtbomb` command is not found, add Go’s bin directory to your PATH:

**macOS / Linux**

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc
```

**Windows**
Add `%USERPROFILE%\go\bin` to your PATH environment variable.

---

### Build from source

```bash
git clone https://github.com/jobin-404/debtbomb.git
cd debtbomb
go build -o debtbomb cmd/debtbomb/main.go
```

---

## Usage

### Enforce in CI

```bash
debtbomb check
```

Fails with exit code `1` if any debt bomb is expired.

---

### Warning window

Warn before things explode:

```bash
debtbomb check --warn-in-days 7
```

This surfaces expiring debt in CI without blocking releases yet.

---

### Listing debt

```bash
debtbomb list
debtbomb list --expired
debtbomb list --json
```

---

## Syntax

DebtBomb looks for comments containing `@debtbomb`.
It works with any language because it only reads comments.

Supported comment styles:

* `//`
* `#`
* `--`
* `/* */`

---

### Single-line

```go
// @debtbomb(expire=2026-02-10, owner=pricing, ticket=JIRA-123)
```

---

### Multi-line

```go
// @debtbomb
//   expire: 2026-02-10
//   owner: pricing
//   ticket: JIRA-123
//   reason: Temporary surge override
```

---

### Fields

| Field    | Description                |
| -------- | -------------------------- |
| `expire` | **Required.** YYYY-MM-DD   |
| `owner`  | Team or person responsible |
| `ticket` | Issue tracker reference    |
| `reason` | Why this debt exists       |

---

## Ignoring files

Create a `.debtbombignore` file to exclude paths:

```
migrations/
legacy/
src/generated/*.go
```

---

## Automatic exclusions

DebtBomb skips files that are not human-written source.

**Directories**

* `node_modules`, `vendor`, `.venv`, `__pycache__`
* `dist`, `build`, `out`, `target`, `bin`, `pkg`, `obj`
* `.git`, `.svn`, `.hg`
* `.idea`, `.vscode`, `.terraform`

**Files**

* Images, videos, archives, executables
* PDFs and office documents
* Minified files (`.min.js`, `.min.css`)
* Lock files
* Any file larger than **1MB**

This keeps it fast even on large repos.

---

## License

MIT

---
