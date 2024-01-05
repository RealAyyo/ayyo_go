package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting: %v", err)
		os.Exit(1)
	}

	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	go func() {
		err := client.Receive()
		if err == nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			stop()
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("...EOF")
			}
			stop()
		}
	}()

	<-ctx.Done()
}
