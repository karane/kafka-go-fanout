package ios

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

var intopic string = os.Getenv("KAFKA_IN_TOPIC")   //"intopic"
var outtopic string = os.Getenv("KAFKA_OUT_TOPIC") //"outtopic"
var brokerAddress string = os.Getenv("KAFKA_ADDRESS")

type KafkaMessage interface {
	GetTopic() string
	GetContent() string
}

func Produce(ctx context.Context, input chan string) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   outtopic,
	})

	for {
		//take message from channel
		msg := <-input

		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(strconv.Itoa(i)),
			Value: []byte(msg),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++

		// sleep for a second
		time.Sleep(time.Microsecond)
	}
}

func FanoutProduce(ctx context.Context, input chan KafkaMessage) {
	// initialize a counter
	i := 0

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
	})

	for {
		//take message from channel
		msg := <-input

		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte(msg.GetContent()),
			Topic: msg.GetTopic(),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++

		// sleep for a microsecond
		time.Sleep(time.Microsecond)
	}
}

func Consume(ctx context.Context, consumerGroup string, input chan string) {

	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   intopic,
		GroupID: consumerGroup,
		Logger:  l,
	})

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		fmt.Println("received: ", string(msg.Value))

		// pass to channel
		input <- string(msg.Value)
	}
}
