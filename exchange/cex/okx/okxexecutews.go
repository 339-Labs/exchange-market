package okx

import (
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/okx/model"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

func ExecuteWs(config *config.CexExchangeConfig) {

	// 创建okx WebSocket客户端
	client := NewOkxWebSocketClient(config, false) // true表示需要登录

	// 设置全局消息监听器
	client.SetListeners(
		func(message string) {

			jsonMap := common.JSONToMap(message)
			if arg, exists := jsonMap["arg"].(map[string]interface{}); exists {

				channel, _ := arg["channel"].(string)

				dataList, _ := jsonMap["data"].([]interface{})
				data := dataList[0].(map[string]interface{})
				instType := data["instType"].(string)

				switch channel {

				case "tickers":
					if instType == "SWAP" {
						handlerFeature(data)
					} else if instType == "SPOT" {
						handlerSpot(data)
					}
				case "funding-rate":
					if instType == "SWAP" {
						handlerFeatureRate(data)
					}
				case "mark-price":
					if instType == "SWAP" {
						handlerFeatureMark(data)
					}
				}

			}
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

	var reqs []model.SubscribeReq

	// InstId 来区别现货还是合约 BTC-USDT 和 BTC-USD-SWAP
	// todo spotSymbols get feature from db
	var spotSymbols []string
	for _, symbol := range spotSymbols {
		// 订阅特定合约的数据流
		subscribeReq := model.SubscribeReq{
			Channel: "tickers",
			InstId:  symbol,
		}
		reqs = append(reqs, subscribeReq)
	}

	// todo featureSymbols get feature from db
	var featureSymbols []string
	for _, symbol := range featureSymbols {
		// 订阅特定合约的数据流
		subscribeReq := model.SubscribeReq{
			Channel: "tickers",
			InstId:  symbol,
		}
		reqs = append(reqs, subscribeReq)
	}

	err := client.SubscribeList(reqs, func(message string) {
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")

}

func handlerSpot(spot map[string]interface{}) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["instId"], spot["last"])
}

func handlerFeature(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s ", feature["instId"], feature["last"])
}

func handlerFeatureMark(feature map[string]interface{}) {
	log.Info("feature instType: %s ------ ,instId: %s , markPx: %s ", feature["instType"], feature["instId"], feature["markPx"])
}

func handlerFeatureRate(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , fundingRate: %s", feature["instId"], feature["fundingRate"])
}
