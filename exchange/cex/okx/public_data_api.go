package okx

import (
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/exchange/cex"
)

type PublicDataApi interface {
	InstrumentsInfo(instType cex.InstType, instId string) (string, error)
	MarkPrice(instType cex.InstType, instId string) (string, error)
}

func (c *Client) InstrumentsInfo(instType cex.InstType, instId string) (string, error) {
	params := common.NewParams()
	params["instType"] = string(instType)
	params["instId"] = instId
	rsp, err := c.okxApiClient.OkxRestClient.DoGetNoAuth("/api/v5/public/instruments", params)
	return rsp, err
}

func (c *Client) MarkPrice(instType cex.InstType, instId string) (string, error) {
	params := common.NewParams()
	params["instType"] = string(instType)
	params["instId"] = instId
	rsp, err := c.okxApiClient.OkxRestClient.DoGetNoAuth("/api/v5/public/mark-price", params)
	return rsp, err
}
