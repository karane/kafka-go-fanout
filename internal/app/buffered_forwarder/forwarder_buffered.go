package buffered_forwarder

import (
	"context"
	"fmt"
	"kafka-forwarder/internal/pkg/ios"
	"kafka-forwarder/internal/pkg/time_buffer"
	"os"
	"strconv"
	"strings"
	"time"
)

var intopic string = os.Getenv("KAFKA_IN_TOPIC")
var consumerGroup string = os.Getenv("KAFKA_CONSUMER_GROUP")
var brokerAddress string = os.Getenv("KAFKA_ADDRESS")
var bufferTimeoutInSeconds string = os.Getenv("BUFFER_TIMEOUT_IN_SECONDS")

var outChannel chan []time_buffer.Message
var tbuf *time_buffer.Buffer

func Run() {
	fmt.Println("Broker: " + brokerAddress)
	fmt.Println("InTopic: " + intopic)
	fmt.Println("ConsumerGroup: " + consumerGroup)

	timeoutInSeconds, err := strconv.Atoi(bufferTimeoutInSeconds)
	if err != nil {
		panic(err)
	}

	outChannel = make(chan []time_buffer.Message)
	tbuf = time_buffer.New("time-buffer", time.Duration(timeoutInSeconds)*time.Second, outChannel)

	ctx := context.Background()

	input := make(chan string)
	go ios.Consume(ctx, consumerGroup, input)
	fmt.Println("Calling Process ...")
	bufMessagesChannel := Process(ctx, input, tbuf)
	ios.Produce(ctx, bufMessagesChannel)
}

func Process(ctx context.Context, input chan string, tbuf *time_buffer.Buffer) chan string {

	tbuf.StartTimer()
	tbuf.SetMaxLength(2)
	tbuf.Run()

	outProcessChannel := make(chan string)

	fmt.Println("Calling Process:Ingesting ...")
	go func() {
		for {
			//take message from channel
			msg := time_buffer.Message{Topic: "A", Content: <-input}
			fmt.Printf("Received: %v\n", msg)
			tbuf.InsertElement(msg)
		}
	}()

	fmt.Println("Calling Process:Producing ...")
	go func() {
		for {
			bufferedMessages := <-tbuf.OutC
			var messages []string
			for _, msg := range bufferedMessages {
				messages = append(messages, msg.Content)
			}
			aggregatedMessages := fmt.Sprintf("[%v]", strings.Join(messages, ", "))
			outProcessChannel <- aggregatedMessages
		}
	}()

	return outProcessChannel
}

// func Cancel() {
// 	tbuf.Cancel()
// }
