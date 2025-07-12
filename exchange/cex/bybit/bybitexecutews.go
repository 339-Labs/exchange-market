package bybit

import (
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/ethereum/go-ethereum/log"
	"strings"
	"time"
)

type ByBitExClient struct {
	ByBitWebSocketClient *ByBitWebSocketClient
	config               *config.CexExchangeConfig
	spotPriceMap         *maps.PriceMap
	featurePriceMap      *maps.PriceMap
}

func NewByBitExClient(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) (*ByBitExClient, error) {
	// 创建bitget WebSocket客户端
	client := NewByBitWebSocketClient(config, false) // true表示需要登录

	return &ByBitExClient{
		ByBitWebSocketClient: client,
		config:               config,
		spotPriceMap:         spotPriceMap,
		featurePriceMap:      featurePriceMap,
	}, nil
}

func (bb *ByBitExClient) ExecuteSpotWs() {

	// 设置全局消息监听器
	bb.ByBitWebSocketClient.SetListeners(
		func(message string) {
			fmt.Printf("收到消息: %s\n", message)
		},
		func(message string) {
			fmt.Printf("收到错误: %s\n", message)
		},
	)

	// 启动客户端
	if err := bb.ByBitWebSocketClient.Start(); err != nil {
		log.Info("启动客户端失败:", err)
	}

	// 等待登录完成
	time.Sleep(2 * time.Second)

	// todo spotSymbols get spot from db
	var spotSymbols []string
	spotSymbols = append(spotSymbols, "BTCUSDT")
	spotSymbols = append(spotSymbols, "ETHUSDT")

	err := bb.ByBitWebSocketClient.SubscribeList(spotSymbols, func(message string) {
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)
		jsonMap := common.JSONToMap(message)

		topic, _ := jsonMap["topic"].(string)
		ts, _ := jsonMap["ts"].(string)

		if strings.Contains(topic, "tickers") {
			data, _ := jsonMap["data"].(map[string]interface{})
			bb.handlerSpot(data, ts)
		}
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")

}

func (bb *ByBitExClient) ExecuteFeatureWs() {

	// 创建WebSocket客户端
	cf := &config.CexExchangeConfig{
		WsUrl:        bb.config.WsUrlFeature,
		ApiKey:       bb.config.ApiKey,
		ApiSecretKey: bb.config.ApiSecretKey,
	}
	// 创建bitget WebSocket客户端
	client := NewByBitWebSocketClient(cf, false) // true表示需要登录

	// 设置全局消息监听器
	client.SetListeners(
		func(message string) {

			fmt.Printf("收到消息: %s\n", message)

		},
		func(message string) {
			fmt.Printf("收到错误: %s\n", message)
		},
	)

	// 启动客户端
	if err := client.Start(); err != nil {
		log.Info("启动客户端失败:", err)
	}

	// 等待登录完成
	time.Sleep(2 * time.Second)

	// todo spotSymbols get spot from db
	var spotSymbols []string
	spotSymbols = append(spotSymbols, "BTCUSDT")
	spotSymbols = append(spotSymbols, "ETHUSDT")

	err := client.SubscribeList(spotSymbols, func(message string) {
		fmt.Println("订阅价格更新，处理价格: %s\n", message)

		jsonMap := common.JSONToMap(message)

		topic, _ := jsonMap["topic"].(string)
		ts, _ := jsonMap["ts"].(string)

		if strings.Contains(topic, "tickers") {
			data, _ := jsonMap["data"].(map[string]interface{})
			bb.handlerFeature(data, ts)
		}
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")
}

func (bb *ByBitExClient) handlerSpot(spot map[string]interface{}, ts string) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["symbol"], spot["lastPrice"])

	bb.spotPriceMap.Write(spot["symbol"].(string), &maps.PriceData{
		Symbol:    spot["symbol"].(string),
		Price:     spot["lastPrice"].(string),
		Timestamp: ts,
	})

}

func (bb *ByBitExClient) handlerFeature(feature map[string]interface{}, ts string) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["symbol"], feature["lastPrice"], feature["fundingRate"])

	bb.featurePriceMap.Write(feature["symbol"].(string), &maps.PriceData{
		Symbol:      feature["symbol"].(string),
		Price:       feature["lastPrice"].(string),
		FundingRate: feature["fundingRate"].(string),
		MarkPrice:   feature["markPrice"].(string),
		Timestamp:   ts,
	})

}
