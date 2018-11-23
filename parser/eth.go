package parser

import (
	"fmt"
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

// GetLatestValidBlockNumber GetLatestValidBlockNumber
func GetLatestValidBlockNumber() uint64 {
	var _blockNumber string
	latest, err := Call(config.Ethereum.Host, "eth_blockNumber", []interface{}{"latest", false})
	if err != nil {
		fmt.Println("call eth_blockNumber error : ", err.Error())
		return 0
	}
	MapToObject(latest.Result, &_blockNumber)
	_latestBlockNumber, _ := strconv.ParseUint(_blockNumber[2:], 16, 64)
	return _latestBlockNumber - 12
}
