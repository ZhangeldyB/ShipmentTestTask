package mongo

import "time"

type shipmentDoc struct {
	ID              string    `bson:"_id"`
	ReferenceNumber string    `bson:"reference_number"`
	Origin          string    `bson:"origin"`
	Destination     string    `bson:"destination"`
	TransportMode   string    `bson:"transport_mode"`
	OperatorName    string    `bson:"operator_name"`
	OperatorPhone   string    `bson:"operator_phone"`
	UnitIdentifier  string    `bson:"unit_identifier"`
	Status          string    `bson:"status"`
	Amount          float64   `bson:"amount"`
	CarrierRevenue  float64   `bson:"carrier_revenue"`
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
