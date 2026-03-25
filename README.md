# AddrPass

**Stop giving your address to every website.** Store it once, share tokenized links with access control, monitoring, and instant revocation.

[![npm](https://img.shields.io/npm/v/@addrpass/sdk)](https://www.npmjs.com/package/@addrpass/sdk)
[![License: AGPL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue)](LICENSE)

## What is AddrPass?

AddrPass is an open-source address management platform. Instead of copying your address into every website, you store it once and share **tokenized references** — links, QR codes, or short codes — with controlled access, real-time monitoring, and instant revocation.

**For individuals**: Protect your home address from data breaches, package theft, and unwanted exposure.

**For businesses**: Integrate via SDK or API. E-commerce sites store tokens instead of addresses. Delivery companies get QR codes instead of printed labels. Everyone reduces PII liability.

## Quick Start

### Use the cloud (free)

1. Go to [addrpass.com/register](https://addrpass.com/register)
2. Add your address
3. Create a share → get a link or QR code
4. Monitor who accesses it

### Self-host (free, unlimited)

```bash
git clone https://github.com/addrpass/addrpass.git
cd addrpass/docker
docker compose up
```

### Integrate via SDK

```bash
npm install @addrpass/sdk
```

**Backend** (Node.js):
```typescript
import { AddrPassClient } from "@addrpass/sdk";

const addrpass = new AddrPassClient({
  clientId: "ap_your_id",
  clientSecret: "aps_your_secret",
});

const { address } = await addrpass.resolve(shareToken);
```

**Frontend** (React):
```tsx
import { AddrPassProvider, AddrPassButton } from "@addrpass/sdk/react";

<AddrPassProvider clientId="ap_xxx" redirectUri="/callback" scope="delivery">
  <AddrPassButton onToken={(r) => ship(r.share_token)} />
</AddrPassProvider>
```

**Vanilla JS**:
```html
<script src="https://api.addrpass.com/widget.js"></script>
<script>
  new AddrPass({
    clientId: "ap_xxx",
    onToken: (data) => console.log(data.share_token)
  }).renderButton("#checkout");
</script>
```

## Features

- **Encrypted vault** — store multiple addresses securely
- **Tokenized sharing** — links, QR codes, short codes
- **Scoped access** — full, delivery (no phone), zone (city only), verify (exists only)
- **Access monitoring** — who, when, from where, which business
- **Expiration & limits** — auto-expire, max accesses, PIN protection
- **Instant revocation** — kill the token, address stays safe
- **Delegation chains** — user → e-commerce → delivery → driver
- **Shipping labels** — QR + zone code, no plaintext address
- **OAuth2** — authorization code flow for e-commerce checkout
- **Webhooks** — HMAC-signed access notifications
- **Rate limiting** — per-IP and per-API-key
- **Self-hostable** — Docker Compose, one command

## Architecture

```
addrpass.com (Cloudflare Workers) — Static frontend
api.addrpass.com (Contabo VPS, Germany) — Go API + PostgreSQL
```

- **API**: Go + chi + pgx (no ORM)
- **Frontend**: Next.js + Tailwind CSS
- **Database**: PostgreSQL 16
- **Auth**: JWT + bcrypt + OAuth2
- **Deployment**: Docker + Caddy (auto-SSL)

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full technical deep dive.

## Project Structure

```
apps/
  api/              Go API server (31 endpoints)
  web/              Next.js frontend (landing, dashboard, auth, consent)
packages/
  sdk-js/           @addrpass/sdk — npm package (MIT)
docker/
  docker-compose.yml
```

## API Reference

31 endpoints covering auth, addresses, shares, resolution, QR codes, businesses, API keys, OAuth2, webhooks, delegations, and shipping labels.

Full spec: [openapi.yaml](apps/api/openapi.yaml)

## SDK

| Package | Install | Docs |
|---------|---------|------|
| **@addrpass/sdk** | `npm install @addrpass/sdk` | [README](packages/sdk-js/README.md) |

The SDK provides:
- `AddrPassClient` — backend client with full API coverage
- `<AddrPassProvider>` + `<AddrPassButton>` — React components
- `useAddrPass()` — React hook for custom UI
- Vanilla JS widget via `widget.js`

## License

- **Core** (API, web app): [AGPL-3.0](LICENSE) — open source, copyleft
- **SDK** (npm package): [MIT](packages/sdk-js/README.md) — use freely in any project

## Links

- **Website**: [addrpass.com](https://addrpass.com)
- **API**: [api.addrpass.com](https://api.addrpass.com/health)
- **npm**: [@addrpass/sdk](https://www.npmjs.com/package/@addrpass/sdk)
- **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
- **Research**: [RESEARCH.md](RESEARCH.md)
