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

type BitGetExClient struct {
	BitGetWebSocketClient *BitGetWebSocketClient
	config                *config.CexExchangeConfig
	spotPriceMap          *maps.PriceMap
	featurePriceMap       *maps.PriceMap
}

func NewBitGetExClient(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) (*BitGetExClient, error) {
	// 创建bitget WebSocket客户端
	client := NewBitGetWebSocketClient(config, false) // true表示需要登录
	return &BitGetExClient{
		BitGetWebSocketClient: client,
		config:                config,
		spotPriceMap:          spotPriceMap,
		featurePriceMap:       featurePriceMap,
	}, nil
}

// bitget 通过 InstId 来区别现货还是合约 BTC-USDT 和 BTC-USD-SWAP
func (bg *BitGetExClient) ExecuteWs() {

	// 创建bitget WebSocket客户端
	client := NewBitGetWebSocketClient(bg.config, false) // true表示需要登录

	// 设置全局消息监听器
	client.SetListeners(
		func(message string) {

			fmt.Println("全局消息: %s\n", message)

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
	spotSymbols = append(spotSymbols, "ETHUSDT")
	spotSymbols = append(spotSymbols, "BTCUSDT")
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
	featureSymbols = append(featureSymbols, "ETHUSDT")
	featureSymbols = append(featureSymbols, "BTCUSDT")
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

		jsonMap := common.JSONToMap(message)
		arg, _ := jsonMap["arg"].(map[string]interface{})
		channel, _ := arg["channel"].(string)
		instType, _ := arg["instType"].(string)

		if dataList, ok := jsonMap["data"].([]interface{}); ok {
			data := dataList[0].(map[string]interface{})
			if channel == "ticker" && instType == "SPOT" {
				bg.handlerSpot(data)
			} else if channel == "ticker" && instType == "USDT-FUTURES" {
				bg.handlerFeature(data)
			}
		}
	})

	if err != nil {
		log.Error("订阅失败: %v", err)
	}
	log.Info("执行")
}

func (bg *BitGetExClient) handlerSpot(spot map[string]interface{}) {

	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["instId"], spot["lastPr"])
	bg.spotPriceMap.Write(spot["instId"].(string), &maps.PriceData{
		Symbol:    spot["instId"].(string),
		Price:     spot["lastPr"].(string),
		Timestamp: spot["ts"].(string),
	})

}

func (bg *BitGetExClient) handlerFeature(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["instId"], feature["lastPr"], feature["fundingRate"])
	bg.featurePriceMap.Write(feature["instId"].(string), &maps.PriceData{
		Symbol:      feature["instId"].(string),
		Price:       feature["lastPr"].(string),
		FundingRate: feature["fundingRate"].(string),
		MarkPrice:   feature["markPrice"].(string),
		Timestamp:   feature["ts"].(string),
	})
}
