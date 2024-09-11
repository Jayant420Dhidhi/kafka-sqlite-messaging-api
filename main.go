package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/gorilla/mux"
    _"github.com/mattn/go-sqlite3"
)

// Message is the struct for the incoming JSON message
type Message struct {
    Topic   string `json:"topic"`
    Message string `json:"message"`
}

// MessageWithTopic is a struct that includes both message and topic
type MessageWithTopic struct {
    Topic   string `json:"topic"`
    Message string `json:"message"`
}

// In-memory map to store messages for different topics
var db *sql.DB
var mu sync.Mutex

// Initialize the database
func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./messages.db")
    if err != nil {
        log.Fatal(err)
    }

    // Create table if it does not exist
    tableCreationQuery := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        topic TEXT,
        message TEXT
    );
    `
    _, err = db.Exec(tableCreationQuery)
    if err != nil {
        log.Fatal(err)
    }
}

// Close the database connection
func closeDB() {
    if err := db.Close(); err != nil {
        log.Fatal(err)
    }
}

// KafkaProducer is a function to produce messages to Kafka
func KafkaProducer(topic string, message string) error {
    // Store message in the SQLite database
    mu.Lock()
    _, err := db.Exec("INSERT INTO messages (topic, message) VALUES (?, ?)", topic, message)
    mu.Unlock()
    if err != nil {
        return err
    }

    // Simulate Kafka producer behavior
    p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
    if err != nil {
        return err
    }
    defer p.Close()

    // Delivery report handler
    go func() {
        for e := range p.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
                } else {
                    fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
                }
            }
        }
    }()

    // Produce message
    p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(message),
    }, nil)

    // Wait for message deliveries
    p.Flush(15 * 1000)
    return nil
}

// KafkaConsumer is a function to consume messages from Kafka (now using in-memory store)
func KafkaConsumer(topic string) ([]MessageWithTopic, error) {
    mu.Lock()
    rows, err := db.Query("SELECT message FROM messages WHERE topic = ?", topic)
    mu.Unlock()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []MessageWithTopic
    for rows.Next() {
        var msg string
        if err := rows.Scan(&msg); err != nil {
            return nil, err
        }
        result = append(result, MessageWithTopic{
            Topic:   topic,
            Message: msg,
        })
    }
    return result, nil
}

// API handler to produce messages (POST)
func produceMessageHandler(w http.ResponseWriter, r *http.Request) {
    var msg Message
    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    err = KafkaProducer(msg.Topic, msg.Message)
    if err != nil {
        http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Message sent to Kafka successfully")
}

// API handler to get messages (GET)
func getMessageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    topic := vars["topic"]
    messages, err := KafkaConsumer(topic)
    if err != nil {
        http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(messages)
}

// API handler to update messages (PUT)
func updateMessageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    topic := vars["topic"]
    var msg Message
    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Simulate an update (Kafka doesn't support updates directly)
    fmt.Printf("Updating message in topic %s: %s\n", topic, msg.Message)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Message updated successfully (simulated)")
}

// API handler to delete messages (DELETE)
func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    topic := vars["topic"]

    // Simulate deletion (Kafka doesn't support deletion directly)
    mu.Lock()
    _, err := db.Exec("DELETE FROM messages WHERE topic = ?", topic)
    mu.Unlock()
    if err != nil {
        http.Error(w, "Failed to delete messages", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Deleting messages from topic %s (simulated)\n", topic)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Messages deleted successfully (simulated)")
}

func main() {
    initDB()
    defer closeDB()

    r := mux.NewRouter()

    // Defined routes for each HTTP method
    r.HandleFunc("/produce", produceMessageHandler).Methods("POST")
    r.HandleFunc("/messages/{topic}", getMessageHandler).Methods("GET")
    r.HandleFunc("/messages/{topic}", updateMessageHandler).Methods("PUT")
    r.HandleFunc("/messages/{topic}", deleteMessageHandler).Methods("DELETE")

    fmt.Println("Starting server on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
