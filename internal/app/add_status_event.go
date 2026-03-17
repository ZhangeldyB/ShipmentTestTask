package app

import (
	"context"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

type AddStatusEventUseCase struct {
	repo ShipmentRepository
}

func NewAddStatusEventUseCase(repo ShipmentRepository) *AddStatusEventUseCase {
	return &AddStatusEventUseCase{repo: repo}
}

func (uc *AddStatusEventUseCase) Execute(ctx context.Context, shipmentID string, newStatus domain.Status, note string) (*domain.Shipment, error) {
	shipment, err := uc.repo.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	event, err := shipment.ApplyEvent(newStatus, note)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, shipment); err != nil {
		return nil, err
	}
	if err := uc.repo.SaveEvent(ctx, event); err != nil {
		return nil, err
	}

	return shipment, nil
}
