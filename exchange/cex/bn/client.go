package bn

import (
	config2 "github.com/339-Labs/binance-api-sdk-go/config"
	"github.com/339-Labs/binance-api-sdk-go/pkg/client"
	"github.com/339-Labs/binance-api-sdk-go/pkg/client/v"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex"
)

const CexName = "BN"

type Client struct {
	spotBnClient     *client.BnApiClient
	mixBnClient      *client.BnApiClient
	mixMarketClient  *v.MixMarketClient
	spotMarketClient *v.SpotMarketClient
}

func NewClient(config config.CexExchangeConfig) (BnClient, error) {
	bnConfig := config2.NewBnConfig(config.ApiKey, config.ApiSecretKey, int(config.TimeOut), "")
	return &Client{
		spotBnClient:     new(client.BnApiClient).Init(bnConfig, string(cex.Spot)),
		mixBnClient:      new(client.BnApiClient).Init(bnConfig, string(cex.SWAP)),
		mixMarketClient:  new(v.MixMarketClient).Init(bnConfig),
		spotMarketClient: new(v.SpotMarketClient).Init(bnConfig),
	}, nil
}

type BnClient interface {
	MarketDataAPI
}
