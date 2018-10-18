package main

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/Eric-GreenComb/parse-eth/common"
	"github.com/Eric-GreenComb/parse-eth/config"
	"github.com/Eric-GreenComb/parse-eth/parser"
	"github.com/Eric-GreenComb/parse-eth/persist"
)

func main() {

	session, err := mgo.Dial(config.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	mongo := persist.Mongo{}
	mongo.SetCollection(session.DB(config.MongoDB.DB).C(config.MongoDB.Block), session.DB(config.MongoDB.DB).C(config.MongoDB.Token))

	_startNum := uint64(config.Ethereum.BlockNum)
	_mongoNum := mongo.GetSyncedBlockCount()

	if _startNum < _mongoNum {
		_startNum = _mongoNum
	}

	fmt.Println("start block num : ", _startNum)

	sync := make(chan int, 1)
	go mongo.Sync(_startNum, parser.GetLatestValidBlockNumber(), sync)

	// 周期同步
	for {
		select {
		case <-sync:
			log.Println("syncing task is completed.")
			time.Sleep(time.Duration(config.Server.Timer) * time.Second) // TODO: using event listen
			mongo.Sync(mongo.GetSyncedBlockCount(), parser.GetLatestValidBlockNumber(), sync)
		}
	}
}
