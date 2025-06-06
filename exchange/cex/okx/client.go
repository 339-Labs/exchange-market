package okx

import (
	"github.com/339-Labs/exchange-market/config"
	config2 "github.com/339-Labs/okx-api-sdk-go/config"
	"github.com/339-Labs/okx-api-sdk-go/pkg/client"
	v5 "github.com/339-Labs/okx-api-sdk-go/pkg/client/v5"
)

const CexName = "Okx"

type Client struct {
	okxApiClient     *client.OkxApiClient
	okxAccountClient *v5.AccountClient
	okxMarketClient  *v5.MarketClient
	okxTradeClient   *v5.TradeClient
}

func NewClient(config config.CexExchangeConfig) (OkxClient, error) {
	okxConfig := config2.NewOkxConfig(config.ApiKey, config.ApiSecretKey, config.Passphrase, int(config.TimeOut), "", config.WsUrl)
	apiClient := new(client.OkxApiClient).Init(okxConfig)
	accountClient := new(v5.AccountClient).Init(okxConfig)
	marketClient := new(v5.MarketClient).Init(okxConfig)
	tradeClient := new(v5.TradeClient).Init(okxConfig)
	return &Client{
		okxApiClient:     apiClient,
		okxAccountClient: accountClient,
		okxMarketClient:  marketClient,
		okxTradeClient:   tradeClient,
	}, nil
}

type OkxClient interface {
	// trade market data
	TradeAPI
	// public market data
	PublicDataApi
}
