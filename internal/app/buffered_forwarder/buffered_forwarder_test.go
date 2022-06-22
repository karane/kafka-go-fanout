package buffered_forwarder

import (
	"context"
	"kafka-forwarder/internal/pkg/time_buffer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {

	ctx := context.Background()
	input := make(chan string)
	go func() {
		input <- "karane"
		input <- "keila"
		input <- "jose"
		input <- "lucileide"
	}()

	outChannel = make(chan []time_buffer.Message)
	tbuf := time_buffer.New("time-buffer", 1*time.Second, outChannel)

	output := Process(ctx, input, tbuf)
	assert.Equal(t, "[karane, keila]", <-output)
	assert.Equal(t, "[jose, lucileide]", <-output)
	// Cancel()
}

// func TestCancel(t *testing.T) {

// 	ctx := context.Background()
// 	input := make(chan string)
// 	go func() {
// 		input <- "karane"
// 		input <- "keila"
// 		input <- "jose"
// 		input <- "lucileide"
// 	}()

// 	outChannel = make(chan []time_buffer.Message)
// 	tbuf := time_buffer.New("time-buffer", 10*time.Second, outChannel)

// 	output := Process(ctx, input, tbuf)
// 	Cancel()

// 	select {
// 	case <-output:
// 		assert.Fail(t, "Cancel did not work!")
// 	case <-time.After(1 * time.Second):
// 	}

// }
