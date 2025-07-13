package bitget

import (
	"context"
	"errors"
	"github.com/339-Labs/exchange-market/common"
	symbol2 "github.com/339-Labs/exchange-market/database/symbol"
	"time"
)

var TypeTransformError = errors.New("type transform error")

func (c *Client) InitSpotSymbol() error {

	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v2/spot/public/symbols", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		jsonMap := common.BytesToMap(resp.Body)

		if msg, ok := jsonMap["msg"].(string); ok && msg == "success" {

			datas, _ := jsonMap["data"].([]map[string]interface{})

			symbols := make([]symbol2.MarketSymbol, 5000, len(datas))

			for _, data := range datas {

				symbol, _ := data["symbol"].(string)
				baseCoin, _ := data["baseCoin"].(string)
				quoteCoin, _ := data["quoteCoin"].(string)
				var marketSymbol = symbol2.MarketSymbol{
					Symbol:        symbol,
					UnifiedSymbol: baseCoin + "/" + quoteCoin,
					InstType:      "Spot",
					Exchange:      "BitGet",
					ChainId:       "999999",
					Base:          baseCoin,
					Quote:         quoteCoin,
					Timestamp:     uint64(time.Now().UnixMilli()),
				}
				symbols = append(symbols, marketSymbol)
			}

			return nil
		}
	}
	return errors.New("bitget: init spot symbol error")
}

func (c *Client) InitFeatureSymbol() error {
	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v2/mix/market/contracts?productType=USDT-FUTURES", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		jsonMap := common.BytesToMap(resp.Body)

		if msg, ok := jsonMap["msg"].(string); ok && msg == "success" {

			datas, _ := jsonMap["data"].([]map[string]interface{})

			symbols := make([]symbol2.MarketSymbol, 5000, len(datas))

			for _, data := range datas {

				symbol, _ := data["symbol"].(string)
				baseCoin, _ := data["baseCoin"].(string)
				quoteCoin, _ := data["quoteCoin"].(string)
				var marketSymbol = symbol2.MarketSymbol{
					Symbol:        symbol,
					UnifiedSymbol: baseCoin + "/" + quoteCoin,
					InstType:      "Feature",
					Exchange:      "BitGet",
					ChainId:       "999999",
					Base:          baseCoin,
					Quote:         quoteCoin,
					Timestamp:     uint64(time.Now().UnixMilli()),
				}
				symbols = append(symbols, marketSymbol)
			}

			return nil
		}
	}
	return errors.New("bitget: init feature symbol error")
}
