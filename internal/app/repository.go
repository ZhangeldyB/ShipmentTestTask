package app

import (
	"context"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

type ShipmentRepository interface {
	Save(ctx context.Context, s *domain.Shipment) error
	FindByID(ctx context.Context, id string) (*domain.Shipment, error)
	FindByReferenceNumber(ctx context.Context, ref string) (*domain.Shipment, error)
	SaveEvent(ctx context.Context, e *domain.ShipmentEvent) error
	FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error)
}
