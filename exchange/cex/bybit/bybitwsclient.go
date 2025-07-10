package bybit

import (
	"encoding/json"
	"fmt"
	"github.com/339-Labs/exchange-market/common"
	"github.com/339-Labs/exchange-market/common/signer"
	"github.com/339-Labs/exchange-market/common/ws"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/bybit/constants"
	"github.com/339-Labs/exchange-market/exchange/cex/bybit/model"
	"github.com/ethereum/go-ethereum/log"
	"sync"
	"time"
)

// BybitMessageHandler bybit消息处理器
type BybitMessageHandler struct {
	// 配置和认证
	Config      *config.CexExchangeConfig
	Signer      *signer.Signer
	LoginStatus bool
	NeedLogin   bool

	// 消息处理
	Listener      OnReceive
	ErrorListener OnReceive
	ScribeMap     map[string]OnReceive
	AllSubscribe  *model.Set

	// 同步
	mu sync.RWMutex

	// WebSocket客户端引用
	wsClient *ws.GenericWebSocketClient
}

// OnReceive 消息接收回调函数类型
type OnReceive func(message string)

// NewByBitMessageHandler 创建新的bybit消息处理器
func NewByBitMessageHandler(config *config.CexExchangeConfig, needLogin bool) *BybitMessageHandler {
	handler := &BybitMessageHandler{
		Config:       config,
		NeedLogin:    needLogin,
		LoginStatus:  false,
		ScribeMap:    make(map[string]OnReceive),
		AllSubscribe: model.NewSet(),
		Signer:       new(signer.Signer).Init(config.ApiSecretKey),
	}

	return handler
}

// SetWebSocketClient 设置WebSocket客户端引用
func (h *BybitMessageHandler) SetWebSocketClient(client *ws.GenericWebSocketClient) {
	h.wsClient = client
}

// SetListeners 设置消息监听器
func (h *BybitMessageHandler) SetListeners(msgListener OnReceive, errorListener OnReceive) {
	h.Listener = msgListener
	h.ErrorListener = errorListener
}

// HandleMessage 处理普通消息
func (h *BybitMessageHandler) HandleMessage(message string) error {
	jsonMap := common.JSONToMap(message)

	// 检查是否有错误代码
	if success, exists := jsonMap["success"]; exists {
		if !success.(bool) {
			return fmt.Errorf("received error msg: %d", message)
		}
	}

	// 处理登录响应
	if event, exists := jsonMap["op"]; exists && event == "auth" {
		return h.handleLoginResponse(message)
	}

	// 处理数据消息
	if _, exists := jsonMap["data"]; exists {
		return h.handleDataMessage(message, jsonMap)
	}

	// 处理其他消息
	return h.handleOtherMessage(message)
}

// HandleError 处理错误消息
func (h *BybitMessageHandler) HandleError(message string) error {
	log.Error("Received error message: %s", message)

	if h.ErrorListener != nil {
		h.ErrorListener(message)
	}

	return nil
}

// HandleSpecialMessage 处理特殊消息（如pong）
func (h *BybitMessageHandler) HandleSpecialMessage(message string) (handled bool, err error) {
	if message == "pong" {
		log.Info("Received pong: %s", message)
		return true, nil
	}

	return false, nil
}

// handleLoginResponse 处理登录响应
func (h *BybitMessageHandler) handleLoginResponse(message string) error {
	log.Info("Login response: %s", message)

	h.mu.Lock()
	h.LoginStatus = true
	h.mu.Unlock()

	if h.Listener != nil {
		h.Listener(message)
	}

	return nil
}

// handleDataMessage 处理数据消息
func (h *BybitMessageHandler) handleDataMessage(message string, jsonMap map[string]interface{}) error {
	listener := h.getListener(jsonMap["topic"])
	if listener != nil {
		listener(message)
	}

	return nil
}

// handleOtherMessage 处理其他消息
func (h *BybitMessageHandler) handleOtherMessage(message string) error {
	log.Info("Received other message: %s", message)

	if h.Listener != nil {
		h.Listener(message)
	}

	return nil
}

// getListener 获取特定订阅的监听器
func (h *BybitMessageHandler) getListener(topic interface{}) OnReceive {
	if topic == nil {
		return h.Listener
	}

	h.mu.RLock()
	listener, exists := h.ScribeMap[topic.(string)]
	h.mu.RUnlock()

	if !exists {
		return h.Listener
	}

	return listener
}

// IsLoggedIn 检查是否已登录
func (h *BybitMessageHandler) IsLoggedIn() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.LoginStatus
}

// AddSubscription 添加订阅
func (h *BybitMessageHandler) AddSubscription(req string, listener OnReceive) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.ScribeMap[req] = listener
	h.AllSubscribe.Add(req)
}

// RemoveSubscription 移除订阅
func (h *BybitMessageHandler) RemoveSubscription(req string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.ScribeMap, req)
	h.AllSubscribe.Remove(req)
}

