package handlers

import (
	"bufio"
	"fmt"
	"strconv"

	"crypto-member/sse"
	"github.com/gofiber/fiber/v2"
)

func SSEHandler(c *fiber.Ctx) error {

	var userIDPtr *uint

	userIDStr := c.Query("user_id")

	if userIDStr != "" {
		id64, err := strconv.ParseUint(userIDStr, 10, 64)
		if err == nil {
			id := uint(id64)
			userIDPtr = &id
		}
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	client := &sse.Client{
		UserID: userIDPtr, // ✅ pakai pointer yang benar
		Stream: make(chan []byte, 10),
	}

	sse.NotificationHub.AddClient(client)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		defer sse.NotificationHub.RemoveClient(client)

		for {
			select {

			case msg := <-client.Stream:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.Flush()

			case <-c.Context().Done():
				return
			}
		}
	})

	return nil
}