package bitget

import (
	"github.com/339-Labs/exchange-market/common"
	"log"
)

func (c *Client) CoinsInfo(coin string) (string, error) {
	rsp, err := c.v2SpotMarketClient.Coins()
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}

	return rsp, nil
}

func (c *Client) SpotSymbols(symbol string) (string, error) {
	params := common.NewParams()
	params["symbol"] = symbol
	rsp, err := c.v2SpotMarketClient.Symbols(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}
	return rsp, nil
}

func (c *Client) SpotTickers(symbol string) (string, error) {
	params := common.NewParams()
	params["symbol"] = symbol
	rsp, err := c.v2SpotMarketClient.Tickers(params)

	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}

	return rsp, nil
}
