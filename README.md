# TonOS DNS Daemon

A lightweight, authoritative DNS daemon written in Go for the **TonOS** ecosystem. It intercepts `.dapp` decentralized domain requests and routes them locally, while forwarding standard internet traffic to public resolvers.

## Features

- **Decentralized TLD Resolution:** Intercepts and resolves `.dapp` domains locally.
- **Authoritative Responses:** Acts as the primary nameserver for local TonOS services.
- **Fallback Mechanism:** Seamlessly proxies standard DNS queries (e.g., `.com`, `.org`) to Cloudflare (`1.1.1.1`).

## Installation & Running

Make sure you have Go installed, then run the daemon with superuser privileges (required to bind to port 53):

```bash
sudo go run main.go
```
