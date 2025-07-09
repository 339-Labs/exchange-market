package bn

import (
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/signer"
	"github.com/339-Labs/exchange-market/common/ws"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/bn/constants"
	"github.com/339-Labs/exchange-market/exchange/cex/bn/model"
	"github.com/ethereum/go-ethereum/log"
	"strings"
	"sync"
	"time"
)

// BnMessageHandler 币安消息处理器
type BnMessageHandler struct {
	// 配置和认证
	Config      *config.CexExchangeConfig
	Signer      *signer.Signer
	LoginStatus bool
	NeedLogin   bool

	Listener      OnReceive
	ErrorListener OnReceive
	StreamMap     map[string]OnReceive // 币安使用流名称作为key
	AllSubscribe  *model.Set

	// 同步
	mu sync.RWMutex

	// WebSocket客户端引用
	wsClient *ws.GenericWebSocketClient
}

// OnReceive 消息接收回调函数类型
type OnReceive func(message string)

// NewBnMessageHandler 创建新的币安消息处理器
func NewBnMessageHandler(config *config.CexExchangeConfig, needLogin bool) *BnMessageHandler {
	handler := &BnMessageHandler{
		Config:       config,
		NeedLogin:    needLogin,
		LoginStatus:  false,
		StreamMap:    make(map[string]OnReceive),
		AllSubscribe: model.NewSet(),
		Signer:       new(signer.Signer).Init(config.ApiSecretKey),
	}

	return handler
}

// SetWebSocketClient 设置WebSocket客户端引用
func (h *BnMessageHandler) SetWebSocketClient(client *ws.GenericWebSocketClient) {
	h.wsClient = client
}

// SetListeners 设置消息监听器
func (h *BnMessageHandler) SetListeners(msgListener OnReceive, errorListener OnReceive) {
	h.Listener = msgListener
	h.ErrorListener = errorListener
}

// HandleMessage 处理普通消息
func (h *BnMessageHandler) HandleMessage(message string) error {
	jsonMap := common.JSONToMap(message)

	// 检查是否有错误
	if errorCode, exists := jsonMap["code"]; exists {
		if code, ok := errorCode.(float64); ok && int(code) != 200 {
			return fmt.Errorf("received error code: %d", int(code))
		}
	}

	// 检查是否有错误
	if errorCode, exists := jsonMap["status"]; exists {
		if code, ok := errorCode.(float64); ok && int(code) != 200 {
			return fmt.Errorf("received error code: %d", int(code))
		}
	}

	// 处理订阅确认响应
	if result, exists := jsonMap["result"]; exists {
		if result == nil {
			return h.handleSubscribeResponse(message, jsonMap)
		}
	}

	// 处理数据消息
	if stream, exists := jsonMap["stream"]; exists {
		return h.handleDataMessage(message, stream.(string))
	}

	// 处理其他消息
	return h.handleOtherMessage(message)
}

// HandleError 处理错误消息
func (h *BnMessageHandler) HandleError(message string) error {
	log.Error("Received error message: %s", message)

	if h.ErrorListener != nil {
		h.ErrorListener(message)
	}

	return nil
}

// HandleSpecialMessage 处理特殊消息（如pong）
func (h *BnMessageHandler) HandleSpecialMessage(message string) (handled bool, err error) {
	// 处理pong消息
	if message == "ping" {
		log.Info("Received ping from Binance")
		return true, nil
	}

	// 检查是否是订阅确认消息
	if strings.Contains(message, "result") && strings.Contains(message, "id") {
		log.Info("Binance subscription response: %s", message)
		return true, nil
	}

	return false, nil
}

// handleSubscribeResponse 处理订阅响应
func (h *BnMessageHandler) handleSubscribeResponse(message string, jsonMap map[string]interface{}) error {
	log.Info("Subscribe response: %s", message)

	if h.Listener != nil {
		h.Listener(message)
	}

	return nil
}

// handleDataMessage 处理数据消息
func (h *BnMessageHandler) handleDataMessage(message string, stream string) error {
	listener := h.getListener(stream)
	if listener != nil {
		listener(message)
	}
	return nil
}

