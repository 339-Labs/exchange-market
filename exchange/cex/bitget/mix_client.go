package bitget

import (
	"github.com/339-Labs/exchange-market/common"
	"log"
)

func (c *Client) AllTickers(productType ProductType) (string, error) {
	params := common.NewParams()
	params["productType"] = string(productType)

	rsp, err := c.v2MixMarketClient.Tickers(params)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return "", err
	}
	return rsp, nil
}
