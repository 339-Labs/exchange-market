package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/common/tasks"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type OkxtTask struct {
	spotPriceMap    *maps.PriceMap
	featurePriceMap *maps.PriceMap
	resourceCtx     context.Context

	resourceCancel context.CancelFunc
	tasks          tasks.Group
	ticker         *time.Ticker
}

func NewOkxTask(cfg *context.Context, shutdown context.CancelCauseFunc, duration time.Duration, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) (*OkxtTask, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &OkxtTask{
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("okx ws error "))
		}},
		ticker:          time.NewTicker(duration),
		spotPriceMap:    spotPriceMap,
		featurePriceMap: featurePriceMap,
	}, nil
}

func (t *OkxtTask) Start() error {
	log.Info("okx task started")
	t.tasks.Go(func() error {
		for {

			select {

			case <-t.ticker.C:
				// todo  okx ws data handler, spot and feature

			case <-t.resourceCtx.Done():
				log.Info("stop okx task in work")
				return nil

			}

		}
	})
	return nil
}

func (t *OkxtTask) Close() error {
	var result error
	t.resourceCancel()
	t.ticker.Stop()
	log.Info("okx task stopped")
	if err := t.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("okx ws task wait error"))
	}
	log.Info("okx task stopped success")
	return nil
}
