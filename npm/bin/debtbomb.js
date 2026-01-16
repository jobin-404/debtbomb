#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

/**
 * DebtBomb CLI Wrapper
 * Proxies arguments to the downloaded Go binary.
 */

function main() {
  const platform = process.platform;
  const binaryName = platform === 'win32' ? 'debtbomb.exe' : 'debtbomb-bin';
  const binaryPath = path.join(__dirname, binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error('DebtBomb binary not found. Please try reinstalling the package.');
    console.error(`Expected binary at: ${binaryPath}`);
    process.exit(1);
  }

  // Proxy all arguments to the binary
  const args = process.argv.slice(2);
  const child = spawn(binaryPath, args, {
    stdio: 'inherit',
    shell: false
  });

  child.on('error', (err) => {
    console.error('Failed to start DebtBomb binary:', err.message);
    process.exit(1);
  });

  child.on('exit', (code) => {
    process.exit(code === null ? 1 : code);
  });

  // Handle termination signals
  const signals = ['SIGINT', 'SIGTERM', 'SIGHUP'];
  signals.forEach((signal) => {
    process.on(signal, () => {
      if (!child.killed) {
        child.kill(signal);
      }
    });
  });
}

main();
