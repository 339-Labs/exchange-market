package service

import (
	"context"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
	"github.com/339-Labs/exchange-market/exchange/cex/okx"
	"github.com/339-Labs/exchange-market/redis"
	"github.com/339-Labs/exchange-market/worker"
	"github.com/ethereum/go-ethereum/log"
	"sync/atomic"
	"time"
)

type HandlerOkx struct {
	OkxExClient *okx.OkxExClient
	OkxtTask    *worker.OkxTask

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewHandlerOkx(config *config.Config, db *database.DB, redis *redis.RedisClient, shutdown context.CancelCauseFunc) (*HandlerOkx, error) {

	spotPriceMap := maps.NewPriceMap(10)
	featurePriceMap := maps.NewPriceMap(10)
	markPriceMap := maps.NewPriceMap(10)
	rateMap := maps.NewPriceMap(10)

	okxExClient, _ := okx.NewOkxExClient(&config.ExchangeConfig.ByBit, spotPriceMap, featurePriceMap, markPriceMap, rateMap)
	okxTask, _ := worker.NewOkxTask(shutdown, time.Second*1, spotPriceMap, featurePriceMap, markPriceMap, rateMap)

	return &HandlerOkx{
		OkxExClient: okxExClient,
		OkxtTask:    okxTask,
		shutdown:    shutdown,
	}, nil
}

func (h *HandlerOkx) Start(ctx context.Context) error {
	h.OkxExClient.ExecuteSpotWs()
	h.OkxExClient.ExecuteFeatureWs()
	h.OkxtTask.Start()
	return nil
}

func (h *HandlerOkx) Stop(ctx context.Context) error {
	h.OkxtTask.Close()
	h.OkxExClient.OkxWebSocketClient.Stop()
	log.Info("stop notify success")
	return nil
}

func (h *HandlerOkx) Stopped() bool {
	return h.stopped.Load()
}
