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
		DriverName:      s.DriverName,
		DriverPhone:     s.DriverPhone,
		UnitNumber:      s.UnitNumber,
		Amount:          s.Amount,
		DriverRevenue:   s.DriverRevenue,
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
