package sse

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(app *fiber.App, hub *Hub) {
	h := &Handler{hub: hub}

	// client connect ke room auction_xxx
	app.Get("/sse/auction/:id", h.Subscribe)
}

func (h *Handler) Subscribe(c *fiber.Ctx) error {
	auctionID := c.Params("id")
	room := fmt.Sprintf("auction_%s", auctionID)

	client := make(chan []byte, 10)

	// register client
	h.hub.Register <- Subscription{
		Room:   room,
		Client: client,
	}

	// ambil context fiber biar aman di goroutine
	ctx := c.Context()

	// unregister kalau koneksi mati
	go func() {
		<-ctx.Done()
		h.hub.Unregister <- Subscription{
			Room:   room,
			Client: client,
		}
	}()

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// kirim event initial
		fmt.Fprintf(w, "event: %s\n", "initial")
		fmt.Fprintf(w, "data: connected to %s\n\n", room)
		w.Flush()

		for {
			select {
			case msg, ok := <-client:
				if !ok {
					return
				}
				fmt.Fprintf(w, "event: %s\n", "bid-updated")
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.Flush()
			case <-ctx.Done():
				return
			}
		}
	})

	return nil
}
