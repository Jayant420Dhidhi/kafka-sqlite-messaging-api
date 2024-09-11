# Kafka Messaging API with Go

This project provides a simple API for producing and consuming messages with Kafka using Go. It includes routes to produce messages, get messages, update messages, and delete messages. The application uses an in-memory store and simulates Kafka behavior.

## Prerequisites

- Go 1.18 or later
- Kafka broker (Apache Kafka)
- SQLite3 (for database persistence, optional)

## Setup

### Kafka Setup

1. **Install Kafka:**

   Download Kafka from the [official website](https://kafka.apache.org/downloads). Follow the installation instructions specific to your operating system.

2. **Start Kafka Broker:**

   Ensure that Kafka is running on the default port (`9092`). Start Kafka using:
   ```sh
   kafka-server-start.sh /path/to/kafka/config/server.properties

## Go Application Setup

### Clone the Repository: https://github.com/Jayant420Dhidhi/kafka-sqlite-messaging-api

## Install Dependencies:

### Install Go dependencies using go mod:

```bash
go mod tidy
```

### Ensure that you have the confluent-kafka-go and gorilla/mux packages installed:
```bash
go get github.com/confluentinc/confluent-kafka-go/kafka
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
```

## Running the Application

### Build the Go application using:

```bash
go build -o kafka-go-api
```


### Start the Go application:
```bash
./kafka-go-api
```
The server will start on port 8080.

### Access the API:

#### Produce a Message (POST):

```bash
curl -X POST http://localhost:8080/produce -d '{"topic": "test", "message": "Hello Kafka"}' -H "Content-Type: application/json"
```
#### Get Messages (GET):

```bash
curl http://localhost:8080/messages/test
```

#### Update a Message (PUT):
```bash
curl -X PUT http://localhost:8080/messages/test -d '{"topic": "test", "message": "Updated message"}' -H "Content-Type: application/json"
```

#### Delete Messages (DELETE):
```bash
curl -X DELETE http://localhost:8080/messages/test
```

## Screenshots

### POST


<img width="1440" alt="POST" src="https://github.com/user-attachments/assets/05809890-86ed-4474-ab61-a72b22653b0d">

### GET


<img width="1440" alt="GET" src="https://github.com/user-attachments/assets/8587fc45-4cc4-417d-bbbc-4e250acc0d76">

### PUT


<img width="1440" alt="PUT" src="https://github.com/user-attachments/assets/2717cef3-45e0-43df-8ca3-0679e827bede">

### Database


<img width="1440" alt="Database" src="https://github.com/user-attachments/assets/701693e5-16fc-4bc0-a409-c90cc73febf7">
