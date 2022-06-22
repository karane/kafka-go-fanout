package fanout_buffer

import (
	"kafka-forwarder/internal/pkg/time_buffer"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertOneElement(t *testing.T) {

	outChannel := make(chan []time_buffer.Message)
	fanoutBuffer := New(outChannel, 3, 2*time.Second)

	fanoutBuffer.Put("topic", "karane")
	assert.Equal(t, []time_buffer.Message{{Topic: "topic", Content: "karane"}}, <-outChannel)
}

func TestInsertTwoElements(t *testing.T) {

	outChannel := make(chan []time_buffer.Message)
	fanoutBuffer := New(outChannel, 3, 2*time.Second)

	fanoutBuffer.Put("topic", "karane")
	fanoutBuffer.Put("topic", "keila")
	assert.Equal(t, []time_buffer.Message{{Topic: "topic", Content: "karane"}, {Topic: "topic", Content: "keila"}}, <-outChannel)
}

func TestInsertThreeElements(t *testing.T) {

	outChannel := make(chan []time_buffer.Message)
	fanoutBuffer := New(outChannel, 3, 2*time.Second)

	fanoutBuffer.Put("topic", "karane")
	fanoutBuffer.Put("topic", "keila")
	fanoutBuffer.Put("topic", "bob")
	assert.Equal(t, []time_buffer.Message{
		{Topic: "topic", Content: "karane"},
		{Topic: "topic", Content: "keila"},
		{Topic: "topic", Content: "bob"}},
		<-outChannel)
}

func TestInsertThreeElementsTopicSizeTwo(t *testing.T) {

	outChannel := make(chan []time_buffer.Message)
	fanoutBuffer := New(outChannel, 2, 2*time.Second)

	go func() {
		fanoutBuffer.Put("topic", "karane")
		fanoutBuffer.Put("topic", "keila")
		fanoutBuffer.Put("topic", "bob")
	}()

	assert.Equal(t, []time_buffer.Message{
		{Topic: "topic", Content: "karane"},
		{Topic: "topic", Content: "keila"}},
		<-outChannel)
	assert.Equal(t, []time_buffer.Message{
		{Topic: "topic", Content: "bob"}},
		<-outChannel)
}

func TestInsertThreeElementsTwoTopics(t *testing.T) {

	outChannel := make(chan []time_buffer.Message)
	fanoutBuffer := New(outChannel, 2, 2*time.Second)

	go func() {
		fanoutBuffer.Put("A", "karane")
		fanoutBuffer.Put("B", "bob")
		fanoutBuffer.Put("A", "keila")
	}()

	assert.Equal(t, []time_buffer.Message{
		{Topic: "A", Content: "karane"},
		{Topic: "A", Content: "keila"}},
		<-outChannel)
	assert.Equal(t, []time_buffer.Message{
		{Topic: "B", Content: "bob"}},
		<-outChannel)
}
