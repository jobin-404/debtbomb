# DebtBomb CLI Reference

This document provides a comprehensive reference for the `debtbomb` command-line interface, including commands, flags, exit codes, and usage examples.

## Overview

DebtBomb is a static analysis tool that scans source code for time-bound technical debt comments (`@debtbomb`). It is designed to be language-agnostic and easily integrable into CI/CD pipelines.

**Syntax:**
```bash
debtbomb <command> [flags]
```

## Commands

### `check`

The `check` command is the primary tool for enforcement. It scans the codebase and exits with a failure code if any debt bombs have expired.

**Usage:**
```bash
debtbomb check [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--warn-in-days` | `int` | `0` | If specified, reports items expiring within N days as warnings. Warnings do not cause a non-zero exit code unless they are already expired. |
| `--json` | `bool` | `false` | Outputs the check result in JSON format. Useful for parsing by other tools. |

**Exit Codes:**

| Code | Description |
|------|-------------|
| `0` | **Success.** No expired debt bombs found. Warnings (if any) are displayed but do not fail the build. |
| `1` | **Failure.** One or more debt bombs have expired, or a critical error occurred during scanning. |

**Use Cases:**

1.  **CI/CD Enforcement (GitHub Actions, GitLab CI, Jenkins):**
    Fail the build if any technical debt has expired.
    ```bash
    debtbomb check
    ```

2.  **Pre-emptive Warning:**
    Warn developers about debt that will expire in the next week, allowing them to address it before it breaks the build.
    ```bash
    debtbomb check --warn-in-days 7
    ```

3.  **JSON Output for Custom Tooling:**
    ```bash
    debtbomb check --json > scan_results.json
    ```

---

### `list`

The `list` command outputs a detailed list of all debt bombs found in the project. It is useful for auditing and manual review.

**Usage:**
```bash
debtbomb list [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--expired` | `bool` | `false` | Filters the output to show ONLY expired debt bombs. |
| `--json` | `bool` | `false` | Outputs the list in JSON format instead of a table. |

**Output (Table):**
Displays a formatted ASCII table with columns:
- **Expires**: Date and relative time remaining (e.g., `2025-12-31 (5d12h)`).
- **Owner**: The assignee of the debt.
- **Ticket**: Related issue tracker reference.
- **Location**: File path and line number.

**Use Cases:**

1.  **Developer Audit:**
    See all technical debt currently tracked in the system.
    ```bash
    debtbomb list
    ```

2.  **Find Expired Items:**
    Quickly locate items that need immediate attention.
    ```bash
    debtbomb list --expired
    ```

---

### `report`

The `report` command generates high-level statistics and metrics about the technical debt in the codebase. It helps engineering managers and leads understand the distribution and volume of debt.

**Usage:**
```bash
debtbomb report [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--json` | `bool` | `false` | Outputs the report in JSON format. |

**Report Sections:**
- **Debt by Owner**: Count of items assigned to specific users or teams.
- **Debt by Folder**: Distribution of debt across modules or directories.
- **Debt by Reason**: Common reasons for debt (if provided in comments).
- **By Urgency**: Breakdown of items by expiration status (Expired, < 30 days, < 90 days, > 90 days).
- **Extremes**: The oldest and newest debt items.

**Use Cases:**

1.  **Monthly Health Check:**
    Review which teams or modules are accumulating the most time-bound debt.
    ```bash
    debtbomb report
    ```

2.  **Dashboard Integration:**
    Export metrics to JSON for visualization in a dashboard (e.g., Grafana, Datadog).
    ```bash
    debtbomb report --json
    ```

---

## Comment Syntax Reference

DebtBomb scans for comments containing `@debtbomb`. It supports single-line and multi-line formats in any language that uses standard comment delimiters (`//`, `#`, `--`, `/* */`).

### Fields

| Field | Required | Format | Description |
|-------|----------|--------|-------------|
| `expire` | **Yes** | `YYYY-MM-DD` | The date when the code becomes invalid/expired. |
| `owner` | No | String | Person or team responsible (e.g., `user`, `@team`). |
| `ticket` | No | String | Issue tracker ID (e.g., `JIRA-123`, `#456`). |
| `reason` | No | String | Context on why the debt exists. |

### Examples

**Go / JS / Java / C++:**
```javascript
// @debtbomb(expire=2024-12-31, owner=frontend, ticket=UI-99)
// TODO: Refactor this component
```

**Python / Ruby / Shell:**
```python
# @debtbomb
#   expire: 2024-12-31
#   owner: backend
#   reason: Temporary hack for migration
```

**SQL / Lua:**
```sql
-- @debtbomb(expire=2025-01-01, ticket=DB-500)
```

---

## Ignore Configuration

To exclude specific files or directories from scanning, create a `.debtbombignore` file in the root of your repository. The syntax matches `.gitignore`.

**Example `.debtbombignore`:**
```text
# Ignore build artifacts
dist/
build/

# Ignore vendor directories
vendor/
node_modules/

# Ignore generated code
src/generated/
```

**Automatic Exclusions:**
DebtBomb automatically excludes common non-source directories (`.git`, `node_modules`, etc.) and binary files to ensure performance.
