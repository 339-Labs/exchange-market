package bitget

import (
	"errors"
	"github.com/339-Labs/exchange-market/api/service"
	"github.com/339-Labs/exchange-market/common/client"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
)

var errBlockChainHTTPError = errors.New("BitGet http request error")

type Client struct {
	config config.Config
	db     database.DB
	resty  client.REST
}

func NewClient(config config.Config, db database.DB) service.HandlerSymbolAdaptor {
	rest := client.NewRESTClient(config.ExchangeConfig.BitGet.ApiUrl)

	return &Client{
		config: config,
		db:     db,
		resty:  rest,
	}
}
