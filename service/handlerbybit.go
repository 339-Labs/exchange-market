package service

import (
	"context"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
	"github.com/339-Labs/exchange-market/exchange/cex/bybit"
	"github.com/339-Labs/exchange-market/redis"
	"github.com/339-Labs/exchange-market/worker"
	"github.com/ethereum/go-ethereum/log"
	"sync/atomic"
	"time"
)

type HandlerByBit struct {
	ByBitExClient *bybit.ByBitExClient
	ByBitTask     *worker.ByBitTask

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewHandlerByBit(config *config.Config, db *database.DB, redis *redis.RedisClient, shutdown context.CancelCauseFunc) (*HandlerByBit, error) {

	spotPriceMap := maps.NewPriceMap(10)
	featurePriceMap := maps.NewPriceMap(10)

	bybitExClient, _ := bybit.NewByBitExClient(&config.ExchangeConfig.ByBit, spotPriceMap, featurePriceMap)
	bitGetTask, _ := worker.NewByBitTask(shutdown, time.Second*1, spotPriceMap, featurePriceMap)

	return &HandlerByBit{
		ByBitExClient: bybitExClient,
		ByBitTask:     bitGetTask,
		shutdown:      shutdown,
	}, nil
}

func (h *HandlerByBit) Start(ctx context.Context) error {
	h.ByBitExClient.ExecuteSpotWs()
	h.ByBitExClient.ExecuteFeatureWs()
	h.ByBitTask.Start()
	return nil
}

func (h *HandlerByBit) Stop(ctx context.Context) error {
	h.ByBitTask.Close()
	h.ByBitExClient.ByBitWebSocketClient.Stop()
	log.Info("stop notify success")
	return nil
}

func (h *HandlerByBit) Stopped() bool {
	return h.stopped.Load()
}
