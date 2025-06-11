package bn

import (
	"github.com/339-Labs/exchange-market/common"
)

type MixMarketDataAPI interface {
	// 合约最新价格
	MixLatestPrice(symbol string) (string, error)
}

func (c *Client) MixLatestPrice(symbol string) (string, error) {
	params := common.NewParams()
	if symbol != "" {
		params["symbol"] = symbol
	}
	rsp, err := c.mixMarketClient.Tickers(params)
	return rsp, err

}
