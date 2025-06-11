package bitget

import (
	"github.com/339-Labs/exchange-market/common"
	"log"
)

type MixMarketDataAPI interface {
	// 获取合约交易对 最新行情
	MixLatestPrice(productType ProductType) (string, error)
	// 获取单个合约交易对 最新行情
	MixSignalLatestPrice(productType ProductType, symbol string) (string, error)
}

func (c *Client) MixLatestPrice(productType ProductType) (string, error) {
	params := common.NewParams()
	params["productType"] = string(productType)

	rsp, err := c.v2MixMarketClient.Tickers(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}
	return rsp, nil
}

func (c *Client) MixSignalLatestPrice(productType ProductType, symbol string) (string, error) {
	params := common.NewParams()
	params["productType"] = string(productType)
	params["symbol"] = symbol

	rsp, err := c.v2MixMarketClient.Tickers(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}
	return rsp, nil
}
