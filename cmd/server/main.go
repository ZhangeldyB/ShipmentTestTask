package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/app"
	grpcinfra "github.com/ZhangeldyB/ShipmentTestTask/internal/infra/grpc"
	mongoinfra "github.com/ZhangeldyB/ShipmentTestTask/internal/infra/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	grpcPort := getenv("GRPC_PORT", "50051")
	databaseURL := mustGetenv("DATABASE_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(databaseURL))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("error disconnecting MongoDB: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("Connected to MongoDB")

	db := client.Database("shipment_service")

	if err := mongoinfra.EnsureIndexes(ctx, db); err != nil {
		log.Fatalf("failed to ensure indexes: %v", err)
	}

	repo := mongoinfra.NewRepository(db)

	createUC := app.NewCreateShipmentUseCase(repo)
	getUC := app.NewGetShipmentUseCase(repo)
	addEventUC := app.NewAddStatusEventUseCase(repo)
	getEventsUC := app.NewGetShipmentEventsUseCase(repo)

	handler := grpcinfra.NewShipmentHandler(createUC, getUC, addEventUC, getEventsUC)
	srv := grpcinfra.NewGRPCServer(handler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down gRPC server...")
	srv.GracefulStop()
	log.Println("Server stopped")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
