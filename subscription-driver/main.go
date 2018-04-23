package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	subpb "github.com/adamsanghera/blurber-protobufs/dist/subscription"
	"github.com/adamsanghera/blurber/subscription"

	"google.golang.org/grpc"
)

func createListenAddr() string {
	listenPort := os.Getenv("SUB_PORT")

	// if listenHost == "" {
	// 	panic("No hostname set")
	// }
	if listenPort == "" {
		panic("No port number set")
	}

	return ":" + listenPort
}

func main() {
	addr := createListenAddr()

	log.Printf("SubServer: Derived address: (%s)", addr)

	srv := subscription.NewLedgerServer(addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	subpb.RegisterSubscriptionDBServer(s, srv)

	log.Printf("SubServer: Starting on (%s)", addr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s\n", err.Error())
	}
}
