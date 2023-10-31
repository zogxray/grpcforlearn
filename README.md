## gRPC Chat Server (for education purposes)

### `chatserver/chatserver.go`

In this code, we are building a basic chat server in Go. The server is designed to handle multiple clients who can connect, send messages, and receive messages from other clients. Let's break down the key components and their functionalities:

### Message Struct

The `Message` struct is used to represent a chat message. It contains the following fields:
- `ClientName`: The name of the client sending the message.
- `MessageBody`: The content of the message.
- `ClientUniqueCode`: A unique code to identify the client who sent the message.

### MessageQueue Struct

The `MessageQueue` struct manages a queue of messages. It employs a mutex for thread safety and includes the following methods:
- `AddMessage(message Message)`: Adds a message to the queue while locking and unlocking the mutex for safe concurrent access.
- `GetMessages() []Message`: Retrieves messages from the queue, locking and unlocking the mutex to ensure thread safety.

### ChatServer Struct

The `ChatServer` struct serves as the core of the chat server and is responsible for managing connected clients. It consists of the following fields:
- `clients`: A map that stores connected clients using their unique client codes.
- `mu`: A mutex to protect access to the `clients` map.

### NewChatServer Function

The `NewChatServer` function initializes a new `ChatServer` instance. It sets up the `clients` map during initialization to store information about connected clients.

### AddMessage and GetMessages Methods

These methods are part of the `MessageQueue` struct and are used to add and retrieve messages from the message queue while ensuring safe concurrent access. They lock and unlock the mutex for protection.

### ChatService Method

The `ChatService` method is called for each new client connection. It performs the following tasks:
- Generates a unique client code for the connecting client.
- Adds the client to the `clients` map, using the client code as the key.
- Sets up a message queue for the client to manage their messages.
- Launches two goroutines:
  - One for receiving messages from the client.
  - Another for sending messages to the client.
- The method returns an error channel, allowing it to handle any potential errors during its execution.

### recieveFromStream Goroutine

This goroutine operates in the background and is responsible for continuously receiving messages from the client. If an error occurs during message reception, it logs the error and sends it to the error channel for handling.

### sendToStream Goroutine

The `sendToStream` goroutine is responsible for broadcasting messages from the message queue to all connected clients (excluding the sender). It manages concurrency using mutexes, ensuring that multiple clients can be handled simultaneously. The goroutine sleeps for 1 millisecond between iterations to prevent excessive CPU usage.

### `server.go`

## Simple gRPC Server in Go

Go code snippet that serves as a basic gRPC server. The server is designed to listen for incoming client connections, offering a platform for communication with connected clients.


### Creating a Network Listener

Using the "net" package, the code attempts to create a network listener for incoming connections. It listens on the specified port, and if any error occurs during this process, it logs an error message indicating the issue.

### Starting the Server

Assuming the network listener has been successfully created, the server logs a message indicating that it is now listening on the specified port.

### Setting Up the gRPC Server

The code initializes a gRPC server using the "grpc.NewServer()" function. This server will handle gRPC communication with connected clients.

### Creating a ChatServer Instance

The `cs` variable is assigned a new instance of a `ChatServer` using the "chatserver.NewChatServer()" function. This chat server likely represents the core functionality of the gRPC server, managing client connections and message exchange.

### Registering the Services Server

The `chatserver.RegisterServicesServer` function is called to register the chat services on the gRPC server. This means that the gRPC server is now aware of the chat-related functionality provided by the `ChatServer` instance.

### Starting the gRPC Server

The gRPC server is started using the "grpcserver.Serve(listen)" method. This method binds the server to the network listener, allowing it to accept incoming gRPC connections.

### Handling Server Start Errors

If any error occurs during the server startup, the code logs an error message indicating the issue. This is crucial for diagnosing problems during server initialization.

### `client.go`

## gRPC Client in Go

Go code snippet that serves as a gRPC client. The client connects to a gRPC server and allows me to send and receive chat messages. Let's break down the code step by step:


### Creating a gRPC Connection

The code establishes a connection to the gRPC server using the "grpc.Dial" function. It connects to the server specified by the `server` variable and uses the "grpc.WithInsecure()" option, indicating that the client doesn't require transport security (e.g., SSL/TLS) for this example.

### Setting Up the gRPC Client

The client is initialized using the "chatserver.NewServicesClient" function, which likely represents the chat-related functionality. The client is now capable of making gRPC calls to the server.

### Creating a ClientHandle

A `clientHandle` structure is defined to manage the client's connection and chat-related functionalities. It includes a gRPC stream and a client name.

### Starting Message Sending and Receiving

Two goroutines are launched:
1. `sendMessage()`: This function allows me to send chat messages to the server. It reads user input from the console, sends the message to the server, and handles potential errors.
2. `receiveMessage()`: This function is responsible for receiving chat messages from the server and displaying them to the console.

### sendMessage() Function

In the `sendMessage()` function:
- User input is read from the console.
- The input is trimmed to remove line breaks and whitespaces.
- If the message is empty, it's skipped.
- The client constructs a message using its name and the user's input.
- The message is sent to the server using the gRPC stream. Any errors are handled and logged.

### receiveMessage() Function

In the `receiveMessage()` function:
- The client listens for incoming chat messages from the server.
- Received messages are logged to the console.
- Any errors in message reception are handled and logged.

### Ongoing Operation

The client enters an infinite loop with `select {}` to keep the application running and continuously send and receive messages from the gRPC server.


### `check.go`

## Understanding a gRPC Client for Sending Repeated Messages

Go code snippet that serves as a gRPC client, designed to connect to a gRPC server and send a specific number of repeated messages.

### Sending Messages

Within the goroutine, the code:
- Reads the desired number of messages from the console.
- For each message number in the specified range:
  - Constructs a message using the client's name and the current message number.
  - Sends the message to the server using the gRPC stream.
  - Logs the sent message for reference.



