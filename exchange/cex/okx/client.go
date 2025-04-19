package okx

import (
	"github.com/339-Labs/exchange-market/common"
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
	Tickers(instType InstType) (string, error)
	Ticker(instId string) (string, error)
}

func (c *Client) Tickers(instType InstType) (string, error) {
	params := common.NewParams()
	params["instType"] = string(instType)
	rsp, err := c.okxMarketClient.Tickers(params)
	if err != nil {
		return "", err
	}
	return rsp, nil
}

func (c *Client) Ticker(instId string) (string, error) {
	params := common.NewParams()
	params["instId"] = instId
	rsp, err := c.okxMarketClient.Ticker(params)
	if err != nil {
		return "", err
	}
	return rsp, nil
}
