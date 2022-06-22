package forwarder

import (
	"context"
	"fmt"
	"kafka-forwarder/internal/pkg/ios"
	"os"
)

var intopic string = os.Getenv("KAFKA_IN_TOPIC")   //"intopic"
var outtopic string = os.Getenv("KAFKA_OUT_TOPIC") //"outtopic"
var brokerAddress string = os.Getenv("KAFKA_ADDRESS")
var consumerGroup string = os.Getenv("KAFKA_CONSUMER_GROUP")

func Run() {
	fmt.Println("Broker: " + brokerAddress)
	fmt.Println("InTopic: " + intopic)
	fmt.Println("OutTopic: " + outtopic)
	fmt.Println("KafkaConsumerGroup: " + consumerGroup)
	ctx := context.Background()
	input := make(chan string)
	go ios.Consume(ctx, consumerGroup, input)
	ios.Produce(ctx, input)
}
