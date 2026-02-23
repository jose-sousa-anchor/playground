package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

// TxPoolStatus represents the status of the transaction pool
type TxPoolStatus struct {
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

// TxPoolContent represents the full content of the transaction pool
type TxPoolContent struct {
	Pending map[string]map[string]*types.Transaction `json:"pending"`
	Queued  map[string]map[string]*types.Transaction `json:"queued"`
}

func main() {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Read WebSocket URL from environment or .env file
	wsURL := os.Getenv("ETH_WS_URL")
	if wsURL == "" {
		fmt.Println("Please set ETH_WS_URL environment variable or create a .env file")
		fmt.Println("Example .env file:")
		fmt.Println("  ETH_WS_URL=wss://your-node-url")
		fmt.Println("\nOr export as environment variable:")
		fmt.Println("  export ETH_WS_URL='wss://your-node-url'")
		os.Exit(1)
	}

	// Choose mode from command line argument
	mode := "watch"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "watch":
		watchPendingTransactions(wsURL)
	case "status":
		getTxPoolStatus(wsURL)
	case "content":
		getTxPoolContent(wsURL)
	case "inspect":
		getTxPoolInspect(wsURL)
	default:
		fmt.Println("Available modes:")
		fmt.Println("  watch   - Watch pending transactions in real-time")
		fmt.Println("  status  - Get txpool status")
		fmt.Println("  content - Get full txpool content")
		fmt.Println("  inspect - Get txpool inspection")
		fmt.Println("\nUsage: go run main.go [mode]")
	}
}

// watchPendingTransactions subscribes to pending transactions and displays them
func watchPendingTransactions(wsURL string) {
	fmt.Println("Connecting to Ethereum node via WebSocket...")

	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected! Watching for pending transactions...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Subscribe to pending transactions
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads: %v", err)
	}
	defer sub.Unsubscribe()

	// Also subscribe to pending transactions using raw RPC
	rpcClient, err := rpc.Dial(wsURL)
	if err != nil {
		log.Fatalf("Failed to connect via RPC: %v", err)
	}
	defer rpcClient.Close()

	pendingTxs := make(chan common.Hash)
	pendingSub, err := rpcClient.EthSubscribe(ctx, pendingTxs, "newPendingTransactions")
	if err != nil {
		log.Fatalf("Failed to subscribe to pending transactions: %v", err)
	}
	defer pendingSub.Unsubscribe()

	count := 0
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case err := <-pendingSub.Err():
			log.Fatalf("Pending subscription error: %v", err)
		case txHash := <-pendingTxs:
			count++
			// Fetch transaction details
			tx, isPending, err := client.TransactionByHash(ctx, txHash)
			if err != nil {
				fmt.Printf("[%d] Hash: %s (details unavailable: %v)\n", count, txHash.Hex(), err)
				continue
			}

			// Display transaction info
			fmt.Printf("[%d] ================================\n", count)
			fmt.Printf("Hash:     %s\n", txHash.Hex())
			fmt.Printf("From:     %s\n", getTxFrom(tx))
			if tx.To() != nil {
				fmt.Printf("To:       %s\n", tx.To().Hex())
			} else {
				fmt.Printf("To:       Contract Creation\n")
			}
			fmt.Printf("Value:    %s ETH\n", weiToEther(tx.Value()))
			fmt.Printf("Gas:      %d\n", tx.Gas())
			fmt.Printf("GasPrice: %s Gwei\n", weiToGwei(tx.GasPrice()))
			if tx.GasFeeCap() != nil {
				fmt.Printf("MaxFee:   %s Gwei\n", weiToGwei(tx.GasFeeCap()))
			}
			if tx.GasTipCap() != nil {
				fmt.Printf("MaxTip:   %s Gwei\n", weiToGwei(tx.GasTipCap()))
			}
			fmt.Printf("Nonce:    %d\n", tx.Nonce())
			fmt.Printf("Pending:  %v\n", isPending)
			fmt.Println()

		case <-ctx.Done():
			fmt.Println("Stopped watching transactions")
			return
		}
	}
}

// getTxPoolStatus gets the current status of the transaction pool
func getTxPoolStatus(wsURL string) {
	rpcClient, err := rpc.Dial(wsURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer rpcClient.Close()

	var status TxPoolStatus
	err = rpcClient.Call(&status, "txpool_status")
	if err != nil {
		log.Fatalf("Failed to get txpool status: %v", err)
	}

	fmt.Println("Transaction Pool Status:")
	fmt.Printf("Pending: %s\n", status.Pending)
	fmt.Printf("Queued:  %s\n", status.Queued)
}

// getTxPoolContent gets the full content of the transaction pool
func getTxPoolContent(wsURL string) {
	rpcClient, err := rpc.Dial(wsURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer rpcClient.Close()

	var content map[string]interface{}
	err = rpcClient.Call(&content, "txpool_content")
	if err != nil {
		log.Fatalf("Failed to get txpool content: %v", err)
	}

	fmt.Println("Transaction Pool Content:")

	// Pretty print JSON
	jsonData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal content: %v", err)
	}

	fmt.Println(string(jsonData))
}

// getTxPoolInspect gets an inspection of the transaction pool
func getTxPoolInspect(wsURL string) {
	rpcClient, err := rpc.Dial(wsURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer rpcClient.Close()

	var inspect map[string]interface{}
	err = rpcClient.Call(&inspect, "txpool_inspect")
	if err != nil {
		log.Fatalf("Failed to get txpool inspection: %v", err)
	}

	fmt.Println("Transaction Pool Inspection:")

	// Pretty print JSON
	jsonData, err := json.MarshalIndent(inspect, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal inspection: %v", err)
	}

	fmt.Println(string(jsonData))
}

// Helper functions

func getTxFrom(tx *types.Transaction) string {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return "unknown"
	}
	return from.Hex()
}

func weiToEther(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	f := new(big.Float).SetInt(wei)
	f = f.Quo(f, big.NewFloat(1e18))
	return f.Text('f', 6)
}

func weiToGwei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	f := new(big.Float).SetInt(wei)
	f = f.Quo(f, big.NewFloat(1e9))
	return f.Text('f', 2)
}
