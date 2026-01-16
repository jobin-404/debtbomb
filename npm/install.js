const fs = require('fs');
const path = require('path');
const https = require('https');

const pkg = require('./package.json');
const VERSION = pkg.version;
const REPO = 'jobin-404/debtbomb';

function getBinaryName() {
  const platform = process.platform;
  const arch = process.arch;

  let os = '';
  let architecture = '';

  switch (platform) {
    case 'darwin':
      os = 'darwin';
      break;
    case 'linux':
      os = 'linux';
      break;
    case 'win32':
      os = 'windows';
      break;
    default:
      console.error(`Unsupported platform: ${platform}`);
      process.exit(1);
  }

  switch (arch) {
    case 'x64':
      architecture = 'amd64';
      break;
    case 'arm64':
      architecture = 'arm64';
      break;
    default:
      console.error(`Unsupported architecture: ${arch}`);
      process.exit(1);
  }

  const extension = os === 'windows' ? '.exe' : '';
  return `debtbomb-${os}-${architecture}${extension}`;
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    
    const request = (downloadUrl) => {
      https.get(downloadUrl, (response) => {
        if (response.statusCode === 302 || response.statusCode === 301) {
          request(response.headers.location);
          return;
        }

        if (response.statusCode !== 200) {
          reject(new Error(`Failed to download binary: ${response.statusCode} ${response.statusMessage}`));
          return;
        }

        response.pipe(file);
        file.on('finish', () => {
          file.close(resolve);
        });
      }).on('error', (err) => {
        fs.unlink(dest, () => reject(err));
      });
    };

    request(url);
  });
}

async function install() {
  const binaryName = getBinaryName();
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;
  const binDir = path.join(__dirname, 'bin');
  const dest = path.join(binDir, process.platform === 'win32' ? 'debtbomb.exe' : 'debtbomb-bin');

  console.log(`Downloading DebtBomb binary from ${url}...`);

  try {
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }

    await download(url, dest);

    if (process.platform !== 'win32') {
      fs.chmodSync(dest, 0o755);
    }

    console.log('DebtBomb binary downloaded successfully.');
  } catch (err) {
    console.error('Error downloading DebtBomb binary:', err.message);
    console.error('Please ensure you have an active internet connection and that the release exists.');
    process.exit(1);
  }
}

install();