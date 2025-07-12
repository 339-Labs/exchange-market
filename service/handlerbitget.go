package service

import (
	"context"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
	"github.com/339-Labs/exchange-market/exchange/cex/bitget"
	"github.com/339-Labs/exchange-market/redis"
	"github.com/339-Labs/exchange-market/worker"
	"github.com/ethereum/go-ethereum/log"
	"sync/atomic"
	"time"
)

type HandlerBitGet struct {
	BitGetExClient *bitget.BitGetExClient
	BitGetTask     *worker.BitGetTask

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewHandlerBitGet(config *config.Config, db *database.DB, redis *redis.RedisClient, shutdown context.CancelCauseFunc) (*HandlerBitGet, error) {

	spotPriceMap := maps.NewPriceMap(10)
	featurePriceMap := maps.NewPriceMap(10)

	bitGetExClient, _ := bitget.NewBitGetExClient(&config.ExchangeConfig.BitGet, spotPriceMap, featurePriceMap)
	bitGetTask, _ := worker.NewBitGetTask(shutdown, time.Second*1, spotPriceMap, featurePriceMap)

	return &HandlerBitGet{
		BitGetExClient: bitGetExClient,
		BitGetTask:     bitGetTask,
		shutdown:       shutdown,
	}, nil
}

func (h *HandlerBitGet) Start(ctx context.Context) error {
	h.BitGetExClient.ExecuteWs()
	h.BitGetTask.Start()
	return nil
}

func (h *HandlerBitGet) Stop(ctx context.Context) error {
	h.BitGetTask.Close()
	h.BitGetExClient.BitGetWebSocketClient.Stop()
	log.Info("stop notify success")
	return nil
}

func (h *HandlerBitGet) Stopped() bool {
	return h.stopped.Load()
}
