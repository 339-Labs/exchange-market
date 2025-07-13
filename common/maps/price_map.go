package maps

import (
	"sync"
	"sync/atomic"
	"time"
)

type PriceData struct {
	Symbol      string `json:"symbol"`
	Price       string `json:"price"`
	FundingRate string `json:"funding_rate"`
	MarkPrice   string `json:"mark_price"`
	Timestamp   string `json:"timestamp"`
}

type PriceMap struct {
	mu   sync.RWMutex
	data map[string]*PriceData // 用指针节省拷贝开销

	// 双buffer
	writeBuffer *sync.Map
	readBuffer  *sync.Map
	bufferMu    sync.Mutex

	// 触发条件
	maxBatchSize int
	maxWaitTime  time.Duration
	counter      int64

	// 控制通道
	done chan struct{}
}

func NewPriceMap(max int) *PriceMap {
	return &PriceMap{
		data: make(map[string]*PriceData, max),
	}
}

// 创建新的双buffer价格映射
func NewDoubleBufferPriceMap(maxBatchSize int, maxWaitTime time.Duration) *PriceMap {
	return &PriceMap{
		data:         make(map[string]*PriceData),
		writeBuffer:  &sync.Map{},
		readBuffer:   &sync.Map{},
		maxBatchSize: maxBatchSize,
		maxWaitTime:  maxWaitTime,
		done:         make(chan struct{}),
	}
}

func (p *PriceMap) Write(key string, value *PriceData) {
	p.mu.Lock()
	p.data[key] = value
	p.mu.Unlock()
}

func (p *PriceMap) WriteBatch(data map[string]*PriceData) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for key, value := range data {
		p.data[key] = value
	}
}

func (p *PriceMap) Read(key string) (*PriceData, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.data[key]
	return val, ok
}

// 读取所有键值对
func (p *PriceMap) ReadAll() map[string]*PriceData {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*PriceData, len(p.data))
	for k, v := range p.data {
		result[k] = v
	}
	return result
}

// 获取所有键
func (p *PriceMap) ReadAllKeys() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	keys := make([]string, 0, len(p.data))
	for k := range p.data {
		keys = append(keys, k)
	}
	return keys
}

// 获取数据数量
func (p *PriceMap) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.data)
}

// 双buffer交换机制 写入数据
func (p *PriceMap) DoubleBufferWriteBatch(key string, value *PriceData) {
	p.writeBuffer.Store(key, value)

	// 原子增加计数
	count := atomic.AddInt64(&p.counter, 1)

	// 达到批量大小时触发刷新
	if count >= int64(p.maxBatchSize) {
		p.triggerFlush()
	}
}

// 触发刷新
func (p *PriceMap) triggerFlush() {
	p.bufferMu.Lock()
	defer p.bufferMu.Unlock()

	// 只有在有数据时才刷新
	if atomic.LoadInt64(&p.counter) == 0 {
		return
	}

	// 交换buffer
	p.writeBuffer, p.readBuffer = p.readBuffer, p.writeBuffer
	atomic.StoreInt64(&p.counter, 0)

	// 异步刷新
	go p.flushReadBuffer()
}

// 刷新读buffer到主数据
func (p *PriceMap) flushReadBuffer() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.readBuffer.Range(func(key, value interface{}) bool {
		p.data[key.(string)] = value.(*PriceData)
		p.readBuffer.Delete(key)
		return true
	})
}

// 双buffer交换机制 启动定期刷新goroutine
func (p *PriceMap) StartPeriodicFlush() {
	go func() {
		ticker := time.NewTicker(p.maxWaitTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if atomic.LoadInt64(&p.counter) > 0 {
					p.triggerFlush()
				}
			case <-p.done:
				return
			}
		}
	}()
}

// 停止所有goroutine
func (p *PriceMap) Stop() {
	close(p.done)
	// 最后刷新一次
	p.triggerFlush()
}