// ByBitWebSocketClient Okx WebSocket客户端
type ByBitWebSocketClient struct {
	*ws.GenericWebSocketClient
	MessageHandler *BybitMessageHandler
}

// NewByBitWebSocketClient 创建新的bybit WebSocket客户端
func NewByBitWebSocketClient(config *config.CexExchangeConfig, needLogin bool) *ByBitWebSocketClient {
	// 创建WebSocket配置
	wsConfig := &ws.ConnectionConfig{
		WsUrl:               config.WsUrl,
		PingInterval:        15 * time.Second,
		ReconnectWaitSecond: float64(constants.ReconnectWaitSecond),
		TimerIntervalSecond: constants.TimerIntervalSecond * time.Second,
		EnableAutoReconnect: true,
		EnablePing:          true,
		PingMsg:             string(json.RawMessage(`{"op":"ping"}`)),
	}

	// 创建通用WebSocket客户端
	genericClient := ws.NewGenericWebSocketClient(wsConfig)

	// 创建Okx消息处理器
	messageHandler := NewByBitMessageHandler(config, needLogin)
	messageHandler.SetWebSocketClient(genericClient)

	// 设置消息处理器
	genericClient.SetMessageHandler(messageHandler)

	// 设置回调函数
	genericClient.SetCallbacks(
		func() {
			log.Info("Okx WebSocket connected")
			// 如果需要登录，则自动登录
			if needLogin {
				messageHandler.Login()
			}
		},
		func() {
			log.Info("Okx WebSocket disconnected")
		},
		func(attempt int) {
			log.Info("Okx WebSocket reconnecting, attempt: %d", attempt)
		},
	)

	return &ByBitWebSocketClient{
		GenericWebSocketClient: genericClient,
		MessageHandler:         messageHandler,
	}
}

// Login 登录
func (h *BybitMessageHandler) Login() error {
	if h.wsClient == nil {
		return fmt.Errorf("WebSocket client is not set")
	}

	timesStamp := common.TimesStampSec() + common.TimesStampSec()

	sign := h.Signer.ByBitSign(h.Config.ApiSecretKey, timesStamp, "")

	var args []interface{}
	args = append(args, h.Config.ApiKey)
	args = append(args, timesStamp)
	args = append(args, sign)

	baseReq := model.WsBaseReq{
		Op:   constants.WsOpLogin,
		Args: args,
	}

	return h.wsClient.SendJSON(baseReq)
}

// Subscribe 订阅
func (c *ByBitWebSocketClient) Subscribe(req string, listener OnReceive) error {
	stream := fmt.Sprintf("tickers.%s", req)
	// 添加到订阅映射
	c.MessageHandler.AddSubscription(stream, listener)

	// 发送订阅请求
	var args []interface{}
	args = append(args, stream)

	baseReq := model.WsBaseReq{
		Op:   constants.WsOpSubscribe,
		Args: args,
	}

	return c.SendJSON(baseReq)
}

// SubscribeList 订阅列表
func (c *ByBitWebSocketClient) SubscribeList(reqs []string, listener OnReceive) error {

	var args []interface{}
	for _, req := range reqs {

		stream := fmt.Sprintf("tickers.%s", req)
		// 添加到订阅映射
		c.MessageHandler.AddSubscription(stream, listener)

		// 发送订阅请求
		args = append(args, stream)
	}

	baseReq := model.WsBaseReq{
		Op:   constants.WsOpSubscribe,
		Args: args,
	}

	return c.SendJSON(baseReq)
}

// Unsubscribe 取消订阅
func (c *ByBitWebSocketClient) Unsubscribe(req string) error {
	stream := fmt.Sprintf("tickers.%s", req)
	// 从订阅映射中移除
	c.MessageHandler.RemoveSubscription(stream)

	// 发送取消订阅请求
	var args []interface{}
	args = append(args, req)

	baseReq := model.WsBaseReq{
		Op:   constants.WsOpUnsubscribe,
		Args: args,
	}

	return c.SendJSON(baseReq)
}

// UnsubscribeList  取消订阅
func (c *ByBitWebSocketClient) UnsubscribeList(req []string) error {
	stream := fmt.Sprintf("tickers.%s", req)
	var args []interface{}
	for _, req := range req {
		// 从订阅映射中移除
		c.MessageHandler.RemoveSubscription(stream)
		// 发送取消订阅请求
		args = append(args, req)
	}

	baseReq := model.WsBaseReq{
		Op:   constants.WsOpUnsubscribe,
		Args: args,
	}

	return c.SendJSON(baseReq)
}

// SetListeners 设置监听器
func (c *ByBitWebSocketClient) SetListeners(msgListener OnReceive, errorListener OnReceive) {
	c.MessageHandler.SetListeners(msgListener, errorListener)
}

// IsLoggedIn 检查是否已登录
func (c *ByBitWebSocketClient) IsLoggedIn() bool {
	return c.MessageHandler.IsLoggedIn()
}
