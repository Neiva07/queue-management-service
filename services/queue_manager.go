package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/segmentio/kafka-go"
)

// var Producer *kafka.Producer
// var Consumer *kafka.Consumer
var writer *kafka.Writer
var reader *kafka.Reader
var conn *kafka.Conn
var queueCount int

var userTopic string

// func init() {
// 	Producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer Producer.Close()

// 	// Delivery report handler for produced messages
// 	go func() {
// 		for e := range Producer.Events() {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
// 				} else {
// 					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
// 				}
// 			}
// 		}
// 	}()

// 	userTopic = "UserTpic"
// }

// func init() {
// 	Consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost",
// 		"group.id":          "myGroup",
// 		"auto.offset.reset": "earliest",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = Consumer.SubscribeTopics([]string{userTopic, "^aRegex.*[Tt]opic"}, nil)

// 	if err != nil {
// 		Consumer.Close()

// 	}

// }

// func EnqueueUser(user *persistence.User) error {
// 	parsedMessage, err := parseMessage(user)

// 	if err != nil {
// 		return err
// 	}

// 	Producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &userTopic, Partition: 0},
// 		Value:          parsedMessage,
// 	}, nil)

// 	return nil
// }

// func PopUsers(quantity uint64) []*persistence.User {

// 	users := make([]*persistence.User, quantity)

// 	for i := 0; i < int(quantity); i++ {
// 		var u *persistence.User

// 		msg, err := Consumer.ReadMessage(-1)
// 		if err != nil {
// 			// The client will automatically try to recover from all errors.
// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 		}

// 		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

// 		json.Unmarshal(msg.Value, u)

// 		users[i] = u
// 	}
// 	return users
// }

func init() {
	partition := 0
	userTopic = "UserTopic"
	var err error

	conn, err = kafka.DialLeader(context.Background(), "tcp", "localhost:9092", userTopic, partition)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		log.Fatal("failed to set deadline:", err)
	}

	writer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    userTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func init() {

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "consumer-group-id",
		Topic:          userTopic,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})

}

func EnqueueUser(user *persistence.User) error {
	parsedMessage, err := parseMessage(user)

	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(user.UserID.String()),
			Value: parsedMessage,
		})
	queueCount++

	return err
}

func PopUsers(quantity uint64) []*persistence.User {

	var users []*persistence.User

	var usersToPop uint64

	if quantity > uint64(queueCount) {
		usersToPop = uint64(queueCount)
		queueCount = 0
	} else {
		usersToPop = quantity
		queueCount -= int(quantity)
	}

	for i := 0; i < int(usersToPop); i++ {
		var u persistence.User
		var err error
		var msg kafka.Message

		if msg, err = reader.ReadMessage(context.Background()); err != nil {
			log.Println(err)
			break
		}
		// msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		fmt.Printf("Message on %s: %s\n", string(msg.Key), string(msg.Value))

		json.Unmarshal(msg.Value, &u)

		users = append(users, &u)
	}

	return users
}

func parseMessage(message interface{}) ([]byte, error) {

	parsedMessage, err := json.Marshal(message)

	if err != nil {
		log.Println("Error parsing message for queue")
		return nil, err
	}

	return parsedMessage, nil
}
