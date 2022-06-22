package main

import (
	"kafka-forwarder/internal/app/fanout"
)

func main() {
	fanout.Run()
}
