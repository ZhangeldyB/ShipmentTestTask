package domain

import "time"

type ShipmentEvent struct {
	ID         string    `bson:"_id"`
	ShipmentID string    `bson:"shipment_id"`
	Status     Status    `bson:"status"`
	Note       string    `bson:"note"`
	OccurredAt time.Time `bson:"occurred_at"`
}
