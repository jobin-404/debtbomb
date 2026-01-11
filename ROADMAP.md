# DebtBomb Roadmap

DebtBomb is already usable today.
The roadmap below focuses on turning it into a **solid, production-ready developer tool**.

---

## Current — v1.0.4

The foundation is in place:

* Scan source files for `@debtbomb` comments
* Parse expiry, owner, ticket, and reason
* Fail CI when expired bombs exist
* `debtbomb check` and `debtbomb list`
* `debtbomb report` for aggregated stats
* `--warn-in-days` support
* Human-readable output
* JSON output
* Ignore common vendor, build, and binary files
* `.debtbombignore` support

This is the **minimum viable enforcement engine**.

---

## v1.1 – Git awareness

Make DebtBomb understand where debt came from:

* Auto-fill owner using `git blame`
* Track when a bomb was added
* Show commit hash for expired bombs
* Optional `--no-git` for faster CI

This turns “expired code” into **owned, traceable debt**.

---

## v1.2 – TODO conversion

Lower the barrier to adoption:

* Detect raw `TODO`, `FIXME`, `HACK`
* Use git history to infer their age
* Compute expiry using a default policy (e.g. 30 days)
* `debtbomb annotate` to convert TODOs into DebtBombs

This lets messy codebases migrate gradually.

---

## v1.3 – Configuration

Make it fit real teams:

* `.debtbombrc` config file
* Default expiry days
* Warning window
* Git integration toggle
* Global ignore rules
* File size limits

---

## v1.4 – CI integrations

Make it drop-in for pipelines:

* GitHub Action
* GitLab CI template
* Jenkins example
* Pre-commit hook

---

## v1.5 – Distribution

Make it easy to install everywhere:

* npm wrapper
* Homebrew formula
* Linux RPM and DEB packages
* Docker image

---

## v2.0 – Stable

DebtBomb is considered stable when:

* Config format is locked
* Output format is stable
* CI integrations are widely used
* Performance is proven on large monorepos

At that point it becomes safe for **company-wide enforcement**.

---

## Guiding principles

DebtBomb will stay:

* Fast
* Language-agnostic
* CI-first
* Opinionated

It is not a task manager.
It is a **time-bomb for technical debt**.

---