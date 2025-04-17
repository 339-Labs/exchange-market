package bitget

/*
* ProductType
USDT-FUTURES USDT专业合约
COIN-FUTURES 混合合约
USDC-FUTURES USDC专业合约
SUSDT-FUTURES USDT专业合约模拟盘
SCOIN-FUTURES 混合合约模拟盘
SUSDC-FUTURES USDC专业合约模拟盘
*/
type ProductType string

const (
	USDT  ProductType = "USDT-FUTURES"
	COIN  ProductType = "COIN-FUTURES"
	USDC  ProductType = "USDC-FUTURES"
	SUSDT ProductType = "SUSDT-FUTURES"
	SCOIN ProductType = "SCOIN-FUTURES"
	SUSDC ProductType = "SUSDC-FUTURES"
)
