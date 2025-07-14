package bitget

import (
	"context"
	"errors"
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/client"
	symbol2 "github.com/339-Labs/exchange-market/database/symbol"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

var TypeTransformError = errors.New("type transform error")

func (c *Client) InitSpotSymbol() error {

	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v2/spot/public/symbols", header)
	if err != nil {
		return err
	}

	symbols, err := c.parseMarketResponse(resp, "Spot")
	log.Info("====", len(symbols))
	return errors.New("bitget: init spot symbol error")
}

// API响应结构体
type APIResponse struct {
	Code        string        `json:"code"`
	Msg         string        `json:"msg"`
	RequestTime int64         `json:"requestTime"`
	Data        []interface{} `json:"data"`
}

// 响应结果结构体
type ResponseResult struct {
	Response *client.RESTResponse // 假设这是您的响应类型
	Error    error
}

func (c *Client) InitFeatureSymbol() error {

	productTypes := []string{"USDT-FUTURES", "USDC_FUTURES", "COIN_FUTURES"}

	header := make(map[string]string, 0)

	var allSymbols []symbol2.MarketSymbol

	// 并发请求所有产品类型
	responses := make(chan ResponseResult, len(productTypes))

	for _, productType := range productTypes {
		go func(pt string) {
			url := fmt.Sprintf("/api/v2/mix/market/contracts?productType=%s", pt)
			resp, err := c.resty.GET(context.Background(), url, header)
			responses <- ResponseResult{Response: resp, Error: err}
		}(productType)
	}

	// 收集所有响应
	for i := 0; i < len(productTypes); i++ {
		result := <-responses
		if result.Error != nil {
			return result.Error
		}

		symbols, err := c.parseMarketResponse(result.Response, "Feature")
		if err != nil {
			return err
		}

		allSymbols = append(allSymbols, symbols...)
	}

	log.Info("====", len(allSymbols))

	return errors.New("bitget: init feature symbol error")
}

// 解析市场响应的通用方法
func (c *Client) parseMarketResponse(resp *client.RESTResponse, InstType string) ([]symbol2.MarketSymbol, error) {
	if !resp.IsSuccess() || resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: status code %d", resp.StatusCode)
	}

	jsonMap := common.BytesToMap(resp.Body)

	// 检查响应码
	code, ok := jsonMap["code"].(string)
	if !ok || code != "00000" {
		return nil, fmt.Errorf("api error: code=%v, msg=%v", jsonMap["code"], jsonMap["msg"])
	}

	// 检查消息
	msg, ok := jsonMap["msg"].(string)
	if !ok || msg != "success" {
		return nil, fmt.Errorf("api response error: %v", jsonMap["msg"])
	}

	// 解析data字段
	datas, ok := jsonMap["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	symbols := make([]symbol2.MarketSymbol, 0, len(datas))
	currentTime := uint64(time.Now().UnixMilli())

	for _, item := range datas {
		data, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		symbol, ok := data["symbol"].(string)
		if !ok {
			continue
		}

		baseCoin, ok := data["baseCoin"].(string)
		if !ok {
			continue
		}

		quoteCoin, _ := data["quoteCoin"].(string)

		marketSymbol := symbol2.MarketSymbol{
			Symbol:        symbol,
			UnifiedSymbol: baseCoin + "/" + quoteCoin,
			InstType:      InstType,
			Exchange:      "BitGet",
			ChainId:       "999999",
			Base:          baseCoin,
			Quote:         quoteCoin,
			Timestamp:     currentTime,
		}

		symbols = append(symbols, marketSymbol)
	}

	return symbols, nil
}