// handleOtherMessage 处理其他消息
func (h *BnMessageHandler) handleOtherMessage(message string) error {
	log.Info("Received other message: %s", message)

	if h.Listener != nil {
		h.Listener(message)
	}

	return nil
}

// getListener 获取特定流的监听器
func (h *BnMessageHandler) getListener(stream string) OnReceive {
	h.mu.RLock()
	listener, exists := h.StreamMap[stream]
	h.mu.RUnlock()

	if !exists {
		return h.Listener
	}

	return listener
}

// IsLoggedIn 检查是否已登录（币安现货不需要登录）
func (h *BnMessageHandler) IsLoggedIn() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.LoginStatus
}

// AddSubscription 添加订阅
func (h *BnMessageHandler) AddSubscription(stream string, listener OnReceive) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.StreamMap[stream] = listener
	h.AllSubscribe.Add(stream)
}

// RemoveSubscription 移除订阅
func (h *BnMessageHandler) RemoveSubscription(stream string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.StreamMap, stream)
	h.AllSubscribe.Remove(stream)
}

// BnWebSocketClient 币安 WebSocket客户端
type BnWebSocketClient struct {
	*ws.GenericWebSocketClient
	MessageHandler *BnMessageHandler
}

// NewBnWebSocketClient 创建新的币安 WebSocket客户端
func NewBnWebSocketClient(config *config.CexExchangeConfig, needLogin bool) *BnWebSocketClient {
	// 创建WebSocket配置
	wsConfig := &ws.ConnectionConfig{
		WsUrl:               config.WsUrl,
		PingInterval:        30 * time.Second,
		ReconnectWaitSecond: float64(constants.ReconnectWaitSecond),
		TimerIntervalSecond: constants.TimerIntervalSecond * time.Second,
		EnableAutoReconnect: true,
		EnablePing:          false,
	}

	// 创建通用WebSocket客户端
	genericClient := ws.NewGenericWebSocketClient(wsConfig)

	// 创建币安消息处理器
	messageHandler := NewBnMessageHandler(config, needLogin)
	messageHandler.SetWebSocketClient(genericClient)

	// 设置消息处理器
	genericClient.SetMessageHandler(messageHandler)

	// 创建客户端实例
	client := &BnWebSocketClient{
		GenericWebSocketClient: genericClient,
		MessageHandler:         messageHandler,
	}

	// 设置回调函数
	genericClient.SetCallbacks(
		func() {
			log.Info("Binance WebSocket connected")
			// 币安现货不需要登录
			if needLogin {
				messageHandler.LoginStatus = true
			}
		},
		func() {
			log.Info("Binance WebSocket disconnected")
		},
		func(attempt int) {
			log.Info("Binance WebSocket reconnecting, attempt: %d", attempt)
		},
	)

	return client
}

// Subscribe 订阅单个流
func (c *BnWebSocketClient) Subscribe(stream string, listener OnReceive) error {
	// 添加到订阅映射
	c.MessageHandler.AddSubscription(stream, listener)

	// 构造币安订阅请求
	subscribeReq := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{stream},
		"id":     time.Now().Unix(),
	}

	return c.SendJSON(subscribeReq)
}

// SubscribeList 订阅多个流
func (c *BnWebSocketClient) SubscribeList(streams []string, listener OnReceive) error {
	// 添加到订阅映射
	for _, stream := range streams {
		c.MessageHandler.AddSubscription(stream, listener)
	}

	// 构造币安订阅请求
	subscribeReq := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": streams,
		"id":     time.Now().Unix(),
	}

	return c.SendJSON(subscribeReq)
}

// Unsubscribe 取消订阅单个流
func (c *BnWebSocketClient) Unsubscribe(stream string) error {
	// 从订阅映射中移除
	c.MessageHandler.RemoveSubscription(stream)

	// 构造币安取消订阅请求
	unsubscribeReq := map[string]interface{}{
		"method": "UNSUBSCRIBE",
		"params": []string{stream},
		"id":     time.Now().Unix(),
	}

	return c.SendJSON(unsubscribeReq)
}

