package ws

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"sync"
	"time"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(message string) error
	HandleError(message string) error
	HandleSpecialMessage(message string) (handled bool, err error)
}

// ConnectionConfig WebSocket连接配置
type ConnectionConfig struct {
	WsUrl                string
	PingInterval         time.Duration // ping间隔
	PingMsg              string
	ReconnectWaitSecond  float64       // 重连等待时间
	TimerIntervalSecond  time.Duration // 定时器间隔
	EnableAutoReconnect  bool          // 是否启用自动重连
	EnablePing           bool          // 是否启用ping
	MaxReconnectAttempts int           // 最大重连尝试次数
}

// DefaultConnectionConfig 默认配置
func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		WsUrl:                "",
		PingInterval:         15 * time.Second,
		ReconnectWaitSecond:  30.0,
		TimerIntervalSecond:  1 * time.Second,
		EnableAutoReconnect:  true,
		EnablePing:           true,
		MaxReconnectAttempts: 5,
	}
}

// GenericWebSocketClient 通用WebSocket客户端
type GenericWebSocketClient struct {
	// 连接相关
	Connection      bool
	WebSocketClient *websocket.Conn
	SendMutex       *sync.Mutex
	Config          *ConnectionConfig

	// 消息处理
	MessageHandler MessageHandler

	// 定时器和状态
	Ticker           *time.Ticker
	PingCron         *cron.Cron
	LastReceivedTime time.Time

	// 回调函数
	OnConnected    func()
	OnDisconnected func()
	OnReconnecting func(attempt int)

	// 控制
	stopChan       chan struct{}
	reconnectCount int
	isRunning      bool
	mu             sync.RWMutex
}

// NewGenericWebSocketClient 创建新的通用WebSocket客户端
func NewGenericWebSocketClient(config *ConnectionConfig) *GenericWebSocketClient {
	if config == nil {
		config = DefaultConnectionConfig()
	}

	return &GenericWebSocketClient{
		Connection:       false,
		SendMutex:        &sync.Mutex{},
		Config:           config,
		Ticker:           time.NewTicker(config.TimerIntervalSecond),
		LastReceivedTime: time.Now(),
		stopChan:         make(chan struct{}),
		reconnectCount:   0,
		isRunning:        false,
	}
}

// SetMessageHandler 设置消息处理器
func (c *GenericWebSocketClient) SetMessageHandler(handler MessageHandler) {
	c.MessageHandler = handler
}

// SetCallbacks 设置回调函数
func (c *GenericWebSocketClient) SetCallbacks(onConnected, onDisconnected func(), onReconnecting func(int)) {
	c.OnConnected = onConnected
	c.OnDisconnected = onDisconnected
	c.OnReconnecting = onReconnecting
}

// Start 启动WebSocket客户端
func (c *GenericWebSocketClient) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning {
		return fmt.Errorf("client is already running")
	}

	if c.MessageHandler == nil {
		return fmt.Errorf("message handler is not set")
	}

	err := c.connect()
	if err != nil {
		return err
	}

	c.isRunning = true

	// 启动读取循环
	go c.readLoop()

	// 启动定时器循环（用于重连检查）
	if c.Config.EnableAutoReconnect {
		go c.timerLoop()
	}

	// 启动ping循环
	if c.Config.EnablePing {
		c.startPing()
	}

	return nil
}

// Stop 停止WebSocket客户端
func (c *GenericWebSocketClient) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isRunning {
		return nil
	}

	c.isRunning = false
	close(c.stopChan)

	// 停止ping
	if c.PingCron != nil {
		c.PingCron.Stop()
	}

	// 停止定时器
	if c.Ticker != nil {
		c.Ticker.Stop()
	}

	// 关闭WebSocket连接
	return c.disconnect()
}

// connect 连接到WebSocket服务器
func (c *GenericWebSocketClient) connect() error {
	log.Info("WebSocket connecting to %s...", c.Config.WsUrl)

	var err error
	c.WebSocketClient, _, err = websocket.DefaultDialer.Dial(c.Config.WsUrl, nil)
	if err != nil {
		log.Error("WebSocket connection failed: %s", err)
		return err
	}

	c.Connection = true
	c.LastReceivedTime = time.Now()
	c.reconnectCount = 0

	log.Info("WebSocket connected successfully")

	if c.OnConnected != nil {
		c.OnConnected()
	}

	return nil
}

