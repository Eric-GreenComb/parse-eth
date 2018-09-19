package parser

import (
	"strconv"

	"github.com/Eric-GreenComb/parse-eth/common"
	"github.com/Eric-GreenComb/parse-eth/config"
)

// GetLatestBlockNumber GetLatestBlockNumber
func GetLatestBlockNumber() uint64 {
	block := common.Block{}
	latest, _ := Call(config.Ethereum.Host, "eth_getBlockByNumber", []interface{}{"latest", false})
	MapToObject(latest.Result, &block)
	latestBlock, _ := strconv.ParseUint(block.Number[2:], 16, 64)
	return latestBlock
}
