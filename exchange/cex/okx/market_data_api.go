package okx

import (
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/exchange/cex"
)

type TradeAPI interface {
	Tickers(instType cex.InstType) (string, error)
	Ticker(instId string) (string, error)
}

func (c *Client) Tickers(instType cex.InstType) (string, error) {
	params := common.NewParams()
	params["instType"] = string(instType)
	rsp, err := c.okxMarketClient.Tickers(params)
	return rsp, err
}

func (c *Client) Ticker(instId string) (string, error) {
	params := common.NewParams()
	params["instId"] = instId
	rsp, err := c.okxMarketClient.Ticker(params)
	return rsp, err
}
