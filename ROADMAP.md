# 💣 DebtBomb Roadmap

DebtBomb is already usable in real CI pipelines.
This roadmap shows how it evolves from a **local enforcement tool** into a **team-level technical-debt system**.

---

## Current — **v0.3.0**

DebtBomb today already provides a full enforcement loop:

**Core engine**

* Scan source files for `@debtbomb` comments
* Parse `expire`, `owner`, `ticket`, and `reason`
* Support both `:` and `=` syntax
* Support inline and multi-line formats
* Fail CI when expired bombs exist
* `debtbomb check` with `--warn-in-days`
* `debtbomb list` (expired, all, JSON)
* `debtbomb report` (owner, folder, urgency, etc)
* `.debtbombignore`
* Automatic exclusion of vendor, build, and binary files

**Team visibility**

* Jira ticket creation or update when bombs expire
* Slack notifications
* Discord notifications
* Microsoft Teams notifications

This already makes technical debt:
**visible, owned, and enforced.**

---

## v0.4 – Git awareness

Make debt traceable to people and commits.

* Auto-detect author and team via `git blame`
* Track when a bomb was introduced
* Show commit hash for expired bombs
* Optional `--no-git` mode for faster CI

This turns:

> “Some debt expired”
> into
> “This expired, in this commit, owned by this team.”

---

## v0.5 – TODO migration

Lower the barrier for existing messy codebases.

* Detect raw `TODO`, `FIXME`, `HACK`
* Use git history to infer how old they are
* Apply default expiry policy (e.g. 30 days)
* `debtbomb annotate` to convert TODOs into real DebtBombs

This lets teams adopt DebtBomb **without rewriting everything**.

---

## v0.6 – Configuration

Make it fit real teams.

* `.debtbombrc` or `debtbomb.toml`
* Default expiry window
* Warning thresholds
* Enable / disable Jira, Slack, etc
* Global ignore rules
* File size limits

---

## v0.7 – CI integrations

Make it drop-in for pipelines.

* Official GitHub Action
* GitLab CI template
* Jenkins example
* Pre-commit hook

---

## v0.8 – Distribution

Make it trivial to install everywhere.

* Homebrew formula
* npm wrapper
* Linux DEB / RPM
* Docker image

---

## v1.0 – Stable

DebtBomb is considered stable when:

* Config format is locked
* Output format is stable
* Jira + chat integrations are widely used
* Performance is proven on large monorepos

At that point it becomes safe for:
**company-wide technical debt enforcement.**

---

## Guiding principles

DebtBomb will always be:

* Fast
* Language-agnostic
* CI-first
* Opinionated

It is not a task manager.
It is a **time-bomb for technical debt**.
