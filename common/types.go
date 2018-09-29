package common

import (
	"math/big"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// MBlock represents a block header in Mongodbq.
type MBlock struct {
	Number    int64  `bson:"number"`
	Hash      string `bson:"hash"`
	Timestamp int64  `bson:"timestamp"`
}

// MTransaction represents a transaction that will serialize to the RPC representation of a transaction
type MTransaction struct {
	BlockNumber int64  `bson:"blocknumber"`
	Hash        string `bson:"hash"`
	From        string `bson:"from"`
	Token       string `bson:"token"`
	To          string `bson:"to"`
	Value       string `bson:"value"`
	Timestamp   int64  `bson:"timestamp"`
}

// Block represents a block header in the Ethereum blockchain.
type Block struct {
	Difficulty      string        `json:"difficulty"`
	Extra           string        `json:"extraData"`
	GasLimit        string        `json:"gasLimit"`
	GasUsed         string        `json:"gasUsed"`
	Hash            string        `json:"hash"`
	Bloom           string        `json:"logsBloom"`
	Coinbase        string        `json:"miner"`
	MixDigest       string        `json:"mixHash"`
	Nonce           string        `json:"nonce"`
	Number          string        `json:"number"`
	ParentHash      string        `json:"parentHash"`
	ReceiptHash     string        `json:"receiptsRoot"`
	UncleHash       string        `json:"sha3Uncles"`
	Size            string        `json:"size"`
	Root            string        `json:"stateRoot"`
	Time            string        `json:"timestamp"`
	TotalDifficulty string        `json:"totalDifficulty"`
	TXs             []Transaction `json:"transactions"`
	TxHash          string        `json:"transactionsRoot"`
	Uncles          []string      `json:"uncles"`
}

// Transaction represents a transaction that will serialize to the RPC representation of a transaction
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func hexToDecimal(s string) bson.Decimal128 {
	bigInt := new(big.Int)
	bigInt.SetString(s, 0)
	bigIntByte, _ := bigInt.MarshalText()
	decimal, _ := bson.ParseDecimal128(string(bigIntByte))
	return decimal
	// return bigInt.Int64()
}

func hexToInt64(s string) int64 {
	bigInt := new(big.Int)
	bigInt.SetString(s, 0)
	return bigInt.Int64()
}

// ToMBlock 转为为mongodb需要的bson格式
func (r *Block) ToMBlock() *MBlock {
	var mb = MBlock{}

	mb.Number = hexToInt64(r.Number)
	mb.Hash = r.Hash
	mb.Timestamp = time.Now().Unix()

	return &mb
}

// ToMTransaction 转为为mongodb需要的bson格式
func (t *Transaction) ToMTransaction() *MTransaction {

	var mt = MTransaction{}

	mt.BlockNumber = hexToInt64(t.BlockNumber)
	mt.Token = t.To
	mt.From = t.From

	return &mt
}
