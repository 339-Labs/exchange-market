package config

type Config struct {
	HttpServerConfig HttpServerConfig `json:"http_server_config"`
	DBConfig         DBConfig         `json:"db_config"`
	ExchangeConfig   ExchangeConfig   `json:"exchange_config"`
}

type HttpServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DBConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type ExchangeConfig struct {
	Bn     CexExchangeConfig `json:"bn"`
	Okx    CexExchangeConfig `json:"okx"`
	ByBit  CexExchangeConfig `json:"bybit"`
	BitGet CexExchangeConfig `json:"bitget"`
	GateIo CexExchangeConfig `json:"gateio"`

	UniswapV2 DexExchangeConfig `json:"unswapV2"`
}

type CexExchangeConfig struct {
	ApiKey       string `json:"api_key"`
	ApiSecretKey string `json:"api_secret_key"`
	ApiUrl       string `json:"api_url"`
	WsUrl        string `json:"ws_url"`
	Passphrase   string `json:"passphrase"`
	TimeOut      int64  `json:"timeout"`
}

type DexExchangeConfig struct {
}