// UnsubscribeList 取消订阅多个流
func (c *BnWebSocketClient) UnsubscribeList(streams []string) error {
	// 从订阅映射中移除
	for _, stream := range streams {
		c.MessageHandler.RemoveSubscription(stream)
	}

	// 构造币安取消订阅请求
	unsubscribeReq := map[string]interface{}{
		"method": "UNSUBSCRIBE",
		"params": streams,
		"id":     time.Now().Unix(),
	}

	return c.SendJSON(unsubscribeReq)
}

// SubscribeDepth 订阅深度数据
func (c *BnWebSocketClient) SubscribeDepth(symbol string, level string, listener OnReceive) error {
	stream := fmt.Sprintf("%s@depth%s", strings.ToLower(symbol), level)
	return c.Subscribe(stream, listener)
}

// SubscribeTicker 订阅24小时价格变动统计
func (c *BnWebSocketClient) SubscribeMiniTicker(symbol string, listener OnReceive) error {
	stream := fmt.Sprintf("%s@miniTicker", strings.ToLower(symbol))
	return c.Subscribe(stream, listener)
}

// SubscribeTicker 订阅全市场最新标记价格
func (c *BnWebSocketClient) SubscribeMarkPriceAll(listener OnReceive) error {
	return c.Subscribe("!markPrice@arr@1s", listener)
}

// SubscribeTicker 订阅全场所有交易对 24小时价格变动统计
func (c *BnWebSocketClient) SubscribeMiniTickerAll(listener OnReceive) error {
	return c.Subscribe("!miniTicker@arr", listener)
}

// SubscribeKline 订阅K线数据
func (c *BnWebSocketClient) SubscribeKline(symbol string, interval string, listener OnReceive) error {
	stream := fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	return c.Subscribe(stream, listener)
}

// SubscribeTrade 订阅成交数据
func (c *BnWebSocketClient) SubscribeTrade(symbol string, listener OnReceive) error {
	stream := fmt.Sprintf("%s@trade", strings.ToLower(symbol))
	return c.Subscribe(stream, listener)
}

// SubscribeBookTicker 订阅最优挂单数据
func (c *BnWebSocketClient) SubscribeBookTicker(symbol string, listener OnReceive) error {
	stream := fmt.Sprintf("%s@bookTicker", strings.ToLower(symbol))
	return c.Subscribe(stream, listener)
}

// SubscribeAllTicker 订阅所有产品的24小时价格变动统计
func (c *BnWebSocketClient) SubscribeAllTicker(listener OnReceive) error {
	stream := "!ticker@arr"
	return c.Subscribe(stream, listener)
}

// SubscribeAllBookTicker 订阅所有产品的最优挂单数据
func (c *BnWebSocketClient) SubscribeAllBookTicker(listener OnReceive) error {
	stream := "!bookTicker"
	return c.Subscribe(stream, listener)
}

// SetListeners 设置监听器
func (c *BnWebSocketClient) SetListeners(msgListener OnReceive, errorListener OnReceive) {
	c.MessageHandler.SetListeners(msgListener, errorListener)
}

// IsLoggedIn 检查是否已登录（币安现货不需要登录）
func (c *BnWebSocketClient) IsLoggedIn() bool {
	return c.MessageHandler.IsLoggedIn()
}

// Login 登录（币安现货不需要登录，但为了兼容性保留）
func (h *BnMessageHandler) Login() error {
	h.mu.Lock()
	h.LoginStatus = true
	h.mu.Unlock()
	return nil
}

// SendPing 发送ping消息（币安使用json格式）
func (c *BnWebSocketClient) SendPing() error {
	return c.Send("ping")
}

// GetAllSubscriptions 获取所有订阅
func (c *BnWebSocketClient) GetAllSubscriptions() []string {
	c.MessageHandler.mu.RLock()
	defer c.MessageHandler.mu.RUnlock()

	subscriptions := make([]string, 0, len(c.MessageHandler.StreamMap))
	for stream := range c.MessageHandler.StreamMap {
		subscriptions = append(subscriptions, stream)
	}
	return subscriptions
}
