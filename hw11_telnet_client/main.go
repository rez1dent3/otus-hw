package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "timeout to connect")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatalf("Not enough arguments for the client to work")
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	connectMessage := []byte("...Connected to " + address + "\n")
	if _, err := os.Stderr.Write(connectMessage); err != nil {
		log.Fatal(err)
	}

	errChannel := make(chan error)
	go func() {
		errChannel <- client.Receive()
	}()

	go func() {
		errChannel <- client.Send()
	}()

	sigintChannel, cancelChannel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancelChannel()

	var message []byte
	select {
	case <-errChannel:
		message = []byte("...Connection was closed by peer\n")
		cancelChannel()
	case <-sigintChannel.Done():
		message = []byte("...EOF\n")
	}

	if err := client.Close(); err != nil {
		log.Printf("Cannot close connection: %v", err)
		return
	}

	if _, err := os.Stderr.Write(message); err != nil {
		log.Printf("Cannot write to stderr: %v", err)
		return
	}
}
