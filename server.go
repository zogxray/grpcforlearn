package main

import (
	"grpcChatServer/chatserver"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	port := "5000"

	listen, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", port, err)
	}

	log.Println("Listening @ :" + port)

	grpcserver := grpc.NewServer()

	cs := chatserver.NewChatServer()
	chatserver.RegisterServicesServer(grpcserver, &cs)

	err = grpcserver.Serve(listen)

	if err != nil {
		log.Fatalf("Falied to start gRPC Server %v", err)
	}

}
