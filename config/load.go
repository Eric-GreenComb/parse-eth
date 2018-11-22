package config

import (
	"strings"

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
	Server.SendAPI = viper.GetBool("server.send_api")
	Server.Timer = viper.GetInt64("server.timer")

	Ethereum.Host = viper.GetString("ethereum.host")
	Ethereum.BlockNum = uint64(viper.GetInt64("ethereum.blocknum"))
	Ethereum.TokenAddress = strings.ToLower(viper.GetString("ethereum.token_addr"))
	Ethereum.ToAddress = strings.ToLower(viper.GetString("ethereum.to_addr"))
	Ethereum.ToAddressRemove0x = Remove0x(Ethereum.ToAddress)
	Ethereum.Token = viper.GetString("ethereum.token")

	MongoDB.Host = viper.GetString("mongo.host")
	MongoDB.DB = viper.GetString("mongo.db")
	MongoDB.Block = viper.GetString("mongo.block")
	MongoDB.Token = viper.GetString("mongo.token")
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
