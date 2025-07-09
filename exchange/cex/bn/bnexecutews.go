package bn

import (
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/config"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

func ExecuteWsSpot(config *config.CexExchangeConfig) {

	// 创建WebSocket客户端
	client := NewBnWebSocketClient(config, false)

	// 设置全局监听器
	client.SetListeners(
		func(message string) {
			log.Info("Global message: %s", message)

			// 解析基本消息结构
			baseMsgs := common.JSONToArrMap(message)
			for _, baseMsg := range baseMsgs {
				eventType, ok := baseMsg["e"].(string)
				if ok {

					// 根据事件类型处理不同的消息
					switch eventType {
					case "kline":
						log.Info("Kline %s", message)
					case "24hrMiniTicker":
						handlerSpot(baseMsg)
					case "depthUpdate":
						log.Info("depthUpdate %s", message)
					case "trade":
						log.Info("Trade %s", message)
					default:
						return
					}

				}
			}

		},
		func(error string) {
			log.Error("Error: %s", error)
		},
	)

	// 连接
	if err := client.Start(); err != nil {
		log.Error("Failed to connect:", err)
	}

	// 等待连接完成
	time.Sleep(2 * time.Second)

	// 订阅所有交易对 24小时价格变动统计
	err := client.SubscribeMiniTickerAll(func(message string) {
		log.Info("mini Ticker data: %s", message)
	})
	if err != nil {
		log.Error("Failed to subscribe ticker: %v", err)
	}

}

func ExecuteWsFeature(conf *config.CexExchangeConfig) {

	// 创建WebSocket客户端
	cf := &config.CexExchangeConfig{
		WsUrl:        conf.WsUrlFeature,
		ApiKey:       conf.ApiKey,
		ApiSecretKey: conf.ApiSecretKey,
	}
	client := NewBnWebSocketClient(cf, false)

	// 设置全局监听器
	client.SetListeners(
		func(message string) {
			log.Info("Global message: %s", message)

			// 解析基本消息结构
			baseMsgs := common.JSONToArrMap(message)
			for _, baseMsg := range baseMsgs {
				eventType, ok := baseMsg["e"].(string)
				if ok {

					// 根据事件类型处理不同的消息
					switch eventType {
					case "kline":
						log.Info("Kline %s", message)
					case "24hrMiniTicker":
						handlerFeature(baseMsg)
					case "depthUpdate":
						log.Info("depthUpdate %s", message)
					case "trade":
						log.Info("Trade %s", message)
					case "markPriceUpdate":
						handlerFeatureMark(baseMsg)
					default:
						return
					}

				}
			}

		},
		func(error string) {
			log.Error("Error: %s", error)
		},
	)

	// 连接
	if err := client.Start(); err != nil {
		log.Error("Failed to connect:", err)
	}

	// 等待连接完成
	time.Sleep(2 * time.Second)

	// 全部交易对 订阅24小时价格变动统计
	err := client.SubscribeAllTicker(func(message string) {
		log.Info("mini Ticker data: %s", message)
	})
	if err != nil {
		log.Info("Failed to subscribe ticker: %v", err)
	}

	// 全部交易对标记价格 订阅24小b价格变动统计
	err = client.SubscribeMiniTickerAll(func(message string) {
		log.Info("mini Ticker data: %s", message)
	})
	if err != nil {
		log.Info("Failed to subscribe ticker: %v", err)
	}

}

func handlerSpot(spot map[string]interface{}) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["s"], spot["c"])
}

func handlerFeature(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["s"], feature["c"])
}

func handlerFeatureMark(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["s"], feature["p"], feature["r"])
}
