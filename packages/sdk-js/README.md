# @addrpass/sdk

Official JavaScript/TypeScript SDK for [AddrPass](https://addrpass.com) — tokenized address sharing with access control, monitoring, and revocation.

## Install

```bash
npm install @addrpass/sdk
```

## Backend Usage (Node.js / Server)

### Resolve a shared address token

```typescript
import { AddrPassClient } from "@addrpass/sdk";

const client = new AddrPassClient();

// Resolve a public share token — no auth needed
const { address, scope } = await client.resolve("faM9_4ZyBEdv-RKInfgLdk3K");
console.log(address.line1, address.city, address.country);
```

### Authenticate with API keys

```typescript
import { AddrPassClient } from "@addrpass/sdk";

const client = new AddrPassClient({
  clientId: "ap_your_client_id",
  clientSecret: "aps_your_client_secret",
});

// Client auto-authenticates on first request
const { address } = await client.resolve("token_here");
```

### OAuth2 Authorization Code Flow (E-commerce)

```typescript
import { AddrPassClient } from "@addrpass/sdk";

const client = new AddrPassClient({
  clientId: "ap_your_client_id",
  clientSecret: "aps_your_client_secret",
});

// Step 1: Generate authorization URL for the customer
const authURL = client.getAuthorizationURL({
  redirectUri: "https://yourshop.com/callback",
  scope: "delivery",
  state: "order_12345",
});
// Redirect customer to authURL

// Step 2: Exchange the code (in your callback handler)
const { access_token, share_token, scope } = await client.exchangeCode(
  code,
  "https://yourshop.com/callback"
);

// Step 3: Resolve the address
const { address } = await client.resolve(share_token);
// Ship to: address.line1, address.city, ...
```

### Create a shipping label

```typescript
const label = await client.createLabel(share.id);
console.log(label.label.reference_code); // "AP-7X3K-9M2P"
console.log(label.label.zone_code);      // "DE-BER-101"
console.log(label.qr_code_url);          // QR code image URL

// Get label image URL
const imageURL = client.labelImageURL(label.label.reference_code);
```

### Manage delegations

```typescript
// Delegate a share to a delivery company
const delegation = await client.createDelegation({
  share_id: share.id,
  to_business_id: "delivery-company-uuid",
  scope: "delivery",
  note: "Package #1234",
});

// Revoke delegation after delivery
await client.revokeDelegation(delegation.id);
```

## Frontend Usage (React)

### Setup provider

```tsx
import { AddrPassProvider, AddrPassButton } from "@addrpass/sdk/react";

function App() {
  return (
    <AddrPassProvider
      clientId="ap_your_client_id"
      clientSecret="aps_your_client_secret"  // For server-side exchange
      redirectUri="https://yourshop.com/callback"
      scope="delivery"
    >
      <Checkout />
    </AddrPassProvider>
  );
}
```

### Checkout button

```tsx
import { AddrPassButton } from "@addrpass/sdk/react";

function Checkout() {
  return (
    <AddrPassButton
      onToken={(result) => {
        console.log("Share token:", result.share_token);
        // Send to your backend to resolve the address
      }}
      onError={(error) => console.error(error)}
    />
  );
}
```

### Custom UI with hook

```tsx
import { useAddrPass } from "@addrpass/sdk/react";

function CustomCheckout() {
  const { authorize, result, loading, error } = useAddrPass();

  if (result) {
    return <p>Address shared! Token: {result.share_token}</p>;
  }

  return (
    <div>
      <button onClick={authorize} disabled={loading}>
        {loading ? "Sharing..." : "Share your address"}
      </button>
      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
}
```

### Vanilla JS (no React)

```html
<script src="https://api.addrpass.com/widget.js"></script>
<div id="addrpass-btn"></div>
<script>
  new AddrPass({
    clientId: "ap_your_client_id",
    scope: "delivery",
    onToken: function(data) {
      console.log("Share token:", data.share_token);
    }
  }).renderButton("#addrpass-btn");
</script>
```

## API Reference

### `AddrPassClient`

| Method | Description |
|--------|-------------|
| `resolve(token, pin?)` | Resolve share token to address |
| `qrCodeURL(token)` | Get QR code image URL |
| `exchangeCode(code, redirectUri)` | Exchange OAuth code for tokens |
| `getAuthorizationURL(params)` | Generate OAuth consent URL |
| `listAddresses()` | List user's addresses |
| `createAddress(data)` | Create an address |
| `getAddress(id)` | Get address by ID |
| `updateAddress(id, data)` | Update address |
| `deleteAddress(id)` | Delete address |
| `createShare(data)` | Create a tokenized share |
| `listShares()` | List user's shares |
| `revokeShare(id)` | Revoke a share |
| `getAccessLogs(shareId)` | Get access logs for a share |
| `createLabel(shareId)` | Create a shipping label |
| `labelImageURL(refCode)` | Get label image URL |
| `createDelegation(data)` | Delegate share to business |
| `listDelegations(shareId)` | List delegations for share |
| `revokeDelegation(id)` | Revoke a delegation |

### React Components

| Export | Description |
|--------|-------------|
| `<AddrPassProvider>` | Context provider with OAuth config |
| `<AddrPassButton>` | Pre-styled "Share via AddrPass" button |
| `useAddrPass()` | Hook for custom UI (authorize, result, loading, error) |

## Self-Hosted

Point the SDK to your own instance:

```typescript
const client = new AddrPassClient({
  baseURL: "https://your-addrpass.example.com",
});
```

```tsx
<AddrPassProvider
  clientId="ap_xxx"
  redirectUri="..."
  appURL="https://your-addrpass.example.com"
  apiURL="https://api.your-addrpass.example.com"
>
```

## License

MIT — use freely in any project.
