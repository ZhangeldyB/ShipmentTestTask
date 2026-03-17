package grpc

import (
	"time"

	pb "github.com/ZhangeldyB/ShipmentTestTask/gen/shipment/v1"
	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
)

func domainShipmentToProto(s *domain.Shipment) *pb.ShipmentResponse {
	return &pb.ShipmentResponse{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          string(s.Status),
		TransportMode:   domainTransportModeToProto(s.TransportMode),
		OperatorName:    s.CarrierInfo.OperatorName,
		OperatorPhone:   s.CarrierInfo.OperatorPhone,
		UnitIdentifier:  s.CarrierInfo.UnitIdentifier,
		Amount:          s.Amount,
		CarrierRevenue:  s.CarrierRevenue,
		CreatedAt:       s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       s.UpdatedAt.Format(time.RFC3339),
	}
}

func domainEventToProto(e *domain.ShipmentEvent) *pb.ShipmentEvent {
	return &pb.ShipmentEvent{
		Id:         e.ID,
		ShipmentId: e.ShipmentID,
		Status:     string(e.Status),
		Note:       e.Note,
		OccurredAt: e.OccurredAt.Format(time.RFC3339),
	}
}

func protoTransportModeToDomain(m pb.TransportMode) domain.TransportMode {
	switch m {
	case pb.TransportMode_TRANSPORT_MODE_TRUCK:
		return domain.TransportModeTruck
	case pb.TransportMode_TRANSPORT_MODE_AIR:
		return domain.TransportModeAir
	case pb.TransportMode_TRANSPORT_MODE_SEA:
		return domain.TransportModeSea
	case pb.TransportMode_TRANSPORT_MODE_RAIL:
		return domain.TransportModeRail
	default:
		return ""
	}
}

func domainTransportModeToProto(m domain.TransportMode) pb.TransportMode {
	switch m {
	case domain.TransportModeTruck:
		return pb.TransportMode_TRANSPORT_MODE_TRUCK
	case domain.TransportModeAir:
		return pb.TransportMode_TRANSPORT_MODE_AIR
	case domain.TransportModeSea:
		return pb.TransportMode_TRANSPORT_MODE_SEA
	case domain.TransportModeRail:
		return pb.TransportMode_TRANSPORT_MODE_RAIL
	default:
		return pb.TransportMode_TRANSPORT_MODE_UNSPECIFIED
	}
}
