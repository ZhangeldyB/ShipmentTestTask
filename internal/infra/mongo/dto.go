package mongo

import "time"

type shipmentDoc struct {
	ID              string    `bson:"_id"`
	ReferenceNumber string    `bson:"reference_number"`
	Origin          string    `bson:"origin"`
	Destination     string    `bson:"destination"`
	Status          string    `bson:"status"`
	DriverName      string    `bson:"driver_name"`
	DriverPhone     string    `bson:"driver_phone"`
	UnitNumber      string    `bson:"unit_number"`
	Amount          float64   `bson:"amount"`
	DriverRevenue   float64   `bson:"driver_revenue"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

type shipmentEventDoc struct {
	ID         string    `bson:"_id"`
	ShipmentID string    `bson:"shipment_id"`
	Status     string    `bson:"status"`
	Note       string    `bson:"note"`
	OccurredAt time.Time `bson:"occurred_at"`
}
