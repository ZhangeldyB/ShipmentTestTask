package grpc

import (
	"context"
	"errors"

	pb "github.com/ZhangeldyB/ShipmentTestTask/gen/shipment/v1"
	"github.com/ZhangeldyB/ShipmentTestTask/internal/app"
	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer

	createShipment    *app.CreateShipmentUseCase
	getShipment       *app.GetShipmentUseCase
	addStatusEvent    *app.AddStatusEventUseCase
	getShipmentEvents *app.GetShipmentEventsUseCase
}

func NewShipmentHandler(
	create *app.CreateShipmentUseCase,
	get *app.GetShipmentUseCase,
	addEvent *app.AddStatusEventUseCase,
	getEvents *app.GetShipmentEventsUseCase,
) *ShipmentHandler {
	return &ShipmentHandler{
		createShipment:    create,
		getShipment:       get,
		addStatusEvent:    addEvent,
		getShipmentEvents: getEvents,
	}
}

func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.ShipmentResponse, error) {
	shipment, err := h.createShipment.Execute(ctx, app.CreateShipmentInput{
		Origin:        req.Origin,
		Destination:   req.Destination,
		TransportMode: protoTransportModeToDomain(req.TransportMode),
		CarrierInfo: domain.CarrierInfo{
			OperatorName:   req.OperatorName,
			OperatorPhone:  req.OperatorPhone,
			UnitIdentifier: req.UnitIdentifier,
		},
		Amount:         req.Amount,
		CarrierRevenue: req.CarrierRevenue,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return domainShipmentToProto(shipment), nil
}

func (h *ShipmentHandler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.ShipmentResponse, error) {
	shipment, err := h.getShipment.Execute(ctx, req.ShipmentId)
	if err != nil {
		return nil, mapError(err)
	}
	return domainShipmentToProto(shipment), nil
}

func (h *ShipmentHandler) AddStatusEvent(ctx context.Context, req *pb.AddStatusEventRequest) (*pb.ShipmentResponse, error) {
	shipment, err := h.addStatusEvent.Execute(ctx, req.ShipmentId, domain.Status(req.NewStatus), req.Note)
	if err != nil {
		return nil, mapError(err)
	}
	return domainShipmentToProto(shipment), nil
}

func (h *ShipmentHandler) GetShipmentEvents(ctx context.Context, req *pb.GetShipmentEventsRequest) (*pb.GetShipmentEventsResponse, error) {
	events, err := h.getShipmentEvents.Execute(ctx, req.ShipmentId)
	if err != nil {
		return nil, mapError(err)
	}

	pbEvents := make([]*pb.ShipmentEvent, len(events))
	for i, e := range events {
		pbEvents[i] = domainEventToProto(e)
	}
	return &pb.GetShipmentEventsResponse{Events: pbEvents}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrShipmentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidTransition):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrTerminalState):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrDuplicateStatus):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrInvalidRevenue),
		errors.Is(err, domain.ErrInvalidTransportMode),
		errors.Is(err, domain.ErrInvalidCarrier):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
