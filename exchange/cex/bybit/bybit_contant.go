package bybit

/*
*
统一帐户
	spot   	现货
	linear USDT  永续, USDT交割, USDC永续, USDC交割
	inverse    反向合约，包含反向永续, 反向交割
	option    期权
经典帐户
	spot 	现货
	linear USDT	永续
	inverse 	反向合约，包含反向永续, 反向交割
	option 	期权
*/

const (
	Spot    = "Spot"
	Linear  = "linear"
	Inverse = "inverse"
	Option  = "option"
)
