package main

import (
	"context"
	"flag"
	"fmt"
	"time"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type TxPoolStatus struct {
	Pending string `json:"pending"` // Number of pending transactions
	Queued  string `json:"queued"`  // Number of queued transactions
}

func main() {
	rpcUrl, _ := getParameters()

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	one := big.NewInt(1)

	currentBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get the start block: %v", err)
	}

	currentBlockHeight := currentBlock.Number()
	tryGetPool := true

	for ; ;  {
		currentBlock, err = client.BlockByNumber(context.Background(), currentBlockHeight)
		if err != nil {
			time.Sleep(800000 * time.Microsecond)
			continue
		}

		txCount := len(currentBlock.Transactions())

		var pendingTxCount uint64

		if tryGetPool {
			var status TxPoolStatus
			err = rpcClient.CallContext(context.Background(), &status, "txpool_status")
			if err != nil {
				tryGetPool = false
			}

			if tryGetPool {
				pendingTxCount, err = strconv.ParseUint(status.Pending[2:], 16, 64)
				if err != nil {
					log.Fatalf("Failed to parse hex string: %v", err)
				}
			}
		}

		fmt.Printf("Height: %v\tTxCount: %v\tPendingTx: %v\tBlockTime: %v\tGasLimit: %v\tGasUsed: %v\n",
			currentBlock.Number(),
			txCount,
			pendingTxCount,
			currentBlock.Time(),
			currentBlock.GasLimit(),
			currentBlock.GasUsed())

		currentBlockHeight.Add(currentBlockHeight, one)
		time.Sleep(200000 * time.Microsecond)
	}
}

func getParameters() (string, int) {
	// handle command line flags
	rpcUrl := flag.String("rpc-url", "http://127.0.0.1:8545", "RPC url of the chain")
	count := flag.Int("count", 10000, "The number of transactions to be sent")
	flag.Parse()

	if *count > 1000000 {
		log.Fatal("Too many transactions to be generated and sent")
	}

	return *rpcUrl, *count
}

