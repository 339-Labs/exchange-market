package maps

import "sync"

type PriceData struct {
	Symbol      string `json:"symbol"`
	Price       string `json:"price"`
	Timestamp   string `json:"timestamp"`
	FundingRate string `json:"funding_rate"`
}

type PriceMap struct {
	mu   sync.RWMutex
	data map[string]*PriceData // 用指针节省拷贝开销
}

func NewPriceMap() *PriceMap {
	return &PriceMap{
		data: make(map[string]*PriceData),
	}
}

func (p *PriceMap) Write(key string, value *PriceData) {
	p.mu.Lock()
	p.data[key] = value
	p.mu.Unlock()
}

func (p *PriceMap) Read(key string) (*PriceData, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.data[key]
	return val, ok
}
