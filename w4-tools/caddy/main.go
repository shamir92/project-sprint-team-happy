package main

import (
	"log"
	"net"
	"time"

	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity/pb"
	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/handler"
	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	locationService := service.NewMerchantService()

	handlers := handler.NewHandler(locationService)

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     30 * time.Second, // Close the connection if it's idle for this duration
			MaxConnectionAge:      1 * time.Minute,  // Maximum age of a connection
			MaxConnectionAgeGrace: 1 * time.Minute,  // Allow a grace period for connection closure
			Time:                  30 * time.Second, // Ping the client if it is idle for this duration to ensure the connection is still active
			Timeout:               20 * time.Second, // Wait this long for the ping ack before considering the connection dead
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second, // Minimum time a client should wait before sending a keepalive ping
			PermitWithoutStream: true,             // Allow pings when there are no active streams
		}),
	)
	pb.RegisterMerchantServiceServer(server, handlers)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server is running on port :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
