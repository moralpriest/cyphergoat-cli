# CypherGoat CLI

CypherGoat CLI is a cryptocurrency swap tool powered by [cyphergoat.com](https://cyphergoat.com) API.

## Features

- Interactive swap wizard
- Real-time exchange rate comparisons
- USD value display (calculated via CoinGecko API)
- Privacy coin support (XMR, ARRR, DERO, ZEC, and more)
- Context-aware API calls with timeout protection
- Built with Go 1.25 and modern best practices
- Source-available with SLSA provenance for security

## Security & Verification

For maximum security, we recommend building from source and verifying the binary:

### Why Build from Source?

- **Transparency**: See exactly what code runs
- **Trust**: No third-party build infrastructure
- **Verification**: SLSA provenance confirms binary matches source

### Verification with Cosign

```bash
# Install Cosign
curl -sL https://slsa.dev/install.sh | bash

# Verify SLSA provenance
cosign verify-attestation --type slsaprovenance --policy .github/policy.json cyphergoat
```

This verifies the binary was built from the correct source code using GitHub Actions.

## Installation

### Quick Install

```bash
# One-liner
curl -sL https://raw.githubusercontent.com/moralpriest/cyphergoat-cli/main/install.sh | sh
```

### Verification

Verify your installation using one of these methods:

#### Option 1: GitHub Attestation (Easy)

Requires [GitHub CLI](https://cli.github.com/):

```bash
# Install GitHub CLI if needed
brew install gh  # macOS
sudo apt install gh  # Linux

# Verify the binary
gh attestation verify /usr/local/bin/cyphergoat -R moralpriest/cyphergoat-cli
```

#### Option 2: Self-Signed Key (Cypherpunk)

Requires [Cosign](https://docs.sigstore.dev/cosign/overview):

```bash
# Install Cosign
brew install cosign  # macOS
sudo pacman -S cosign  # Arch Linux

# Download public key
curl -fsSL https://raw.githubusercontent.com/moralpriest/cyphergoat-cli/main/.github/cyphergoat.pub > cyphergoat.pub

# Verify binary
cosign verify --key cyphergoat.pub /usr/local/bin/cyphergoat
```

This provides cryptographic proof that the binary was signed by the project's private key.

### Build from Source (Recommended)

For maximum security, build from source:

#### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Task runner](https://taskfile.dev/installation/)

#### Steps

```bash
# 1. Clone the repository
git clone https://github.com/moralpriest/cyphergoat-cli.git
cd cyphergoat-cli

# 2. Install Task (if not already installed)
task install:task

# 3. Build the binary
task build

# 4. Verify (optional but recommended)
task verify

# 5. Install to system
sudo mv cyphergoat /usr/local/bin/

# Or install to user directory
task install:user
```

### Verify Installation

```bash
# Check version
cyphergoat version

# Should show:
# CypherGoat CLI v1
# Go: go1.25.x
```

## Usage

### Swap Command

Start the interactive swap wizard:

```bash
cyphergoat swap
```

With verbose debug output:

```bash
cyphergoat swap --verbose
cyphergoat swap -v
```

Example output:

```
CypherGoat Exchange

Send coin: btc
Send coin network (leave empty for default): btc
Receive coin: eth
Receive network (leave empty for default): eth
Amount to swap: 1

Fetching Rates from Partnered Exchanges...

Available Exchange Options

#  | Exchange     | You Receive    | Exchange Rate
1  | PegasusSwap  | 0.18522283 ETH | $601.25 USD
2  | ETZSwap      | 0.18509044 ETH | $601.18 USD
3  | ChangeNow    | 0.18450000 ETH | $599.62 USD

Select exchange option (enter number): 1
Your ETH receiving address: 0x...
```

### Version Command

Print version information:

```bash
cyphergoat version
```

Output:
```
CypherGoat CLI v1.0.0
Commit: abc12345
Date: 2025-01-06
Go: go1.25.5
```

## Privacy Coin Support

Fully supports privacy-focused cryptocurrencies:

| Coin | Ticker | Notes |
|------|--------|-------|
| Monero | XMR | Most established privacy coin |
| Pirate Chain | ARRR | Zero-knowledge protocol |
| Dero | DERO | Privacy-focused blockchain | (Pending)
| Zcash | ZEC | Zero-knowledge proofs |
| Wownero | WOW | Monero fork |
| Firo | FIRO | Lelantus protocol |
| Zano | Zano | Fast privacy transactions |
| Dash | DASH | PrivateSend feature |
| Beldex | BDX | Privacy-focused |
| Banano | BAN | Feeless, privacy-ready |

## Price Service

Exchange rates are calculated using real-time prices from CoinGecko API:

- No API key required (free tier)
- 5-minute price caching for performance
- Rate limiting (100ms between calls)
- Stablecoins (USDC, USDT, DAI) handled as 1:1 USD

## Configuration

### API Key

Set your CypherGoat API key using an environment variable:

```bash
export CYPHERGOAT_API_KEY="your_api_key_here"
```

Get your API key from [https://cyphergoat.com](https://cyphergoat.com).

### .env File

Create a `.env` file in the project root:

```
CYPHERGOAT_API_KEY=your_api_key_here
```

## Development

### Build Commands

```bash
task build          # Build the binary (outputs: cyphergoat)
task build:all      # Build for all platforms
task test           # Run tests with race detector
task test:short     # Run tests (no race detector)
task lint           # Run golangci-lint
task lint:install   # Install linter
task tidy           # Clean up go.mod/go.sum
task clean          # Remove built binary
task verify         # Verify SLSA provenance
task help           # Show all commands
```

### Running Locally

```bash
./cyphergoat swap
```

### Testing

```bash
# Run all tests
task test

# Run specific test
go test ./api/ -v
```

### Linting

```bash
task lint
```

### Vulnerability Check

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Requirements

- Go 1.25.5 or later
- Git
- Task runner

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

### Supply Chain Security

This project uses SLSA provenance to verify binary authenticity:

1. Every release includes SLSA attestation
2. Users can verify binary was built from source
3. No trust required in pre-built binaries

### Reporting Vulnerabilities

For security issues, please email security@your-domain.com or use GitHub's private vulnerability reporting.

## License

MIT

## Resources

- [Source Code](https://github.com/moralpriest/cyphergoat-cli)
- [SLSA Framework](https://slsa.dev)
- [Cosign](https://docs.sigstore.dev/cosign/overview)
- [CoinGecko API](https://docs.coingecko.com/)
