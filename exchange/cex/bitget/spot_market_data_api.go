package bitget

import (
	"github.com/339-Labs/exchange-market/common"
	"log"
)

type SpotMarketDataAPI interface {
	// 获取币种信息  如不填写，默认返回全部币种信息
	SymbolsInfo(coin string) (string, error)
	// 获取现货交易对信息，如不填写，默认返回全部交易对信息
	SpotSymbolsInfo(symbol string) (string, error)
	//获取最新价格行情信息 如不填写，默认返回全部交易对信息
	SpotLatestPrice(symbol string) (string, error)
}

func (c *Client) SymbolsInfo(coin string) (string, error) {
	rsp, err := c.v2SpotMarketClient.Coins()
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}

	return rsp, nil
}

func (c *Client) SpotSymbolsInfo(symbol string) (string, error) {
	params := common.NewParams()
	params["symbol"] = symbol
	rsp, err := c.v2SpotMarketClient.Symbols(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}
	return rsp, nil
}

func (c *Client) SpotLatestPrice(symbol string) (string, error) {
	params := common.NewParams()
	params["symbol"] = symbol
	rsp, err := c.v2SpotMarketClient.Tickers(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}

	return rsp, nil
}
