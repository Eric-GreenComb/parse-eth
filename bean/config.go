package bean

// ServerConfig ServerConfig
type ServerConfig struct {
	RecharegeAPI string
	Timer        int64
}

// EthereumConfig EthereumConfig
type EthereumConfig struct {
	Host         string
	BlockNum     uint64
	TokenAddress string
	ToAddress    string
}

// MongoDBConfig MongoDBConfig
type MongoDBConfig struct {
	Host  string
	DB    string
	Block string
	Token string
}
