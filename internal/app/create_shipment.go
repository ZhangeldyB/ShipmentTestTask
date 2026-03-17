package app

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
	"github.com/google/uuid"
)

type CreateShipmentInput struct {
	Origin         string
	Destination    string
	TransportMode  domain.TransportMode
	CarrierInfo    domain.CarrierInfo
	Amount         float64
	CarrierRevenue float64
}

type CreateShipmentUseCase struct {
	repo ShipmentRepository
}

func NewCreateShipmentUseCase(repo ShipmentRepository) *CreateShipmentUseCase {
	return &CreateShipmentUseCase{repo: repo}
}

func (uc *CreateShipmentUseCase) Execute(ctx context.Context, in CreateShipmentInput) (*domain.Shipment, error) {
	if err := in.TransportMode.Validate(); err != nil {
		return nil, err
	}
	if err := in.CarrierInfo.Validate(); err != nil {
		return nil, err
	}
	if in.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}
	if in.CarrierRevenue <= 0 || in.CarrierRevenue > in.Amount {
		return nil, domain.ErrInvalidRevenue
	}

	now := time.Now().UTC()

	shipment := &domain.Shipment{
		ID:              uuid.NewString(),
		ReferenceNumber: generateReferenceNumber(now),
		Origin:          in.Origin,
		Destination:     in.Destination,
		TransportMode:   in.TransportMode,
		CarrierInfo:     in.CarrierInfo,
		Status:          domain.StatusPending,
		Amount:          in.Amount,
		CarrierRevenue:  in.CarrierRevenue,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	event := &domain.ShipmentEvent{
		ID:         uuid.NewString(),
		ShipmentID: shipment.ID,
		Status:     domain.StatusPending,
		Note:       "Shipment created",
		OccurredAt: now,
	}

	if err := uc.repo.Save(ctx, shipment); err != nil {
		return nil, err
	}
	if err := uc.repo.SaveEvent(ctx, event); err != nil {
		return nil, err
	}

	return shipment, nil
}

func generateReferenceNumber(t time.Time) string {
	const hexChars = "0123456789ABCDEF"
	suffix := make([]byte, 4)
	for i := range suffix {
		suffix[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return fmt.Sprintf("SHP-%s-%s", t.Format("20060102"), strings.ToUpper(string(suffix)))
}
