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

func ExecuteSpotWs(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap) {

	// 创建bitget WebSocket客户端
	client := NewByBitWebSocketClient(config, false) // true表示需要登录

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
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)
		jsonMap := common.JSONToMap(message)

		topic, _ := jsonMap["topic"].(string)

		if strings.Contains(topic, "tickers") {
			data, _ := jsonMap["data"].(map[string]interface{})
			handlerSpot(data, spotPriceMap)
		}
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")

}

func ExecuteFeatureWs(conf *config.CexExchangeConfig, featurePriceMap *maps.PriceMap) {

	// 创建WebSocket客户端
	cf := &config.CexExchangeConfig{
		WsUrl:        conf.WsUrlFeature,
		ApiKey:       conf.ApiKey,
		ApiSecretKey: conf.ApiSecretKey,
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

		if strings.Contains(topic, "tickers") {
			data, _ := jsonMap["data"].(map[string]interface{})
			handlerFeature(data, featurePriceMap)
		}
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")
}

func handlerSpot(spot map[string]interface{}, spotPriceMap *maps.PriceMap) {
	fmt.Println("spot ------ ,instId: %s , lastPr: %s", spot["symbol"], spot["lastPrice"])
}

func handlerFeature(feature map[string]interface{}, featurePriceMap *maps.PriceMap) {
	fmt.Println("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["symbol"], feature["lastPrice"], feature["fundingRate"])
}
