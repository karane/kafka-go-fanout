package fanout

import (
	"context"
	"kafka-forwarder/internal/pkg/fanout_buffer"
	"kafka-forwarder/internal/pkg/time_buffer"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessOneTopic(t *testing.T) {

	ctx := context.Background()
	input := make(chan string)
	go func() {
		input <- "A\tkarane\n" +
			"A\tkeila\n" +
			"A\tjose\n" +
			"A\tlucileide\n"
	}()

	outChannel := make(chan []time_buffer.Message)
	topicSize := 2

	fbuf := fanout_buffer.New(outChannel, uint(topicSize))
	output := Process(ctx, fbuf, input)
	assert.Equal(t, "{topic: A, content: [karane, keila]}", <-output)
	assert.Equal(t, "{topic: A, content: [jose, lucileide]}", <-output)
}

func TestProcessTwoTopics(t *testing.T) {

	ctx := context.Background()
	input := make(chan string)
	go func() {
		input <- "A\tkarane\n" +
			"A\tkeila\n" +
			"B\tjose\n" +
			"B\tlucileide"
	}()

	outChannel := make(chan []time_buffer.Message)
	topicSize := 2

	fbuf := fanout_buffer.New(outChannel, uint(topicSize))
	output := Process(ctx, fbuf, input)

	var resp []string
	resp = append(resp, <-output)
	resp = append(resp, <-output)
	sort.Strings(resp)

	assert.Equal(t,
		[]string{
			"{topic: A, content: [karane, keila]}",
			"{topic: B, content: [jose, lucileide]}",
		}, resp)
}

func TestFanoutProcessOneTopic(t *testing.T) {

	ctx := context.Background()
	input := make(chan string)
	go func() {
		input <- "A\tkarane\n" +
			"A\tkeila\n" +
			"A\tjose\n" +
			"A\tlucileide\n"
	}()

	outChannel := make(chan []time_buffer.Message)
	topicSize := 2

	fbuf := fanout_buffer.New(outChannel, uint(topicSize))
	output := FanoutProcess(ctx, fbuf, input)
	queuedMessage := <-output
	assert.Equal(t, "A", queuedMessage.GetTopic())
	assert.Equal(t, "{topic: A, content: [karane, keila]}", queuedMessage.GetContent())
	queuedMessage = <-output
	assert.Equal(t, "A", queuedMessage.GetTopic())
	assert.Equal(t, "{topic: A, content: [jose, lucileide]}", queuedMessage.GetContent())
}

func TestFanoutProcessTwoTopics(t *testing.T) {

	ctx := context.Background()
	input := make(chan string)
	go func() {
		input <- "A\tkarane\n" +
			"A\tkeila\n" +
			"B\tjose\n" +
			"B\tlucileide"
	}()

	outChannel := make(chan []time_buffer.Message)
	topicSize := 2

	fbuf := fanout_buffer.New(outChannel, uint(topicSize))
	output := FanoutProcess(ctx, fbuf, input)

	var resp []string
	resp = append(resp, (<-output).GetContent())
	resp = append(resp, (<-output).GetContent())
	sort.Strings(resp)

	assert.Equal(t,
		[]string{
			"{topic: A, content: [karane, keila]}",
			"{topic: B, content: [jose, lucileide]}",
		}, resp)
}
