package flags

import "github.com/urfave/cli/v2"

const envVarPrefix = "MARKET"

func prefixEnvVars(name string) []string {
	return []string{envVarPrefix + "_" + name}
}

var (
	MigrationsFlag = &cli.StringFlag{
		Name:    "migrations-dir",
		Value:   "./migrations",
		Usage:   "path for database migrations",
		EnvVars: prefixEnvVars("MIGRATIONS_DIR"),
	}

	// http service
	HttpServerHostFlag = &cli.StringFlag{
		Name:    "http-host",
		Usage:   "http server host",
		EnvVars: prefixEnvVars("HTTP_HOST"),
	}
	HttpServerPortFlag = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "http server port",
		EnvVars: prefixEnvVars("HTTP_PORT"),
	}

	// Slave DB  flags
	SlaveDbHostFlag = &cli.StringFlag{
		Name:     "slave-db-host",
		Usage:    "The host of the slave database",
		EnvVars:  prefixEnvVars("SLAVE_DB_HOST"),
		Required: true,
	}
	SlaveDbPortFlag = &cli.IntFlag{
		Name:     "slave-db-port",
		Usage:    "The port of the slave database",
		EnvVars:  prefixEnvVars("SLAVE_DB_PORT"),
		Required: true,
	}
	SlaveDbUserFlag = &cli.StringFlag{
		Name:     "slave-db-user",
		Usage:    "The user of the slave database",
		EnvVars:  prefixEnvVars("SLAVE_DB_USER"),
		Required: true,
	}
	SlaveDbPasswordFlag = &cli.StringFlag{
		Name:     "slave-db-password",
		Usage:    "The host of the slave database",
		EnvVars:  prefixEnvVars("SLAVE_DB_PASSWORD"),
		Required: true,
	}
	SlaveDbNameFlag = &cli.StringFlag{
		Name:     "slave-db-name",
		Usage:    "The db name of the slave database",
		EnvVars:  prefixEnvVars("SLAVE_DB_NAME"),
		Required: true,
	}

	// redis flags
	RedisAddressFlag = &cli.StringFlag{
		Name:    "redis-address",
		Usage:   "The address of the redis",
		EnvVars: prefixEnvVars("REDIS_ADDRESS"),
	}
	RedisPasswordFlag = &cli.StringFlag{
		Name:    "redis-password",
		Usage:   "The password of the redis",
		EnvVars: prefixEnvVars("REDIS_PASSWORD"),
	}
	RedisUserNameFlag = &cli.StringFlag{
		Name:    "redis-user-name",
		Usage:   "The username of the redis",
		EnvVars: prefixEnvVars("REDIS_USER_NAME"),
	}

	// bn flags
	BnApiKeyFlag = &cli.StringFlag{
		Name:    "bn-api-key",
		Usage:   "The apikey of the bn",
		EnvVars: prefixEnvVars("BN_API_KEY"),
	}
	BnApiSecretKeyFlag = &cli.StringFlag{
		Name:    "bn-api-secret-key",
		Usage:   "The secret of the bn",
		EnvVars: prefixEnvVars("BN_API_SECRET_KEY"),
	}
	BnApiUrlFlag = &cli.StringFlag{
		Name:    "bn-api-url",
		Usage:   "The api url of the bn",
		EnvVars: prefixEnvVars("BN_API_URL"),
	}
	BnWsUrlFlag = &cli.StringFlag{
		Name:     "bn-ws-url",
		Usage:    "The ws url of the bn",
		EnvVars:  prefixEnvVars("BN_WS_URL"),
		Required: true,
	}
	BnWsUrlFeature = &cli.StringFlag{
		Name:    "bn-ws-url-feature",
		Usage:   "The ws url of the bn",
		EnvVars: prefixEnvVars("BN_WS_URL_FEATURE"),
	}
	BnPassphrase = &cli.StringFlag{
		Name:    "bn-passphrase",
		Usage:   "The passphrase of the bn",
		EnvVars: prefixEnvVars("BN_PASSPHRASE"),
	}
	BnTimeOut = &cli.IntFlag{
		Name:     "bn-timeout",
		Usage:    "The timeout of the bn",
		EnvVars:  prefixEnvVars("BN_TIMEOUT"),
		Required: true,
	}

	// okx flags
	OkxApiKeyFlag = &cli.StringFlag{
		Name:    "okx-api-key",
		Usage:   "The apikey of the okx",
		EnvVars: prefixEnvVars("OKX_API_KEY"),
	}
	OkxApiSecretKeyFlag = &cli.StringFlag{
		Name:    "okx-api-secret-key",
		Usage:   "The secret of the okx",
		EnvVars: prefixEnvVars("OKX_API_SECRET_KEY"),
	}
	OkxApiUrlFlag = &cli.StringFlag{
		Name:    "okx-api-url",
		Usage:   "The api url of the okx",
		EnvVars: prefixEnvVars("OKX_API_URL"),
	}
	OkxWsUrlFlag = &cli.StringFlag{
		Name:     "okx-ws-url",
		Usage:    "The ws url of the okx",
		EnvVars:  prefixEnvVars("OKX_WS_URL"),
		Required: true,
	}
	OkxPassphrase = &cli.StringFlag{
		Name:    "okx-passphrase",
		Usage:   "The passphrase of the okx",
		EnvVars: prefixEnvVars("OKX_PASSPHRASE"),
	}
	OkxTimeOut = &cli.IntFlag{
		Name:     "okx-timeout",
		Usage:    "The timeout of the okx",
		EnvVars:  prefixEnvVars("OKX_TIMEOUT"),
		Required: true,
	}

	// bybit flags
	ByBitApiKeyFlag = &cli.StringFlag{
		Name:    "bybit-api-key",
		Usage:   "The apikey of the bybit",
		EnvVars: prefixEnvVars("BYBIT_API_KEY"),
	}
	ByBitApiSecretKeyFlag = &cli.StringFlag{
		Name:    "bybit-api-secret-key",
		Usage:   "The secret of the bybit",
		EnvVars: prefixEnvVars("BYBIT_API_SECRET_KEY"),
	}
	ByBitApiUrlFlag = &cli.StringFlag{
		Name:    "bybit-api-url",
		Usage:   "The api url of the bybit",
		EnvVars: prefixEnvVars("BYBIT_API_URL"),
	}
	ByBitWsUrlFlag = &cli.StringFlag{
		Name:     "bybit-ws-url",
		Usage:    "The ws url of the bybit",
		EnvVars:  prefixEnvVars("BYBIT_WS_URL"),
		Required: true,
	}
	ByBitWsUrlFeature = &cli.StringFlag{
		Name:    "bybit-ws-url-feature",
		Usage:   "The ws url of the bybit",
		EnvVars: prefixEnvVars("BYBIT_WS_URL_FEATURE"),
	}
	ByBitPassphrase = &cli.StringFlag{
		Name:    "bybit-passphrase",
		Usage:   "The passphrase of the bybit",
		EnvVars: prefixEnvVars("BYBIT_PASSPHRASE"),
	}
	ByBitTimeOut = &cli.IntFlag{
		Name:     "bybit-timeout",
		Usage:    "The timeout of the bybit",
		EnvVars:  prefixEnvVars("BYBIT_TIMEOUT"),
		Required: true,
	}

	// bitget flags
	BitGetApiKeyFlag = &cli.StringFlag{
		Name:    "bitget-api-key",
		Usage:   "The apikey of the bitget",
		EnvVars: prefixEnvVars("BITGET_API_KEY"),
	}
	BitGetApiSecretKeyFlag = &cli.StringFlag{
		Name:    "bitget-api-secret-key",
		Usage:   "The secret of the bitget",
		EnvVars: prefixEnvVars("BITGET_API_SECRET_KEY"),
	}
	BitGetApiUrlFlag = &cli.StringFlag{
		Name:    "bitget-api-url",
		Usage:   "The api url of the bitget",
		EnvVars: prefixEnvVars("BITGET_API_URL"),
	}
	BitGetWsUrlFlag = &cli.StringFlag{
		Name:     "bitget-ws-url",
		Usage:    "The ws url of the bitget",
		EnvVars:  prefixEnvVars("BITGET_WS_URL"),
		Required: true,
	}
	BitGetPassphrase = &cli.StringFlag{
		Name:    "bitget-passphrase",
		Usage:   "The passphrase of the bitget",
		EnvVars: prefixEnvVars("BITGET_PASSPHRASE"),
	}
	BitGetTimeOut = &cli.IntFlag{
		Name:     "bitget-timeout",
		Usage:    "The timeout of the bitget",
		EnvVars:  prefixEnvVars("BITGET_TIMEOUT"),
		Required: true,
	}
)

