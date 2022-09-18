package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type client struct {
	destination string
	timeout     time.Duration
	in          io.ReadCloser
	out         io.Writer
	connection  net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		destination: address,
		timeout:     timeout,
		in:          in,
		out:         out,
	}
}

func (t *client) Connect() error {
	conn, err := net.Dial("tcp", t.destination)
	if err != nil {
		return err
	}

	t.connection = conn

	return nil
}

func (t *client) Close() error {
	defer func() {
		_ = t.in.Close()
	}()

	return t.connection.Close()
}

func (t *client) Send() error {
	if _, err := io.Copy(t.connection, t.in); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (t *client) Receive() error {
	if _, err := io.Copy(t.out, t.connection); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
