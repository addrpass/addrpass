# AddrPass

Your address, your control. Share it with a token, not a copy.

AddrPass is an open-source address management platform that replaces plaintext address sharing with tokenized, access-controlled links and QR codes. Store your addresses once, share references with full monitoring, expiration, and revocation.

## The Problem

Your home address is copied across dozens of services, printed on every package, and leaked in every data breach. You have zero visibility into who stores it, accesses it, or loses it.

## How It Works

1. **Store** your addresses in your encrypted vault
2. **Share** via tokenized links, QR codes, or short codes with access control
3. **Monitor** who accessed your address, when, and from what device
4. **Revoke** access anytime. Update your address once, all tokens resolve to the new one.

## Use Cases

- **Online shopping** -- share a token instead of typing your address into every site
- **Package delivery** -- QR code on label instead of printed plaintext address
- **Temporary sharing** -- expiring links for contractors, guests, one-time deliveries
- **Business integration** -- API for e-commerce, delivery companies, healthcare, finance
- **Address changes** -- update once, all active shares resolve to the new address

## Project Structure

```
addrpass/
  apps/
    api/          # Core API server
    web/          # Web dashboard
    docs/         # Documentation site
    scanner/      # Driver/scanner mobile app
  packages/
    sdk-js/       # JavaScript SDK
    sdk-python/   # Python SDK
    shared/       # Shared types and utilities
    qr/           # QR code generation library
  docker/         # Docker Compose for self-hosting
```

## Self-Hosting

```bash
git clone https://github.com/addrpass/addrpass.git
cd addrpass
docker compose up
```

## License

- Core (API, web app, shared): [AGPL-3.0](LICENSE)
- SDKs and client libraries: [MIT](packages/sdk-js/LICENSE)

## Links

- Website: [addrpass.com](https://addrpass.com)
- Documentation: [docs.addrpass.com](https://docs.addrpass.com)
