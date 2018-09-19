package persist

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Eric-GreenComb/parse-eth/common"
	"github.com/Eric-GreenComb/parse-eth/config"
	"github.com/Eric-GreenComb/parse-eth/parser"
)

// Mongo Mongo
type Mongo struct {
	Block *mgo.Collection
	Token *mgo.Collection
}

// SetCollection SetCollection
func (m *Mongo) SetCollection(block *mgo.Collection, token *mgo.Collection) *Mongo {
	m.Block = block
	m.Token = token
	return m
}

// InsertBlockInfo InsertBlockInfo
func (m *Mongo) InsertBlockInfo(block interface{}) error {
	if err := m.Block.Insert(block); err != nil {
		return err
	}
	return nil
}

// GetSyncedBlockCount GetSyncedBlockCount
func (m *Mongo) GetSyncedBlockCount() uint64 {
	result := common.MBlock{}
	m.Block.Find(bson.M{}).Sort("-number").Limit(1).One(&result)
	return uint64(result.Number)
}

// InsertTokenTransfer InsertTokenTransfer
func (m *Mongo) InsertTokenTransfer(tokenTransfer interface{}) error {
	if err := m.Token.Insert(tokenTransfer); err != nil {
		return err
	}
	return nil
}

// Sync Sync
func (m *Mongo) Sync(syncedNumber, latestBlock uint64, c chan int) {
	block := common.Block{}
	if syncedNumber > 0 {
		// 从下一个块开始同步
		syncedNumber++
	}

	for i := syncedNumber; i <= latestBlock; i++ {

		number := fmt.Sprintf("0x%s", strconv.FormatUint(uint64(i), 16))
		resp, err := parser.Call(config.Ethereum.Host, "eth_getBlockByNumber", []interface{}{number, true})
		if err != nil {
			log.Fatal(err)
		}

		if err := parser.MapToObject(resp.Result, &block); err != nil {
			log.Fatalln(err)
		}

		mBlock := block.ToMBlock()

		if err := m.InsertBlockInfo(mBlock); err != nil {
			log.Fatal(err)
		}

		fmt.Println("block : ", i, block.Number, len(block.TXs))

		for _, _tx := range block.TXs {

			if _tx.To == config.Ethereum.TokenAddress {

				_addr, _value, err := parser.ParseTokenTransfer(_tx.Input)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				if _addr != config.Ethereum.ToAddress {
					continue
				}

				mTransaction := _tx.ToMTransaction()
				mTransaction.To = _addr
				mTransaction.Value = _value
				mTransaction.Timestamp = time.Now().Unix()

				if err := m.InsertTokenTransfer(mTransaction); err != nil {
					log.Fatal(err)
				}
				fmt.Println(_tx.BlockNumber, _tx.From, _addr, _value)
			}
		}
	}

	c <- 1
}
