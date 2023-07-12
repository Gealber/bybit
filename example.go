package main

import (
	"context"
	"log"

	"github.com/Gealber/bybit/config"
	bybitHttp "github.com/Gealber/bybit/http"
	bybitWs "github.com/Gealber/bybit/websocket"
)

func main() {
	ctx := context.Background()
	cfg := config.Config()

	err := httpExample(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// WEBSOCKET EXAMPLE:
	// websocketExample(ctx, cfg)
}

func websocketExample(ctx context.Context, cfg *config.AppConfig) {
	wb := bybitWs.NewClient(cfg)

	tickerSubsciption := bybitWs.Request{
		Op: "subscribe",
		Args: []interface{}{
			bybitWs.TickersTONUSDTTopic,
		},
	}

	subscriptions := []bybitWs.Request{
		tickerSubsciption,
	}
	handlers := map[string]bybitWs.Handler{
		bybitWs.TickersTONUSDTTopic: bybitWs.NewTickersHandler(),
	}

	if err := wb.Run(ctx, subscriptions, handlers); err != nil {
		log.Panic(err)
	}
}

func httpExample(ctx context.Context, cfg *config.AppConfig) error {
	client, err := bybitHttp.New(cfg)
	if err != nil {
		return err
	}

	return client.PlaceCascadeOrders(bybitHttp.SellDirection, bybitHttp.TonUSDTSymbol, 0.0001, 4000)
}
