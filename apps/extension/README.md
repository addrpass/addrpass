# AddrPass Chrome Extension

Browser extension for AddrPass — autofill addresses from your vault or share tokenized links, directly from any website.

## Features

- **Form Detection**: Automatically detects address forms on any website using autocomplete attributes, field names, labels, and placeholder heuristics
- **Autofill**: Fill address fields with one click from your AddrPass vault
- **Address Picker**: Choose which address to fill when you have multiple
- **Share**: Create a tokenized share link directly from the popup
- **Copy**: Copy your address to clipboard
- **Multi-language Detection**: Recognizes address fields in English, Turkish, German, French, Spanish, and more

## How It Works

1. **Content Script** (`content/content.js`) — Injected into every page:
   - Scans for address-related input fields using multiple heuristics
   - Shows a floating "Fill with AddrPass" badge when a form is detected
   - Fills fields with proper event dispatching (works with React, Vue, Angular)
   - Handles both `<input>` and `<select>` elements (country dropdowns)

2. **Background Worker** (`background/background.js`) — Service worker:
   - Manages authentication (JWT storage)
   - Makes API calls to api.addrpass.com
   - Caches addresses locally for fast popup rendering
   - Coordinates between popup and content script

3. **Popup** (`popup/`) — Extension popup UI:
   - Login form
   - Address list with Fill / Copy / Share actions
   - Links to dashboard

## Form Detection Strategy

Fields are classified using 4 signal levels (in priority order):

1. **`autocomplete` attribute** — e.g., `autocomplete="street-address"` (strongest)
2. **`name`/`id` attributes** — e.g., `name="address_line_1"`, `id="shipping-city"`
3. **`placeholder` text** — e.g., `placeholder="Enter your city"`
4. **Associated `<label>` text** — e.g., `<label for="addr">Street Address</label>`

Detected field types: `line1`, `line2`, `city`, `state`, `postCode`, `country`, `phone`

A form is considered an address form when **2+ address fields** are detected.

## Install (Development)

1. Open Chrome → `chrome://extensions/`
2. Enable "Developer mode" (top right)
3. Click "Load unpacked"
4. Select the `apps/extension/` directory

## Build for Chrome Web Store

```bash
cd apps/extension
zip -r addrpass-extension.zip . -x ".*" "README.md"
```

Upload the zip at [Chrome Web Store Developer Dashboard](https://chrome.google.com/webstore/devconsole).

## Permissions

- `storage` — Store JWT token and cached addresses locally
- `activeTab` — Access the current tab to detect and fill forms
- Host permissions for `api.addrpass.com` and `addrpass.com` only
