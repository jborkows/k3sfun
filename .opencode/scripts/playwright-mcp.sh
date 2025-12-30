#!/bin/bash
BROWSER_PATH="$(which chromium || which google-chrome || which chromium-browser)"
exec npx @playwright/mcp@latest --executable-path "$BROWSER_PATH" "$@"
