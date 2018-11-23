package persist

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Eric-GreenComb/contrib/net/http"
	"github.com/Eric-GreenComb/parse-eth/bean"
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
			log.Println(err.Error())
		}

		if err := parser.MapToObject(resp.Result, &block); err != nil {
			log.Println(err.Error())
		}

		mBlock := block.ToMBlock()

		if err := m.InsertBlockInfo(mBlock); err != nil {
			log.Println(err.Error())
		}

		log.Println("block : ", i, block.Number, len(block.TXs))

		for _, _tx := range block.TXs {
			// if tracker eth
			if config.Ethereum.TokenAddress == "0x" {
				m.TrackEth(_tx)
			} else {
				m.TrackToken(_tx)
			}
		}
	}

	c <- 1
}

// TrackToken TrackToken
func (m *Mongo) TrackToken(_tx common.Transaction) error {

	if strings.ToLower(_tx.To) != config.Ethereum.TokenAddress {
		return nil
	}

	_addr, _value, err := parser.ParseTokenTransfer(_tx.Input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if strings.ToLower(_addr) != config.Ethereum.ToAddressRemove0x {
		return nil
	}
	log.Println("======== addr==ToAddressRemove0x : ", _addr)

	mTransaction := _tx.ToMTransaction()
	mTransaction.Hash = _tx.Hash
	mTransaction.To = _addr
	mTransaction.Value = _value
	mTransaction.Timestamp = time.Now().Unix()

	if err := m.InsertTokenTransfer(mTransaction); err != nil {
		log.Println(err.Error())
		return err
	}

	if !config.Server.SendAPI {
		return nil
	}

	var _reCharge bean.ReCharge
	_reCharge.Address = _tx.From
	_reCharge.Nums = _value
	_reCharge.CreateTime = time.Now().Unix()
	_reCharge.IsAuth = config.Ethereum.Token

	_postJSON, _ := json.Marshal(_reCharge)

	http.PostJSONString(config.Server.RecharegeAPI, string(_postJSON))

	log.Println(_tx.BlockNumber, _tx.From, _addr, _value)

	return nil
}

// TrackEth TrackEth
func (m *Mongo) TrackEth(_tx common.Transaction) error {

	if strings.ToLower(_tx.To) != config.Ethereum.ToAddress {
		return nil
	}

	log.Println("======== addr==ToAddress : ", _tx.To)

	mTransaction := _tx.ToMTransaction()
	mTransaction.Hash = _tx.Hash
	mTransaction.To = _tx.To

	_value := new(big.Int)
	_value, _ = _value.SetString(Remove0x(_tx.Value), 16)

	mTransaction.Value = _value.String()
	mTransaction.Timestamp = time.Now().Unix()
	_input, _ := hex.DecodeString(Remove0x(_tx.Input))
	mTransaction.Input = string(_input)

	fmt.Println(mTransaction)

	if err := m.InsertTokenTransfer(mTransaction); err != nil {
		log.Println(err.Error())
		return err
	}

	if !config.Server.SendAPI {
		return nil
	}

	var _reCharge bean.ReCharge
	_reCharge.Address = _tx.From
	_reCharge.Nums = mTransaction.Value
	_reCharge.CreateTime = time.Now().Unix()
	_reCharge.IsAuth = config.Ethereum.Token

	_postJSON, _ := json.Marshal(_reCharge)

	http.PostJSONString(config.Server.RecharegeAPI, string(_postJSON))

	log.Println(_tx.BlockNumber, _tx.From, _tx.To, _value.String())

	return nil
}

// Remove0x Remove0x
func Remove0x(s string) string {
	var _ret string
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			_ret = s[2:]
		}
	}
	return _ret
}
