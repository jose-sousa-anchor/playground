# Ethereum Mempool Monitor

A Go program for monitoring the Ethereum mempool (transaction pool) in real-time.

## Features

- **Watch Mode**: Stream pending transactions in real-time with full details
- **Status Mode**: Check current pending and queued transaction counts
- **Content Mode**: View complete mempool contents
- **Inspect Mode**: Get text summary of mempool transactions

## Prerequisites

- Go 1.19 or higher
- Access to an Ethereum node with WebSocket support (QuickNode, Infura, Alchemy, or your own node)

## Installation

```bash
cd /home/user/playground/mempool
go mod tidy
```

## Configuration

### Option 1: Using .env file (Recommended)

Create a `.env` file in the project directory:

```bash
cp .env.example .env
```

Then edit `.env` and set your WebSocket URL:

```
ETH_WS_URL=wss://your-node-url
```

### Option 2: Using environment variable

Alternatively, set your Ethereum node WebSocket URL as an environment variable:

```bash
export ETH_WS_URL='wss://your-node-url'
```

### Node URL Examples:
- QuickNode: `wss://your-endpoint.quiknode.pro/YOUR-TOKEN/`
- Infura: `wss://mainnet.infura.io/ws/v3/YOUR-PROJECT-ID`
- Alchemy: `wss://eth-mainnet.g.alchemy.com/v2/YOUR-API-KEY`
- Local node: `ws://localhost:8546`

## Usage

### Watch Pending Transactions (Default)

Stream all pending transactions in real-time:

```bash
go run main.go watch
# or simply
go run main.go
```

Output shows:
- Transaction hash
- From/To addresses
- Value in ETH
- Gas limit and price
- Max fee and priority fee (EIP-1559)
- Nonce

### Get TxPool Status

Check how many transactions are pending and queued:

```bash
go run main.go status
```

### View TxPool Content

Get full details of all transactions in the mempool:

```bash
go run main.go content
```

### Inspect TxPool

Get a text summary of mempool transactions:

```bash
go run main.go inspect
```

## Example Output

```
[1] ================================
Hash:     0x1234...5678
From:     0xabcd...ef01
To:       0x9876...5432
Value:    0.100000 ETH
Gas:      21000
GasPrice: 25.00 Gwei
MaxFee:   30.00 Gwei
MaxTip:   2.00 Gwei
Nonce:    42
Pending:  true
```

## Use Cases

1. **MEV Opportunities**: Monitor large trades or arbitrage opportunities
2. **Gas Optimization**: Analyze current gas prices to optimize your transactions
3. **Security Monitoring**: Detect suspicious transactions targeting specific contracts
4. **Network Analysis**: Study transaction patterns and mempool behavior

## Building

Build a standalone binary:

```bash
go build -o mempool-monitor
./mempool-monitor watch
```

## Troubleshooting

If `txpool_*` methods fail, your node may not support them. These are optional Geth methods that some node providers disable. The `watch` mode uses standard `newPendingTransactions` subscription which is more widely supported.

## Notes

- The mempool is local to each node, so contents may vary between nodes
- High transaction volume can produce a lot of output
- Press Ctrl+C to stop monitoring
