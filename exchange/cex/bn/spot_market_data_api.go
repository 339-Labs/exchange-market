package bn

import (
	"github.com/339-Labs/exchange-market/common"
)

type MarketDataAPI interface {
	// 现货最新价格
	SpotLatestPrice(symbol string) (string, error)
}

func (c *Client) SpotLatestPrice(symbol string) (string, error) {
	params := common.NewParams()
	if symbol != "" {
		params["symbol"] = symbol
	}
	rsp, err := c.spotMarketClient.Tickers(params)
	return rsp, err
}
