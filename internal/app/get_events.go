package app

import (
	"context"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

type GetShipmentEventsUseCase struct {
	repo ShipmentRepository
}

func NewGetShipmentEventsUseCase(repo ShipmentRepository) *GetShipmentEventsUseCase {
	return &GetShipmentEventsUseCase{repo: repo}
}

func (uc *GetShipmentEventsUseCase) Execute(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	if _, err := uc.repo.FindByID(ctx, shipmentID); err != nil {
		return nil, err
	}
	return uc.repo.FindEventsByShipmentID(ctx, shipmentID)
}
