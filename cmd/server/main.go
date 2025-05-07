package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdhityaRamadhanus/zdb/tcp"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	srv := tcp.NewServer("tcp", "localhost:9000")
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-termChan
		cancel()
	}()

	log.Error().Err(srv.Run(ctx)).Msg("Shutdown server")
}
