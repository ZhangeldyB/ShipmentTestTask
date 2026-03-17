package mongo

import "github.com/ZhangeldyB/ShipmentTestTask/internal/domain"

func toShipmentDoc(s *domain.Shipment) *shipmentDoc {
	return &shipmentDoc{
		ID:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          string(s.Status),
		DriverName:      s.DriverName,
		DriverPhone:     s.DriverPhone,
		UnitNumber:      s.UnitNumber,
		Amount:          s.Amount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func toDomainShipment(doc *shipmentDoc) *domain.Shipment {
	return &domain.Shipment{
		ID:              doc.ID,
		ReferenceNumber: doc.ReferenceNumber,
		Origin:          doc.Origin,
		Destination:     doc.Destination,
		Status:          domain.Status(doc.Status),
		DriverName:      doc.DriverName,
		DriverPhone:     doc.DriverPhone,
		UnitNumber:      doc.UnitNumber,
		Amount:          doc.Amount,
		DriverRevenue:   doc.DriverRevenue,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}
}

func toEventDoc(e *domain.ShipmentEvent) *shipmentEventDoc {
	return &shipmentEventDoc{
		ID:         e.ID,
		ShipmentID: e.ShipmentID,
		Status:     string(e.Status),
		Note:       e.Note,
		OccurredAt: e.OccurredAt,
	}
}

func toDomainEvent(doc *shipmentEventDoc) *domain.ShipmentEvent {
	return &domain.ShipmentEvent{
		ID:         doc.ID,
		ShipmentID: doc.ShipmentID,
		Status:     domain.Status(doc.Status),
		Note:       doc.Note,
		OccurredAt: doc.OccurredAt,
	}
}
