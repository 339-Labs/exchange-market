package main

import (
	"fmt"
	"github.com/339-Labs/exchange-market/exchange/cex/bitget"
	"github.com/339-Labs/exchange-market/exchange/cex/bitget/model"
	"log"
	"time"

	"github.com/339-Labs/exchange-market/config"
)

func main() {
	// 示例1: 使用okx特定的WebSocket客户端
	bnExample()
}

// okxExample 使用okx WebSocket客户端的示例
func bnExample() {

	// 配置
	config := &config.CexExchangeConfig{
		WsUrl:        "wss://ws.bitget.com/v2/ws/public",
		ApiKey:       "your-api-key",
		ApiSecretKey: "your-api-secret",
	}

	// 创建bitget WebSocket客户端
	client := bitget.NewBitGetWebSocketClient(config, false) // true表示需要登录

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
		log.Printf("启动客户端失败:", err)
	}

	// 等待登录完成
	time.Sleep(2 * time.Second)

	// 订阅特定的数据流
	subscribeReq := model.SubscribeReq{
		Channel:  "ticker",
		InstId:   "BTCUSDT",
		InstType: "SPOT",
	}
	var reqs []model.SubscribeReq
	reqs = append(reqs, subscribeReq)

	subscribeReq1 := model.SubscribeReq{
		Channel:  "ticker",
		InstId:   "ETHUSDT",
		InstType: "SPOT",
	}
	reqs = append(reqs, subscribeReq1)

	err := client.SubscribeList(reqs, func(message string) {
		fmt.Printf("订阅价格更新，处理价格: %s\n", message)
	})

	if err != nil {
		log.Fatalln("订阅失败: %v", err)
	}
	log.Println("执行")

	// 等待一段时间再关闭
	time.Sleep(120 * time.Second)

	// 关闭连接
	client.Stop()

}