var requireFlags = []cli.Flag{
	MigrationsFlag,
	HttpServerHostFlag,
	HttpServerPortFlag,

	SlaveDbHostFlag,
	SlaveDbPortFlag,
	SlaveDbNameFlag,
	SlaveDbPasswordFlag,
	SlaveDbUserFlag,

	RedisAddressFlag,
	RedisPasswordFlag,
	RedisUserNameFlag,
}
var optionalFlags = []cli.Flag{
	BnApiKeyFlag,
	BnApiSecretKeyFlag,
	BnApiUrlFlag,
	BnWsUrlFlag,
	BnWsUrlFeature,
	BnPassphrase,
	BnTimeOut,

	OkxApiKeyFlag,
	OkxApiSecretKeyFlag,
	OkxApiUrlFlag,
	OkxWsUrlFlag,
	OkxPassphrase,
	OkxTimeOut,

	ByBitApiKeyFlag,
	ByBitApiSecretKeyFlag,
	ByBitApiUrlFlag,
	ByBitWsUrlFlag,
	ByBitWsUrlFeature,
	ByBitPassphrase,
	ByBitTimeOut,

	BitGetApiKeyFlag,
	BitGetApiSecretKeyFlag,
	BitGetApiUrlFlag,
	BitGetWsUrlFlag,
	BitGetPassphrase,
	BitGetTimeOut,
}

var Flags []cli.Flag

func init() {
	Flags = append(requireFlags, optionalFlags...)
}
