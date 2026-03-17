package grpc

import (
	pb "github.com/ZhangeldyB/ShipmentTestTask/gen/shipment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(handler *ShipmentHandler) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterShipmentServiceServer(srv, handler)
	reflection.Register(srv)
	return srv
}
