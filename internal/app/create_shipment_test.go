package app

import (
	"context"
	"testing"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

// mockRepo is a test double for ShipmentRepository.
type mockRepo struct {
	saveFn           func(ctx context.Context, s *domain.Shipment) error
	findByIDFn       func(ctx context.Context, id string) (*domain.Shipment, error)
	findByRefFn      func(ctx context.Context, ref string) (*domain.Shipment, error)
	saveEventFn      func(ctx context.Context, e *domain.ShipmentEvent) error
	findEventsByIDFn func(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error)
}

func (m *mockRepo) Save(ctx context.Context, s *domain.Shipment) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, s)
	}
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrShipmentNotFound
}

func (m *mockRepo) FindByReferenceNumber(ctx context.Context, ref string) (*domain.Shipment, error) {
	if m.findByRefFn != nil {
		return m.findByRefFn(ctx, ref)
	}
	return nil, domain.ErrShipmentNotFound
}

func (m *mockRepo) SaveEvent(ctx context.Context, e *domain.ShipmentEvent) error {
	if m.saveEventFn != nil {
		return m.saveEventFn(ctx, e)
	}
	return nil
}

func (m *mockRepo) FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	if m.findEventsByIDFn != nil {
		return m.findEventsByIDFn(ctx, shipmentID)
	}
	return nil, nil
}

var validInput = CreateShipmentInput{
	Origin:        "Almaty",
	Destination:   "Astana",
	TransportMode: domain.TransportModeTruck,
	CarrierInfo: domain.CarrierInfo{
		OperatorName:   "Aibek",
		OperatorPhone:  "+7700000000",
		UnitIdentifier: "KZ-001",
	},
	Amount:         1000,
	CarrierRevenue: 700,
}

func TestCreateShipment_Validation(t *testing.T) {
	cases := []struct {
		name    string
		input   CreateShipmentInput
		wantErr error
	}{
		{
			name: "invalid transport mode",
			input: CreateShipmentInput{
				Origin: "A", Destination: "B",
				TransportMode: "HORSE",
				CarrierInfo:   domain.CarrierInfo{OperatorName: "X", UnitIdentifier: "Y"},
				Amount:        100, CarrierRevenue: 50,
			},
			wantErr: domain.ErrInvalidTransportMode,
		},
		{
			name: "missing operator name",
			input: CreateShipmentInput{
				Origin: "A", Destination: "B",
				TransportMode: domain.TransportModeTruck,
				CarrierInfo:   domain.CarrierInfo{UnitIdentifier: "KZ-001"},
				Amount:        100, CarrierRevenue: 50,
			},
			wantErr: domain.ErrInvalidCarrier,
		},
		{
			name: "amount is zero",
			input: CreateShipmentInput{
				Origin: "A", Destination: "B",
				TransportMode: domain.TransportModeTruck,
				CarrierInfo:   domain.CarrierInfo{OperatorName: "X", UnitIdentifier: "Y"},
				Amount:        0, CarrierRevenue: 0,
			},
			wantErr: domain.ErrInvalidAmount,
		},
		{
			name: "carrier revenue exceeds amount",
			input: CreateShipmentInput{
				Origin: "A", Destination: "B",
				TransportMode: domain.TransportModeTruck,
				CarrierInfo:   domain.CarrierInfo{OperatorName: "X", UnitIdentifier: "Y"},
				Amount:        100, CarrierRevenue: 150,
			},
			wantErr: domain.ErrInvalidRevenue,
		},
		{
			name: "carrier revenue is zero",
			input: CreateShipmentInput{
				Origin: "A", Destination: "B",
				TransportMode: domain.TransportModeTruck,
				CarrierInfo:   domain.CarrierInfo{OperatorName: "X", UnitIdentifier: "Y"},
				Amount:        100, CarrierRevenue: 0,
			},
			wantErr: domain.ErrInvalidRevenue,
		},
	}

	repo := &mockRepo{}
	uc := NewCreateShipmentUseCase(repo)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tc.input)
			if err != tc.wantErr {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestCreateShipment_ValidInput(t *testing.T) {
	repo := &mockRepo{}
	uc := NewCreateShipmentUseCase(repo)

	shipment, err := uc.Execute(context.Background(), validInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipment.Status != domain.StatusPending {
		t.Fatalf("expected PENDING, got %s", shipment.Status)
	}
	if shipment.ReferenceNumber == "" {
		t.Fatal("expected non-empty ReferenceNumber")
	}
	if shipment.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if shipment.TransportMode != domain.TransportModeTruck {
		t.Fatalf("expected TRUCK, got %s", shipment.TransportMode)
	}
	if shipment.CarrierInfo.OperatorName != validInput.CarrierInfo.OperatorName {
		t.Fatalf("expected operator %s, got %s", validInput.CarrierInfo.OperatorName, shipment.CarrierInfo.OperatorName)
	}
}
