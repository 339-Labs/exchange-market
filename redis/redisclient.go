package redis

import (
	"context"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 线程安全的Redis客户端
type RedisClient struct {
	rdb    *redis.Client
	mu     sync.RWMutex
	pool   *redis.Ring // 可选：使用Redis集群
	closed bool
}

// RedisConfig Redis配置
type RedisConfig struct {
	Address      string
	Username     string
	Password     string
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
	MaxIdleConns int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewRedisClient 创建线程安全的Redis客户端
func NewRedisClient(config config.RedisConfig) (*RedisClient, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Address,
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxIdleConns: 50,
		DB:           0,
	})

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		rdb:    rdb,
		closed: false,
	}, nil
}

// Close 线程安全地关闭Redis连接
func (r *RedisClient) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true
	return r.rdb.Close()
}

// isClientClosed 检查客户端是否已关闭
func (r *RedisClient) isClientClosed() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.closed
}

// === 高并发行情数据操作方法 ===

// SetPriceData 并发安全地存储行情数据
func (r *RedisClient) SetPriceData(ctx context.Context, priceData *maps.PriceData) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	key := "market_data:" + priceData.Symbol

	// 使用HMSET原子性地设置所有字段
	data := map[string]interface{}{
		"symbol":       priceData.Symbol,
		"price":        priceData.Price,
		"funding_rate": priceData.FundingRate,
		"mark_price":   priceData.MarkPrice,
		"timestamp":    priceData.Timestamp,
	}

	return r.rdb.HMSet(ctx, key, data).Err()
}

// SetPriceDataWithTTL 并发安全地存储行情数据并设置过期时间
func (r *RedisClient) SetPriceDataWithTTL(ctx context.Context, priceData *maps.PriceData, ttl time.Duration) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	key := "market_data:" + priceData.Symbol

	// 使用管道确保原子性
	pipe := r.rdb.Pipeline()

	data := map[string]interface{}{
		"symbol":       priceData.Symbol,
		"price":        priceData.Price,
		"funding_rate": priceData.FundingRate,
		"mark_price":   priceData.MarkPrice,
		"timestamp":    priceData.Timestamp,
	}

	pipe.HMSet(ctx, key, data)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

