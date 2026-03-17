package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID              string    `bson:"_id"`
	ReferenceNumber string    `bson:"reference_number"`
	Origin          string    `bson:"origin"`
	Destination     string    `bson:"destination"`
	Status          Status    `bson:"status"`
	DriverName      string    `bson:"driver_name"`
	DriverPhone     string    `bson:"driver_phone"`
	UnitNumber      string    `bson:"unit_number"`
	Amount          float64   `bson:"amount"`
	DriverRevenue   float64   `bson:"driver_revenue"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

// ApplyEvent validates the status transition and returns a new ShipmentEvent.
// It also updates the shipment's Status and UpdatedAt on success.
func (s *Shipment) ApplyEvent(newStatus Status, note string) (*ShipmentEvent, error) {
	if IsTerminal(s.Status) {
		return nil, ErrTerminalState
	}
	if s.Status == newStatus {
		return nil, ErrDuplicateStatus
	}
	if !CanTransition(s.Status, newStatus) {
		return nil, ErrInvalidTransition
	}

	now := time.Now().UTC()
	event := &ShipmentEvent{
		ID:         uuid.NewString(),
		ShipmentID: s.ID,
		Status:     newStatus,
		Note:       note,
		OccurredAt: now,
	}

	s.Status = newStatus
	s.UpdatedAt = now

	return event, nil
}
