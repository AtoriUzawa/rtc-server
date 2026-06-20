package main

import (
	"log"
	"sync"

	"github.com/AtoriUzawa/vlink-server/internal/app"
)

func main() {
	app := app.New()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()

		if err := app.RunHTTP("0.0.0.0:18888"); err != nil {
			log.Fatalf("http server start failed: %v\n", err)
			return
		}
	}()

	go func() {
		defer wg.Done()

		if err := app.RunWS("0.0.0.0:18887"); err != nil {
			log.Fatalf("websocket server start failed: %v\n", err)
		}
	}()

	wg.Wait()
}
