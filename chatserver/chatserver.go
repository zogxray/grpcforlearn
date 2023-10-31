package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Message struct {
	ClientName       string
	MessageBody      string
	ClientUniqueCode int
}

type MessageQueue struct {
	mu    sync.RWMutex
	queue []Message
}

type ChatServer struct {
	clients map[int]Services_ChatServiceServer
	mu      sync.Mutex
}

func NewChatServer() ChatServer {
	return ChatServer{
		clients: make(map[int]Services_ChatServiceServer),
	}
}

func (mq *MessageQueue) AddMessage(message Message) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.queue = append(mq.queue, message)
}

func (mq *MessageQueue) GetMessages() []Message {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	return mq.queue
}

func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {
	rand.Seed(time.Now().UnixNano())

	clientUniqueCode := rand.Intn(1e6)

	errch := make(chan error, 1)

	mq := &MessageQueue{}

	log.Printf("Adding clent %d", clientUniqueCode)

	is.mu.Lock()
	is.clients[clientUniqueCode] = csi
	is.mu.Unlock()

	log.Printf("Add clent %d :: %v", clientUniqueCode, csi)

	defer func() {
		is.mu.Lock()
		delete(is.clients, clientUniqueCode)
		is.mu.Unlock()
		log.Printf("Removed client %d", clientUniqueCode)
	}()

	go recieveFromStream(csi, clientUniqueCode, mq, errch)
	go sendToStream(is, clientUniqueCode, mq, errch)

	return <-errch
}

func recieveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, mq *MessageQueue, errch_ chan error) {
	for {
		mssg, err := csi_.Recv()

		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch_ <- err
			return
		}

		mq.AddMessage(Message{
			ClientName:       mssg.Name,
			MessageBody:      mssg.Body,
			ClientUniqueCode: clientUniqueCode_,
		})

		log.Printf("Add message %s", mssg.Body)
	}
}

func sendToStream(is *ChatServer, clientUniqueCode_ int, mq *MessageQueue, errch_ chan error) {
	for {
		messages := mq.GetMessages()

		is.mu.Lock()
		mq.mu.RLock()

		for uniqueCode, client := range is.clients {
			if uniqueCode == clientUniqueCode_ {
				continue
			}

			for _, message := range messages {
				err := client.Send(&FromServer{
					Name:    message.ClientName,
					Message: message.MessageBody,
				})

				log.Printf("Send message %s", message.MessageBody)

				if err != nil {
					errch_ <- err

				}
			}
		}
		mq.mu.RUnlock()
		is.mu.Unlock()

		mq.mu.Lock()
		mq.queue = mq.queue[:0]
		mq.mu.Unlock()

		time.Sleep(1 * time.Millisecond)
	}
}
