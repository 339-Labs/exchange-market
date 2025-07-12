package main

import (
	"context"
	"fmt"
	"github.com/339-Labs/exchange-market/common/cliapp"
	"github.com/339-Labs/exchange-market/common/opio"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
	flags2 "github.com/339-Labs/exchange-market/flags"
	"github.com/339-Labs/exchange-market/redis"
	"github.com/339-Labs/exchange-market/service"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

func NewCli(GitCommit string, GitData string) *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              GitCommit,
		Description:          "exchange market data",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "migrate",
				Description: fmt.Sprintf("migrate the database to the latest version"),
				Flags:       flags,
				Action:      runMigrations,
			},
			{
				Name:        "run bn",
				Description: fmt.Sprintf("run bn task"),
				Flags:       flags,
				Action:      cliapp.LifecycleCmd(runBnTask),
			},
			{
				Name:        "run okx",
				Description: fmt.Sprintf("run okx task"),
				Flags:       flags,
				Action:      cliapp.LifecycleCmd(runOkxTask),
			},
			{
				Name:        "run bybit",
				Description: fmt.Sprintf("run bybit task"),
				Flags:       flags,
				Action:      cliapp.LifecycleCmd(runBybitTask),
			},
			{
				Name:        "run bitget",
				Description: fmt.Sprintf("run bitget task"),
				Flags:       flags,
				Action:      cliapp.LifecycleCmd(runBitgetTask),
			},
		},
	}
}

func runMigrations(ctx *cli.Context) error {
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	log.Info("running migrations...")
	config, err := config.NewConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return err
	}

	db, err := database.NewDB(&config.SlaveDBConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Error("fail to close database", "err", err)
		}
	}(db)
	return db.ExecuteSQLMigration(config.Migrations)
}

func runBnTask(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	config, err := config.NewConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	db, err := database.NewDB(&config.SlaveDBConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	redis, err := redis.NewRedisClient(config.RedisConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	return service.NewHandlerBN(config, db, redis, shutdown)
}

func runOkxTask(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	config, err := config.NewConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	db, err := database.NewDB(&config.SlaveDBConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	redis, err := redis.NewRedisClient(config.RedisConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	return service.NewHandlerOkx(config, db, redis, shutdown)
}

func runBybitTask(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {

	config, err := config.NewConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	db, err := database.NewDB(&config.SlaveDBConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	redis, err := redis.NewRedisClient(config.RedisConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	return service.NewHandlerByBit(config, db, redis, shutdown)
}

func runBitgetTask(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	config, err := config.NewConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	db, err := database.NewDB(&config.SlaveDBConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	redis, err := redis.NewRedisClient(config.RedisConfig)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	return service.NewHandlerBitGet(config, db, redis, shutdown)
}
