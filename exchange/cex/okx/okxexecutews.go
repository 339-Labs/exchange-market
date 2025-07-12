package okx

import (
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/okx/model"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type OkxExClient struct {
	OkxWebSocketClient *OkxWebSocketClient
	config             *config.CexExchangeConfig
	spotPriceMap       *maps.PriceMap
	featurePriceMap    *maps.PriceMap
	markPriceMap       *maps.PriceMap
	rateMap            *maps.PriceMap
}

func NewOkxExClient(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap, markPriceMap *maps.PriceMap, rateMap *maps.PriceMap) (*OkxExClient, error) {
	// 创建okx WebSocket客户端
	client := NewOkxWebSocketClient(config, false) // true表示需要登录

	return &OkxExClient{
		OkxWebSocketClient: client,
		config:             config,
		spotPriceMap:       spotPriceMap,
		featurePriceMap:    featurePriceMap,
		markPriceMap:       markPriceMap,
		rateMap:            rateMap,
	}, nil
}

func (okx *OkxExClient) ExecuteSpotWs() {

	// 设置全局消息监听器
	okx.OkxWebSocketClient.SetListeners(
		func(message string) {
			fmt.Printf("收到消息: %s\n", message)
		},
		func(message string) {
			fmt.Printf("收到错误: %s\n", message)
		},
	)

	// 启动客户端
	if err := okx.OkxWebSocketClient.Start(); err != nil {
		log.Info("启动客户端失败:", err)
	}

	// 等待登录完成
	time.Sleep(2 * time.Second)

	var reqs []model.SubscribeReq

	// InstId 来区别现货还是合约 BTC-USDT 和 BTC-USD-SWAP
	// todo spotSymbols get feature from db
	var spotSymbols []string
	spotSymbols = append(spotSymbols, "BTC-USDT")
	spotSymbols = append(spotSymbols, "ETH-USDT")
	for _, symbol := range spotSymbols {
		// 订阅特定合约的数据流
		subscribeReq := model.SubscribeReq{
			Channel: "tickers",
			InstId:  symbol,
		}
		reqs = append(reqs, subscribeReq)
	}

	err := okx.OkxWebSocketClient.SubscribeList(reqs, func(message string) {
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)

		jsonMap := common.JSONToMap(message)
		if arg, exists := jsonMap["arg"].(map[string]interface{}); exists {

			channel, _ := arg["channel"].(string)

			dataList, _ := jsonMap["data"].([]interface{})
			data := dataList[0].(map[string]interface{})
			instType := data["instType"].(string)

			switch channel {

			case "tickers":
				if instType == "SWAP" {
					okx.handlerFeature(data)
				} else if instType == "SPOT" {
					okx.handlerSpot(data)
				}
			case "funding-rate":
				if instType == "SWAP" {
					okx.handlerFeatureRate(data)
				}
			case "mark-price":
				if instType == "SWAP" {
					okx.handlerFeatureMark(data)
				}
			}
		}

	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")

}

func (okx *OkxExClient) ExecuteFeatureWs() {

	// 设置全局消息监听器
	okx.OkxWebSocketClient.SetListeners(
		func(message string) {
			fmt.Printf("收到消息: %s\n", message)
		},
		func(message string) {
			fmt.Printf("收到错误: %s\n", message)
		},
	)

	// 启动客户端
	if err := okx.OkxWebSocketClient.Start(); err != nil {
		log.Info("启动客户端失败:", err)
	}

	// 等待登录完成
	time.Sleep(2 * time.Second)

	var reqs []model.SubscribeReq

	// InstId 来区别现货还是合约 BTC-USDT 和 BTC-USD-SWAP

	// todo featureSymbols get feature from db
	var featureSymbols []string
	featureSymbols = append(featureSymbols, "BTC-USDT-SWAP")
	featureSymbols = append(featureSymbols, "ETH-USDT-SWAP")
	for _, symbol := range featureSymbols {
		// 订阅特定合约的数据流
		subscribeReq := model.SubscribeReq{
			Channel: "tickers",
			InstId:  symbol,
		}
		reqs = append(reqs, subscribeReq)
	}

	err := okx.OkxWebSocketClient.SubscribeList(reqs, func(message string) {
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)

		jsonMap := common.JSONToMap(message)
		if arg, exists := jsonMap["arg"].(map[string]interface{}); exists {

			channel, _ := arg["channel"].(string)

			dataList, _ := jsonMap["data"].([]interface{})
			data := dataList[0].(map[string]interface{})
			instType := data["instType"].(string)

			switch channel {

			case "tickers":
				if instType == "SWAP" {
					okx.handlerFeature(data)
				} else if instType == "SPOT" {
					okx.handlerSpot(data)
				}
			case "funding-rate":
				if instType == "SWAP" {
					okx.handlerFeatureRate(data)
				}
			case "mark-price":
				if instType == "SWAP" {
					okx.handlerFeatureMark(data)
				}
			}
		}

	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")

}

func (okx *OkxExClient) handlerSpot(spot map[string]interface{}) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["instId"], spot["last"])

	okx.spotPriceMap.Write(spot["instId"].(string), &maps.PriceData{
		Symbol:    spot["instId"].(string),
		Price:     spot["last"].(string),
		Timestamp: spot["ts"].(string),
	})

}

func (okx *OkxExClient) handlerFeature(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s ", feature["instId"], feature["last"])

	okx.featurePriceMap.Write(feature["instId"].(string), &maps.PriceData{
		Symbol:    feature["instId"].(string),
		Price:     feature["last"].(string),
		Timestamp: feature["ts"].(string),
	})

}

func (okx *OkxExClient) handlerFeatureMark(feature map[string]interface{}) {
	log.Info("feature instType: %s ------ ,instId: %s , markPx: %s ", feature["instType"], feature["instId"], feature["markPx"])

	okx.markPriceMap.Write(feature["instId"].(string), &maps.PriceData{
		Symbol:    feature["instId"].(string),
		MarkPrice: feature["markPx"].(string),
		Timestamp: feature["ts"].(string),
	})
}

func (okx *OkxExClient) handlerFeatureRate(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , fundingRate: %s", feature["instId"], feature["fundingRate"])

	okx.rateMap.Write(feature["instId"].(string), &maps.PriceData{
		Symbol:      feature["instId"].(string),
		FundingRate: feature["fundingRate"].(string),
		Timestamp:   feature["ts"].(string),
	})
}
