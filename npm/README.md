# debtbomb

[![npm version](https://img.shields.io/npm/v/debtbomb.svg)](https://www.npmjs.com/package/debtbomb)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An npm wrapper for [DebtBomb](https://github.com/jobin-404/debtbomb), a Go-based CLI tool to manage and track technical debt directly in your code using "bombs" that expire.

## Features

- **No Go required**: Automatically downloads the correct binary for your OS (macOS, Linux, Windows).
- **Per-project installation**: Keep your debt tracking versioned with your project.
- **CI/CD Ready**: Works seamlessly in GitHub Actions, Vercel, Netlify, etc.
- **Zero Configuration**: Just install and run.

## Installation

```bash
npm install --save-dev debtbomb
# or
yarn add -D debtbomb
# or
pnpm add -D debtbomb
```

## Usage

### Run via npx
```bash
npx debtbomb list
```

### Run via npm scripts
Add these to your `package.json`:
```json
{
  "scripts": {
    "debt:check": "debtbomb check",
    "debt:report": "debtbomb report"
  }
}
```

## Commands

- `debtbomb check`: Scan for expired debtbombs and exit with code 1 if found (perfect for CI).
- `debtbomb list`: List all debtbombs found in the project.
- `debtbomb report`: Show aggregated statistics about your technical debt.
- `debtbomb notify`: Send notifications (Jira/Slack) about expiring debt.

## How it works

This package is a thin wrapper. On installation, a `postinstall` script detects your operating system and architecture, then downloads the appropriate official release binary from GitHub. All arguments passed to the `debtbomb` command are transparently proxied to the Go binary.

## License

MIT
