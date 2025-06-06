package bn

import (
	"errors"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/exchange/cex"
)

type MarketDataAPI interface {
	Tickers(instType cex.InstType, symbol string) (string, error)
}

func (c *Client) Tickers(instType cex.InstType, symbol string) (string, error) {

	params := common.NewParams()
	if symbol != "" {
		params["symbol"] = symbol
	}
	if instType == cex.Spot {
		rsp, err := c.spotMarketClient.Tickers(params)
		return rsp, err
	} else if instType == cex.SWAP {
		rsp, err := c.mixMarketClient.Tickers(params)
		return rsp, err
	}
	return "", errors.New("invalid inst type")

}
