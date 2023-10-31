package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcChatServer/chatserver"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func main() {
	fmt.Println("Connect localhost:5001...")

	server := "localhost:5000"

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

	ch := clientHandle{stream: stream, clientName: "Check"}

	go func() {
		for {
			message, err := readConsoleInput()

			if message == 0 {
				continue
			}

			if err != nil {
				log.Fatalf("Error reading console input: %v", err)
			}

			for i := 1; i <= message; i++ {
				clientMassageBox := &chatserver.FromClient{
					Name: ch.clientName,
					Body: string(strconv.Itoa(i)),
				}

				err = ch.stream.Send(clientMassageBox)

				if err != nil {
					log.Printf("Error while send message to server :: %v", err)
				}

				log.Printf("Send message to server :: %s", clientMassageBox.Body)
			}

		}
	}()

	select {}
}

func readConsoleInput() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter message to send: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	message = strings.Trim(message, "\r\n")

	num, err := strconv.Atoi(message)
	if err != nil {
		return 0, err
	}

	return num, nil
}
