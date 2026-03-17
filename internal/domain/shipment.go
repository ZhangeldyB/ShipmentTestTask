package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	TransportMode   TransportMode
	CarrierInfo     CarrierInfo
	Status          Status
	Amount          float64
	CarrierRevenue  float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
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
