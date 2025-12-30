# OpenCode Configuration

## Prerequisites

### Playwright MCP Server

The Playwright MCP server requires a Chromium-based browser to be installed on your system.

#### Required packages

Install one of the following browsers:

**Arch Linux / Manjaro / Tuxedo (pacman):**
```bash
sudo pacman -S chromium
```

**Debian / Ubuntu (apt):**
```bash
sudo apt install chromium-browser
# or
sudo apt install google-chrome-stable
```

**Fedora (dnf):**
```bash
sudo dnf install chromium
```

#### Node.js dependencies

Playwright MCP requires `npx` to be available:

```bash
# Ensure Node.js and npm are installed
node --version
npm --version
```

### Verification

After installation, verify the browser is accessible:

```bash
which chromium || which google-chrome || which chromium-browser
```

This should return a path like `/usr/bin/chromium`.

## Configuration

The `opencode.json` configures MCP servers. The Playwright server uses a wrapper script (`scripts/playwright-mcp.sh`) that dynamically locates the system browser.

## Troubleshooting

If Playwright fails to launch:

1. Ensure a supported browser is installed (see above)
2. Check the browser is in your PATH: `which chromium`
3. Restart OpenCode after making configuration changes
