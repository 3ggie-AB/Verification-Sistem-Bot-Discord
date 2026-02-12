package service

import (
	"encoding/json"
	"crypto-member/models"
	"crypto-member/sse"
)

func SendNotification(n models.Notification) {

	go func() {

		data, _ := json.Marshal(n)

		sse.NotificationHub.Broadcast(n.UserID, data)

	}()
}