// disconnect 断开WebSocket连接
func (c *GenericWebSocketClient) disconnect() error {
	if c.WebSocketClient == nil {
		return nil
	}

	log.Info("WebSocket disconnecting...")

	err := c.WebSocketClient.Close()
	c.WebSocketClient = nil
	c.Connection = false

	if err != nil {
		log.Error("WebSocket disconnect error: %s", err)
	} else {
		log.Info("WebSocket disconnected")
	}

	if c.OnDisconnected != nil {
		c.OnDisconnected()
	}

	return err
}

// reconnect 重连
func (c *GenericWebSocketClient) reconnect() {
	c.reconnectCount++

	if c.Config.MaxReconnectAttempts > 0 && c.reconnectCount > c.Config.MaxReconnectAttempts {
		log.Error("Max reconnection attempts reached (%d)", c.Config.MaxReconnectAttempts)
		return
	}

	log.Info("WebSocket reconnecting... (attempt %d)", c.reconnectCount)

	if c.OnReconnecting != nil {
		c.OnReconnecting(c.reconnectCount)
	}

	c.disconnect()

	// 等待一段时间后重连
	time.Sleep(time.Duration(c.Config.ReconnectWaitSecond) * time.Second)

	err := c.connect()
	if err != nil {
		log.Error("Reconnection failed: %s", err)
	}
}

// Send 发送消息
func (c *GenericWebSocketClient) Send(data string) error {
	if c.WebSocketClient == nil {
		return fmt.Errorf("no connection available")
	}

	log.Info("Sending message: %s", data)

	c.SendMutex.Lock()
	defer c.SendMutex.Unlock()

	err := c.WebSocketClient.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		log.Error("Failed to send message: %s, error: %s", data, err)
		return err
	}

	return nil
}

// SendJSON 发送JSON消息
func (c *GenericWebSocketClient) SendJSON(data interface{}) error {
	if c.WebSocketClient == nil {
		return fmt.Errorf("no connection available")
	}

	c.SendMutex.Lock()
	defer c.SendMutex.Unlock()

	err := c.WebSocketClient.WriteJSON(data)
	if err != nil {
		log.Error("Failed to send JSON message: %v, error: %s", data, err)
		return err
	}

	return nil
}

// readLoop 读取循环
func (c *GenericWebSocketClient) readLoop() {
	for {
		select {
		case <-c.stopChan:
			return
		default:
			if c.WebSocketClient == nil {
				log.Info("Read error: no connection available")
				time.Sleep(c.Config.TimerIntervalSecond)
				continue
			}

			_, buf, err := c.WebSocketClient.ReadMessage()
			if err != nil {
				log.Info("Read error: %s", err)
				continue
			}

			c.LastReceivedTime = time.Now()
			message := string(buf)

			log.Info("Received message: %s", message)

			// 首先检查是否是特殊消息（如pong）
			if handled, err := c.MessageHandler.HandleSpecialMessage(message); err != nil {
				log.Error("Error handling special message: %s", err)
			} else if handled {
				continue
			}

			// 处理普通消息
			if err := c.MessageHandler.HandleMessage(message); err != nil {
				log.Error("Error handling message: %s", err)
				c.MessageHandler.HandleError(message)
			}
		}
	}
}

// timerLoop 定时器循环，用于检查连接状态和重连
func (c *GenericWebSocketClient) timerLoop() {
	for {
		select {
		case <-c.stopChan:
			return
		case <-c.Ticker.C:
			if !c.isRunning {
				return
			}

			elapsedSecond := time.Since(c.LastReceivedTime).Seconds()

			if elapsedSecond > c.Config.ReconnectWaitSecond {
				log.Info("Connection timeout, reconnecting...")
				go c.reconnect()
			}
		}
	}
}

// startPing 启动ping循环
func (c *GenericWebSocketClient) startPing() {
	c.PingCron = cron.New()

	// 计算cron表达式
	interval := int(c.Config.PingInterval.Seconds())
	cronExpr := fmt.Sprintf("*/%d * * * * *", interval)

	_ = c.PingCron.AddFunc(cronExpr, c.ping)
	c.PingCron.Start()
}

// ping 发送ping消息
func (c *GenericWebSocketClient) ping() {
	ping := "ping"
	if c.Config.PingMsg != "" {
		ping = c.Config.PingMsg
	}
	if err := c.Send(ping); err != nil {
		log.Error("Failed to send ping: %s", err)
	}
}

// IsConnected 检查是否连接
func (c *GenericWebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Connection
}

// IsRunning 检查是否运行中
func (c *GenericWebSocketClient) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isRunning
}
