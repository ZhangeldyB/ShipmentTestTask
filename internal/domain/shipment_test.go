package domain

import (
	"testing"
	"time"
)

func newShipment(status Status) *Shipment {
	return &Shipment{
		ID:            "test-id",
		TransportMode: TransportModeTruck,
		CarrierInfo: CarrierInfo{
			OperatorName:   "John Doe",
			UnitIdentifier: "KZ-001",
		},
		Status:    status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func TestApplyEvent_ValidTransitions(t *testing.T) {
	cases := []struct {
		name string
		from Status
		to   Status
	}{
		{"PENDING -> ASSIGNED", StatusPending, StatusAssigned},
		{"ASSIGNED -> PICKED_UP", StatusAssigned, StatusPickedUp},
		{"PICKED_UP -> IN_TRANSIT", StatusPickedUp, StatusInTransit},
		{"IN_TRANSIT -> DELIVERED", StatusInTransit, StatusDelivered},
		{"IN_TRANSIT -> FAILED", StatusInTransit, StatusFailed},
		{"PENDING -> CANCELLED", StatusPending, StatusCancelled},
		{"ASSIGNED -> CANCELLED", StatusAssigned, StatusCancelled},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newShipment(tc.from)
			event, err := s.ApplyEvent(tc.to, "test note")
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if event == nil {
				t.Fatal("expected event, got nil")
			}
			if s.Status != tc.to {
				t.Fatalf("expected status %s, got %s", tc.to, s.Status)
			}
			if event.Status != tc.to {
				t.Fatalf("expected event status %s, got %s", tc.to, event.Status)
			}
		})
	}
}

func TestApplyEvent_InvalidTransitions(t *testing.T) {
	cases := []struct {
		name string
		from Status
		to   Status
	}{
		{"PENDING -> IN_TRANSIT", StatusPending, StatusInTransit},
		{"PENDING -> DELIVERED", StatusPending, StatusDelivered},
		{"PICKED_UP -> DELIVERED", StatusPickedUp, StatusDelivered},
		{"IN_TRANSIT -> ASSIGNED", StatusInTransit, StatusAssigned},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newShipment(tc.from)
			_, err := s.ApplyEvent(tc.to, "")
			if err != ErrInvalidTransition {
				t.Fatalf("expected ErrInvalidTransition, got: %v", err)
			}
		})
	}
}

func TestApplyEvent_TerminalState(t *testing.T) {
	terminals := []Status{StatusDelivered, StatusFailed, StatusCancelled}

	for _, terminal := range terminals {
		t.Run(string(terminal), func(t *testing.T) {
			s := newShipment(terminal)
			_, err := s.ApplyEvent(StatusAssigned, "")
			if err != ErrTerminalState {
				t.Fatalf("expected ErrTerminalState, got: %v", err)
			}
		})
	}
}

func TestApplyEvent_DuplicateStatus(t *testing.T) {
	s := newShipment(StatusPending)
	_, err := s.ApplyEvent(StatusPending, "")
	if err != ErrDuplicateStatus {
		t.Fatalf("expected ErrDuplicateStatus, got: %v", err)
	}
}

func TestTransportMode_Validate(t *testing.T) {
	valid := []TransportMode{TransportModeTruck, TransportModeAir, TransportModeSea, TransportModeRail}
	for _, m := range valid {
		if err := m.Validate(); err != nil {
			t.Fatalf("expected valid mode %s, got error: %v", m, err)
		}
	}

	if err := TransportMode("HORSE").Validate(); err != ErrInvalidTransportMode {
		t.Fatalf("expected ErrInvalidTransportMode, got: %v", err)
	}
}

func TestCarrierInfo_Validate(t *testing.T) {
	cases := []struct {
		name    string
		carrier CarrierInfo
		wantErr error
	}{
		{
			name:    "missing operator name",
			carrier: CarrierInfo{UnitIdentifier: "KZ-001"},
			wantErr: ErrInvalidCarrier,
		},
		{
			name:    "missing unit identifier",
			carrier: CarrierInfo{OperatorName: "John"},
			wantErr: ErrInvalidCarrier,
		},
		{
			name:    "valid carrier",
			carrier: CarrierInfo{OperatorName: "John", UnitIdentifier: "KZ-001"},
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.carrier.Validate(); err != tc.wantErr {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
