package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		text := fmt.Sprintf("%s\n", scanner.Text())

		_, err := c.conn.Write([]byte(text))
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		text := fmt.Sprintf("%s\n", scanner.Text())
		_, err := c.out.Write([]byte(text))
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
