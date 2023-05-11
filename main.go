package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

func main() {

	//以太坊节点
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum: %v", err)
	}
	log.Println("Success to connect to the Ethereum")

	// 获取最新的区块号
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the latest block number: %v", err)
	}
	log.Printf("The max blockNumber is %d", blockNumber)

	// 创建文件
	file, err := os.Create("transactions.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// EOA地址
	address := common.HexToAddress("8e5cC460F20916422Ab6223c93953454b48FF17e")

	// 遍历区块链
	for i := uint64(0); i <= blockNumber; i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			log.Fatal(err)
		}

		// 交易时间戳
		blockTime := block.Time()
		tm := time.Unix(int64(blockTime), 0)
		blockTimeStr := tm.Format("2006/01/02 15:04:05")

		// 遍历区块中的所有交易
		for _, tx := range block.Transactions() {
			if tx.Value().Sign() > 0 {
				// 获取发送方地址
				from, err := client.TransactionSender(context.Background(), tx, block.Hash(), 0)
				if err != nil {
					//log.Printf("Failed to get transaction sender: %v", err)
					continue
				}

				if strings.ToLower(tx.To().Hex()) == strings.ToLower(address.Hex()) || strings.ToLower(from.Hex()) == strings.ToLower(address.Hex()) {
					// 交易记录日志
					str := fmt.Sprintf("%s %s %s %s\n", blockTimeStr, from.Hex(), address.Hex(), tx.Value().String())
					log.Printf(str)

					// 写入文件
					_, err = file.WriteString(str)
					if err != nil {
						log.Printf("Failed to write transaction to file: %v", err)
						continue
					}
				}
			}
		}
	}
}
