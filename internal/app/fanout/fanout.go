package fanout

import (
	"context"
	"fmt"
	"kafka-forwarder/internal/pkg/fanout_buffer"
	"kafka-forwarder/internal/pkg/ios"
	"kafka-forwarder/internal/pkg/time_buffer"
	"os"
	"strconv"
	"strings"
	"time"
)

var brokerAddress string = os.Getenv("KAFKA_ADDRESS")
var intopic string = os.Getenv("KAFKA_IN_TOPIC")
var consumerGroup string = os.Getenv("KAFKA_CONSUMER_GROUP")
var topicBufferSize string = os.Getenv("TOPIC_BUFFER_SIZE")
var bufferTimeoutInSeconds string = os.Getenv("BUFFER_TIMEOUT_IN_SECONDS")

var outChannel = make(chan []time_buffer.Message)

type QueueMessage struct {
	topic   string
	content string
}

func (q QueueMessage) GetTopic() string {
	return q.topic
}

func (q QueueMessage) GetContent() string {
	return q.content
}

func NewQueueMessage(topic string, content string) QueueMessage {
	return QueueMessage{
		topic:   topic,
		content: content,
	}
}

func Run() {
	fmt.Println("Broker: " + brokerAddress)
	fmt.Println("InTopic: " + intopic)
	fmt.Println("ConsumerGroup: " + consumerGroup)
	fmt.Println("TopicBufferSize: " + topicBufferSize)
	fmt.Println("TimeOutInSeconds: " + bufferTimeoutInSeconds)

	bufferSize, err := strconv.Atoi(topicBufferSize)
	if err != nil {
		panic(err)
	}

	timeoutInSeconds, err := strconv.Atoi(bufferTimeoutInSeconds)
	if err != nil {
		panic(err)
	}

	fbuf := fanout_buffer.New(outChannel, uint(bufferSize), time.Duration(timeoutInSeconds)*time.Second)

	ctx := context.Background()

	input := make(chan string)
	go ios.Consume(ctx, consumerGroup, input)
	fmt.Println("Calling Process ...")
	bufMessagesChannel := FanoutProcess(ctx, fbuf, input)
	ios.FanoutProduce(ctx, bufMessagesChannel)
}

func Process(ctx context.Context, fbuf *fanout_buffer.FBuffer, input chan string) chan string {

	outProcessChannel := make(chan string)

	fmt.Println("Calling Process:Ingesting ...")
	go func() {
		for {
			//take message from channel
			rawMsg := <-input
			lines := strings.Split(rawMsg, "\n")

			for _, line := range lines {
				msgParts := strings.Split(line, "\t")

				topic := msgParts[0]
				content := strings.Join(msgParts[1:], "\t")

				fmt.Printf("Received: {topic: %v, content: %v}\n", topic, content)
				fbuf.Put(topic, content)
			}
		}
	}()

	fmt.Println("Calling Process:Producing ...")
	go func() {
		for {
			bufferedMessages := <-fbuf.OutC
			var messages []string
			for _, msg := range bufferedMessages {
				messages = append(messages, msg.Content)
			}
			topic := bufferedMessages[0].Topic
			aggregatedMessages := fmt.Sprintf("{topic: %v, content: [%v]}", topic, strings.Join(messages, ", "))
			outProcessChannel <- aggregatedMessages
		}
	}()

	return outProcessChannel
}

func FanoutProcess(ctx context.Context, fbuf *fanout_buffer.FBuffer, input chan string) chan ios.KafkaMessage {

	outProcessChannel := make(chan ios.KafkaMessage)

	fmt.Println("Calling Process:Ingesting ...")
	go func() {
		for {
			//take message from channel
			rawMsg := <-input
			lines := strings.Split(rawMsg, "\n")

			for _, line := range lines {
				msgParts := strings.Split(line, "\t")

				topic := msgParts[0]
				content := strings.Join(msgParts[1:], "\t")

				fmt.Printf("Received: {topic: %v, content: %v}\n", topic, content)
				fbuf.Put(topic, content)
			}
		}
	}()

	fmt.Println("Calling Process:Producing ...")
	go func() {
		for {
			bufferedMessages := <-fbuf.OutC
			var messages []string
			for _, msg := range bufferedMessages {
				messages = append(messages, msg.Content)
			}
			topic := bufferedMessages[0].Topic
			aggregatedMessages := fmt.Sprintf("{topic: %v, content: [%v]}", topic, strings.Join(messages, ", "))
			queueMessage := NewQueueMessage(topic, aggregatedMessages)
			outProcessChannel <- queueMessage
		}
	}()

	return outProcessChannel
}
