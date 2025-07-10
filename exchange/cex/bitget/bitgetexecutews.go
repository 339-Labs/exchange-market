package bitget

import (
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/bitget/model"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

// bitget 通过 InstId 来区别现货还是合约 BTC-USDT 和 BTC-USD-SWAP
func ExecuteWs(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) {

	// 创建bitget WebSocket客户端
	client := NewBitGetWebSocketClient(config, false) // true表示需要登录

	// 设置全局消息监听器
	client.SetListeners(
		func(message string) {

			jsonMap := common.JSONToMap(message)
			if arg, exists := jsonMap["arg"].(map[string]interface{}); exists {

				channel, _ := arg["channel"].(string)
				instType, _ := arg["instType"].(string)

				dataList, _ := jsonMap["data"].([]interface{})
				data := dataList[0].(map[string]interface{})

				if channel == "ticker" && instType == "SPOT" {
					handlerSpot(data, spotPriceMap)
				} else if channel == "ticker" && instType == "USDT-FUTURES" {
					handlerFeature(data, featurePriceMap)
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

	// todo spotSymbols get spot from db
	var spotSymbols []string
	for _, symbol := range spotSymbols {
		// 订阅特定现货的数据流
		subscribeReq := model.SubscribeReq{
			Channel:  "ticker",
			InstId:   symbol,
			InstType: "SPOT",
		}
		reqs = append(reqs, subscribeReq)
	}

	// todo spotSymbols get feature from db
	var featureSymbols []string
	for _, symbol := range featureSymbols {
		// 订阅特定合约的数据流
		subscribeReq := model.SubscribeReq{
			Channel:  "ticker",
			InstId:   symbol,
			InstType: "USDT-FUTURES",
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

func handlerSpot(spot map[string]interface{}, spotPriceMap *maps.PriceMap) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["instId"], spot["lastPr"])
}

func handlerFeature(feature map[string]interface{}, featurePriceMap *maps.PriceMap) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["instId"], feature["lastPr"], feature["fundingRate"])
}
