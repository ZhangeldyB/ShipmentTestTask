package mongo

import "github.com/ZhangeldyB/ShipmentTestTask/internal/domain"

func toShipmentDoc(s *domain.Shipment) *shipmentDoc {
	return &shipmentDoc{
		ID:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		TransportMode:   string(s.TransportMode),
		OperatorName:    s.CarrierInfo.OperatorName,
		OperatorPhone:   s.CarrierInfo.OperatorPhone,
		UnitIdentifier:  s.CarrierInfo.UnitIdentifier,
		Status:          string(s.Status),
		Amount:          s.Amount,
		CarrierRevenue:  s.CarrierRevenue,
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
		TransportMode:   domain.TransportMode(doc.TransportMode),
		CarrierInfo: domain.CarrierInfo{
			OperatorName:   doc.OperatorName,
			OperatorPhone:  doc.OperatorPhone,
			UnitIdentifier: doc.UnitIdentifier,
		},
		Status:         domain.Status(doc.Status),
		Amount:         doc.Amount,
		CarrierRevenue: doc.CarrierRevenue,
		CreatedAt:      doc.CreatedAt,
		UpdatedAt:      doc.UpdatedAt,
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
