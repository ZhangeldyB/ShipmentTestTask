package domain

import "time"

type ShipmentEvent struct {
	ID         string
	ShipmentID string
	Status     Status
	Note       string
	OccurredAt time.Time
}
