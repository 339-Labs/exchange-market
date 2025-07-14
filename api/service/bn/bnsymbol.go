package bn

import (
	"context"
	"errors"
	"github.com/339-Labs/exchange-market/common"
	symbol2 "github.com/339-Labs/exchange-market/database/symbol"
	"github.com/ethereum/go-ethereum/log"
	"strings"
	"time"
)

func (c *Client) InitSpotSymbol() error {

	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v3/ticker/price", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		arrMap := common.BytesToArrMap(resp.Body)

		symbols := make([]symbol2.MarketSymbol, len(arrMap), len(arrMap))
		for _, v := range arrMap {
			symbol, _ := v["symbol"].(string)

			sym, base, quote, err := c.handlerSymbol(symbol)
			if err != nil {
				continue
			} else {

				var marketSymbol = symbol2.MarketSymbol{
					Symbol:        symbol,
					UnifiedSymbol: sym,
					InstType:      "Spot",
					Exchange:      "Bn",
					ChainId:       "999999",
					Base:          base,
					Quote:         quote,
					Timestamp:     uint64(time.Now().UnixMilli()),
				}
				symbols = append(symbols, marketSymbol)

			}
		}
		log.Info("====", len(symbols))
	}
	return errors.New("bn: init spot symbol error")
}

func (c *Client) InitFeatureSymbol() error {
	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/fapi/v1/ticker/price", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		arrMap := common.BytesToArrMap(resp.Body)

		symbols := make([]symbol2.MarketSymbol, len(arrMap), len(arrMap))
		for _, v := range arrMap {
			symbol, _ := v["symbol"].(string)

			sym, base, quote, err := c.handlerSymbol(symbol)
			if err != nil {
				continue
			} else {

				var marketSymbol = symbol2.MarketSymbol{
					Symbol:        symbol,
					UnifiedSymbol: sym,
					InstType:      "Feature",
					Exchange:      "Bn",
					ChainId:       "999999",
					Base:          base,
					Quote:         quote,
					Timestamp:     uint64(time.Now().UnixMilli()),
				}
				symbols = append(symbols, marketSymbol)

			}
		}
		log.Info("====", len(symbols))
	}
	return errors.New("bn: init feature symbol error")
}

func (c *Client) handlerSymbol(symbol string) (string, string, string, error) {
	if strings.HasSuffix(symbol, "USDT") {
		base := strings.TrimSuffix(symbol, "USDT")
		quote := "USDT"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "BTC") {
		base := strings.TrimSuffix(symbol, "BTC")
		quote := "BTC"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "ETH") {
		base := strings.TrimSuffix(symbol, "ETH")
		quote := "ETH"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "USDC") {
		base := strings.TrimSuffix(symbol, "USDC")
		quote := "USDC"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "BNB") {
		base := strings.TrimSuffix(symbol, "BNB")
		quote := "BNB"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "TRY") {
		base := strings.TrimSuffix(symbol, "TRY")
		quote := "TRY"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "FDUSD") {
		base := strings.TrimSuffix(symbol, "FDUSD")
		quote := "FDUSD"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "JPY") {
		base := strings.TrimSuffix(symbol, "JPY")
		quote := "JPY"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "EUR") {
		base := strings.TrimSuffix(symbol, "EUR")
		quote := "EUR"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "MXN") {
		base := strings.TrimSuffix(symbol, "MXN")
		quote := "MXN"
		return base + "/" + quote, base, quote, nil
	}
	if strings.HasSuffix(symbol, "BRL") {
		base := strings.TrimSuffix(symbol, "BRL")
		quote := "BRL"
		return base + "/" + quote, base, quote, nil
	}
	return "", "", "", errors.New("bn: handler symbol error")
}
