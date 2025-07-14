package bybit

import (
	"github.com/339-Labs/exchange-market/api/service"
	"github.com/339-Labs/exchange-market/common/client"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
)

type Client struct {
	config config.Config
	db     database.DB
	resty  client.REST
}

func NewClient(config config.Config, db database.DB) service.HandlerSymbolAdaptor {
	rest := client.NewRESTClient(config.ExchangeConfig.ByBit.ApiUrl)
	return &Client{
		config: config,
		db:     db,
		resty:  rest,
	}
}
