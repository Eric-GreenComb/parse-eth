package config

import (
	"github.com/spf13/viper"

	"github.com/Eric-GreenComb/parse-eth/bean"
)

// Server Server Config
var Server bean.ServerConfig

// Ethereum Ethereum Config
var Ethereum bean.EthereumConfig

// MongoDB MongoDB Config
var MongoDB bean.MongoDBConfig

func init() {
	readConfig()
	initConfig()
}

func readConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
}

func initConfig() {
	Server.RecharegeAPI = viper.GetString("server.recharege_api")
	Server.Timer = viper.GetInt64("server.timer")

	Ethereum.Host = viper.GetString("ethereum.host")
	Ethereum.BlockNum = uint64(viper.GetInt64("ethereum.blocknum"))
	Ethereum.TokenAddress = viper.GetString("ethereum.token_addr")
	Ethereum.ToAddress = viper.GetString("ethereum.to_addr")

	MongoDB.Host = viper.GetString("mongo.host")
	MongoDB.DB = viper.GetString("mongo.db")
	MongoDB.Block = viper.GetString("mongo.block")
	MongoDB.Token = viper.GetString("mongo.sinoc")
}
