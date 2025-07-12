package constants

const (
	// 币安WebSocket消息类型
	StatusOK            = 200
	Spot                = "Spot"
	Feature             = "Feature"
	EventTicker         = "24hrTicker"        // 交易对详细信息
	EventMiniTicker     = "24hrMiniTicker"    // 交易对精简信息
	EventMarkPrice      = "markPriceUpdate"   // 交易对标记价格
	StreamTickerArr     = "!ticker@arr"       // 交易对详细信息 - 订阅所有交易对
	StreamMiniTickerArr = "!miniTicker@arr"   // 交易对精简信息 - 订阅所有交易对
	StreamMarkPriceArr  = "!markPrice@arr@1s" // 交易对标记价格 - 订阅所有交易对

	/*
	 * http headers
	 */
	ContentType        = "Content-Type"
	BgAccessKey        = "ACCESS-KEY"
	BgAccessSign       = "ACCESS-SIGN"
	BgAccessTimestamp  = "ACCESS-TIMESTAMP"
	BgAccessPassphrase = "ACCESS-PASSPHRASE"
	ApplicationJson    = "application/json"

	EN_US  = "en_US"
	ZH_CN  = "zh_CN"
	LOCALE = "locale="

	/*
	 * http methods
	 */
	GET  = "GET"
	POST = "POST"

	/*
	 * websocket
	 */
	WsAuthMethod        = "GET"
	WsAuthPath          = "/users/self/verify"
	WsOpLogin           = "login"
	WsOpUnsubscribe     = "unsubscribe"
	WsOpSubscribe       = "subscribe"
	TimerIntervalSecond = 5
	ReconnectWaitSecond = 60

	/*
	 * SignType
	 */
	RSA    = "RSA"
	SHA256 = "SHA256"
)
