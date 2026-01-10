# Contributing to DebtBomb

Thanks for your interest in contributing to DebtBomb 🧨
This project exists because real developers get burned by forgotten “temporary” code, and contributions are what will make it better.

DebtBomb is still early and evolving, so feedback, ideas, and small improvements are just as valuable as big features.

---

## What kind of contributions are welcome?

All of these help:

* Bug reports
* Feature ideas
* Documentation improvements
* Test cases
* Performance improvements
* New command ideas
* CI and packaging help (npm, Homebrew, RPM, etc.)

If you’re not sure where to start, check the Issues page — look for “good first issue” or “ideas”.

---

## Getting started

1. Fork the repo
2. Clone your fork
3. Create a branch

```bash
git checkout -b my-change
```

4. Make your change
5. Run the tool locally
6. Commit and open a pull request

---

## Running locally

```bash
go build -o debtbomb cmd/debtbomb/main.go
./debtbomb list
```

You can test against a fake repo by creating a few files with `@debtbomb` comments.

---

## Code style

This is a Go project. Please follow:

* `gofmt`
* Simple, readable code
* Prefer clarity over cleverness

DebtBomb is a CLI tool that needs to be:

* Fast
* Predictable
* Easy to debug

---

## Design principles

When contributing, try to keep these in mind:

* **Language-agnostic** — don’t tie logic to a specific programming language
* **CI-friendly** — output must be clean and scriptable
* **Fast** — large repos should still scan quickly
* **Opinionated** — temporary code should be visible and time-bounded

If a feature makes DebtBomb more complicated but not more useful in CI, it probably doesn’t belong.

---

## Opening issues

If you find a bug or have an idea, feel free to open an issue.

Good issues include:

* What you expected
* What happened
* Example comments or files
* Your OS and DebtBomb version

---

## Discussions

If you’re not sure something should be a feature, open a discussion instead of an issue. Ideas and design feedback are very welcome.

---

## Code of conduct

Be respectful.
Disagree on ideas, not people.
We’re all here to make messy codebases a little less painful.

---