// BatchSetPriceData 批量并发安全地存储多个行情数据
func (r *RedisClient) BatchSetPriceData(ctx context.Context, priceDataList []*maps.PriceData) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(priceDataList) == 0 {
		return nil
	}

	// 使用管道批量操作
	pipe := r.rdb.Pipeline()

	for _, priceData := range priceDataList {
		key := "market_data:" + priceData.Symbol
		data := map[string]interface{}{
			"symbol":       priceData.Symbol,
			"price":        priceData.Price,
			"funding_rate": priceData.FundingRate,
			"mark_price":   priceData.MarkPrice,
			"timestamp":    priceData.Timestamp,
		}
		pipe.HMSet(ctx, key, data)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// BatchSetPriceDataWithTTL 批量存储行情数据并设置过期时间
func (r *RedisClient) BatchSetPriceDataWithTTL(ctx context.Context, priceDataList []*maps.PriceData, ttl time.Duration) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(priceDataList) == 0 {
		return nil
	}

	// 使用管道批量操作
	pipe := r.rdb.Pipeline()

	for _, priceData := range priceDataList {
		key := "market_data:" + priceData.Symbol
		data := map[string]interface{}{
			"symbol":       priceData.Symbol,
			"price":        priceData.Price,
			"funding_rate": priceData.FundingRate,
			"mark_price":   priceData.MarkPrice,
			"timestamp":    priceData.Timestamp,
		}
		pipe.HMSet(ctx, key, data)
		pipe.Expire(ctx, key, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetPriceData 并发安全地获取行情数据
func (r *RedisClient) GetPriceData(ctx context.Context, symbol string) (*maps.PriceData, error) {
	if r.isClientClosed() {
		return nil, redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	key := "market_data:" + symbol

	result, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, redis.Nil
	}

	priceData := &maps.PriceData{
		Symbol:      result["symbol"],
		Price:       result["price"],
		FundingRate: result["funding_rate"],
		MarkPrice:   result["mark_price"],
		Timestamp:   result["timestamp"],
	}

	return priceData, nil
}

// GetMultiplePriceData 并发安全地批量获取多个交易对的行情数据
func (r *RedisClient) GetMultiplePriceData(ctx context.Context, symbols []string) (map[string]*maps.PriceData, error) {
	if r.isClientClosed() {
		return nil, redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(symbols) == 0 {
		return make(map[string]*maps.PriceData), nil
	}

	// 使用管道批量获取
	pipe := r.rdb.Pipeline()

	cmds := make(map[string]*redis.MapStringStringCmd)
	for _, symbol := range symbols {
		key := "market_data:" + symbol
		cmds[symbol] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*maps.PriceData)
	for symbol, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		if len(data) > 0 {
			result[symbol] = &maps.PriceData{
				Symbol:      data["symbol"],
				Price:       data["price"],
				FundingRate: data["funding_rate"],
				MarkPrice:   data["mark_price"],
				Timestamp:   data["timestamp"],
			}
		}
	}

	return result, nil
}

// UpdatePriceDataField 并发安全地更新单个字段
func (r *RedisClient) UpdatePriceDataField(ctx context.Context, symbol string, field string, value string) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	key := "market_data:" + symbol
	return r.rdb.HSet(ctx, key, field, value).Err()
}

// BatchUpdatePriceFields 批量更新多个交易对的字段
func (r *RedisClient) BatchUpdatePriceFields(ctx context.Context, updates map[string]map[string]interface{}) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(updates) == 0 {
		return nil
	}

	// 使用管道批量更新
	pipe := r.rdb.Pipeline()

	for symbol, fields := range updates {
		key := "market_data:" + symbol
		pipe.HMSet(ctx, key, fields)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// === 高并发辅助方法 ===

// GetPriceDataWithWorkers 使用worker池并发获取大量行情数据
func (r *RedisClient) GetPriceDataWithWorkers(ctx context.Context, symbols []string, workerCount int) (map[string]*maps.PriceData, error) {
	if r.isClientClosed() {
		return nil, redis.ErrClosed
	}

	if len(symbols) == 0 {
		return make(map[string]*maps.PriceData), nil
	}

	if workerCount <= 0 {
		workerCount = 10
	}

	// 创建channels
	symbolChan := make(chan string, len(symbols))
	resultChan := make(chan struct {
		symbol string
		data   *maps.PriceData
		err    error
	}, len(symbols))

	// 启动worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for symbol := range symbolChan {
				data, err := r.GetPriceData(ctx, symbol)
				resultChan <- struct {
					symbol string
					data   *maps.PriceData
					err    error
				}{symbol, data, err}
			}
		}()
	}

	// 发送任务
	go func() {
		for _, symbol := range symbols {
			symbolChan <- symbol
		}
		close(symbolChan)
	}()

	// 等待所有worker完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	result := make(map[string]*maps.PriceData)
	for res := range resultChan {
		if res.err != nil && res.err != redis.Nil {
			return nil, res.err
		}
		if res.data != nil {
			result[res.symbol] = res.data
		}
	}

	return result, nil
}

func (r *RedisClient) CachePrice(symbol string, price string, ttl time.Duration) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}
	key := "price:" + symbol
	return r.rdb.Set(context.Background(), key, price, ttl).Err()
}

func (r *RedisClient) GetPrice(symbol string) (string, error) {
	if r.isClientClosed() {
		return "", redis.ErrClosed
	}
	key := "price:" + symbol
	return r.rdb.Get(context.Background(), key).Result()
}

func (r *RedisClient) Set(key string, value string, ttl time.Duration) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}
	return r.rdb.Set(context.Background(), key, value, ttl).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	if r.isClientClosed() {
		return "", redis.ErrClosed
	}
	return r.rdb.Get(context.Background(), key).Result()
}

func (r *RedisClient) TryLock(lockKey string, value string, ttl time.Duration) (bool, error) {
	if r.isClientClosed() {
		return false, redis.ErrClosed
	}
	ok, err := r.rdb.SetNX(context.Background(), lockKey, value, ttl).Result()
	return ok, err
}

func (r *RedisClient) Unlock(lockKey string, value string) error {
	if r.isClientClosed() {
		return redis.ErrClosed
	}
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`
	return r.rdb.Eval(context.Background(), script, []string{lockKey}, value).Err()
}
