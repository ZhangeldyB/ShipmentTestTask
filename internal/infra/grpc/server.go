package grpc

import (
	"time"

	pb "github.com/ZhangeldyB/ShipmentTestTask/gen/shipment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(handler *ShipmentHandler) *grpc.Server {
	srv := grpc.NewServer(
		// Allow up to 1 000 concurrent streams per connection.
		grpc.MaxConcurrentStreams(1000),

		// Keep idle connections alive so clients don't pay a new TCP
		// handshake on every burst of traffic.
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     30 * time.Second,
			MaxConnectionAge:      2 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  10 * time.Second,
			Timeout:               3 * time.Second,
		}),

		// Enforce that clients respect the keepalive policy so misbehaving
		// clients cannot hold connections open indefinitely.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	pb.RegisterShipmentServiceServer(srv, handler)
	reflection.Register(srv)
	return srv
}
