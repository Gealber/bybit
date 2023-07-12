package websocket

import (
	"context"
	"errors"
	"log"
	"os"
)

type TickersHandler struct {
	logger *log.Logger
}

func NewTickersHandler() *TickersHandler {
	return &TickersHandler{
		logger: log.New(os.Stdout, "[ticket-handler]", log.Lshortfile),
	}
}

func (t *TickersHandler) ProcessMsg(ctx context.Context, obj any) error {
	msg, ok := obj.(TickersResponse)
	if !ok {
		return errors.New("invalid type of obj for TickersResponse")
	}

	t.logger.Printf("MSG: %+v\n", msg)

	return nil
}
