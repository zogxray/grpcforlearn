package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcChatServer/chatserver"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Enter server IP:PORT ::: ")

	server := "localhost:5000"

	log.Println("Connecting : " + server)

	conn, err := grpc.Dial(server, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Falied to connect gRPC server :: %v", err)
	}

	defer conn.Close()

	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())

	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	ch := clientHandle{stream: stream}
	ch.clientName = string("Monitor-" + strconv.Itoa(rand.Intn(1e2)))
	go ch.sendMessage()
	go ch.receiveMessage()

	select {}
}

func (ch *clientHandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Falied to read from console :: %v", err)
		}

		message = strings.Trim(message, "\r\n")

		if message == "" {
			continue
		}

		clientMassageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: message,
		}

		err = ch.stream.Send(clientMassageBox)

		if err != nil {
			log.Printf("Error while send message to server :: %v", err)
		}
	}
}

func (ch *clientHandle) receiveMessage() {
	for {
		mssg, err := ch.stream.Recv()

		if err != nil {
			log.Printf("Error while receiving message to server :: %v", err)
		}

		log.Printf("%s : %s \r\n", mssg.Name, mssg.Message)
	}
}
