package main

import (
	"github.com/ethereum/go-ethereum/log"
	"os"
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true)))
}
