package app

import (
	"context"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

type GetShipmentUseCase struct {
	repo ShipmentRepository
}

func NewGetShipmentUseCase(repo ShipmentRepository) *GetShipmentUseCase {
	return &GetShipmentUseCase{repo: repo}
}

func (uc *GetShipmentUseCase) Execute(ctx context.Context, id string) (*domain.Shipment, error) {
	shipment, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}
