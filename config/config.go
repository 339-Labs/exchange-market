package config

import (
	"github.com/339-Labs/exchange-market/flags"
	"github.com/urfave/cli/v2"
)

type Config struct {
	Migrations       string
	HttpServerConfig ServerConfig   `json:"http_server_config"`
	SlaveDBConfig    DBConfig       `json:"slave_db_config"`
	RedisConfig      RedisConfig    `json:"redis_config"`
	ExchangeConfig   ExchangeConfig `json:"exchange_config"`
}

type ServerConfig struct {
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

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Username string `json:"username"`
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
	WsUrlFeature string `json:"ws_url_feature"`
	Passphrase   string `json:"passphrase"`
	TimeOut      int64  `json:"timeout"`
}

type DexExchangeConfig struct {
	RpcUrl   string `json:"rpc_url"`
	WsRpcUrl string `json:"ws_rpc_url"`
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	return &Config{
		Migrations: ctx.String(flags.MigrationsFlag.Name),
		HttpServerConfig: ServerConfig{
			Host: ctx.String(flags.HttpServerHostFlag.Name),
			Port: ctx.Int(flags.HttpServerPortFlag.Name),
		},
		SlaveDBConfig: DBConfig{
			Host: ctx.String(flags.SlaveDbHostFlag.Name),
			Port: ctx.Int(flags.SlaveDbPortFlag.Name),
			User: ctx.String(flags.SlaveDbUserFlag.Name),
			Pass: ctx.String(flags.SlaveDbPasswordFlag.Name),
			Name: ctx.String(flags.SlaveDbNameFlag.Name),
		},
		RedisConfig: RedisConfig{
			Address:  ctx.String(flags.RedisAddressFlag.Name),
			Password: ctx.String(flags.RedisPasswordFlag.Name),
			Username: ctx.String(flags.RedisUserNameFlag.Name),
		},
		ExchangeConfig: ExchangeConfig{
			Bn: CexExchangeConfig{
				ApiKey:       ctx.String(flags.BnApiKeyFlag.Name),
				ApiSecretKey: ctx.String(flags.BnApiSecretKeyFlag.Name),
				ApiUrl:       ctx.String(flags.BnApiUrlFlag.Name),
				WsUrl:        ctx.String(flags.BnWsUrlFlag.Name),
				WsUrlFeature: ctx.String(flags.ByBitWsUrlFeature.Name),
				Passphrase:   ctx.String(flags.BnPassphrase.Name),
				TimeOut:      ctx.Int64(flags.BnTimeOut.Name),
			},
			Okx: CexExchangeConfig{
				ApiKey:       ctx.String(flags.OkxApiKeyFlag.Name),
				ApiSecretKey: ctx.String(flags.BnApiSecretKeyFlag.Name),
				ApiUrl:       ctx.String(flags.OkxApiUrlFlag.Name),
				WsUrl:        ctx.String(flags.OkxWsUrlFlag.Name),
				Passphrase:   ctx.String(flags.OkxPassphrase.Name),
				TimeOut:      ctx.Int64(flags.OkxTimeOut.Name),
			},
			ByBit: CexExchangeConfig{
				ApiKey:       ctx.String(flags.ByBitApiKeyFlag.Name),
				ApiSecretKey: ctx.String(flags.ByBitApiSecretKeyFlag.Name),
				ApiUrl:       ctx.String(flags.ByBitApiUrlFlag.Name),
				WsUrl:        ctx.String(flags.ByBitWsUrlFlag.Name),
				WsUrlFeature: ctx.String(flags.ByBitWsUrlFeature.Name),
				Passphrase:   ctx.String(flags.ByBitPassphrase.Name),
				TimeOut:      ctx.Int64(flags.ByBitTimeOut.Name),
			},
			BitGet: CexExchangeConfig{
				ApiKey:       ctx.String(flags.BitGetApiKeyFlag.Name),
				ApiSecretKey: ctx.String(flags.BitGetApiSecretKeyFlag.Name),
				ApiUrl:       ctx.String(flags.BitGetApiUrlFlag.Name),
				WsUrl:        ctx.String(flags.BitGetWsUrlFlag.Name),
				Passphrase:   ctx.String(flags.BitGetPassphrase.Name),
				TimeOut:      ctx.Int64(flags.BitGetTimeOut.Name),
			},
			GateIo: CexExchangeConfig{},
		},
	}, nil
}
