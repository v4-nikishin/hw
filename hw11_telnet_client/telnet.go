package main

import (
	"bufio"
	"context"
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

type TelnetClientStruct struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientStruct{address: address, timeout: timeout, in: in, out: out}
}

func (tc *TelnetClientStruct) Connect() error {
	dialer := &net.Dialer{}
	tc.ctx, tc.cancel = context.WithTimeout(context.Background(), tc.timeout)

	var err error
	tc.conn, err = dialer.DialContext(tc.ctx, "tcp", tc.address)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TelnetClientStruct) Close() error {
	tc.cancel()
	if tc.conn != nil {
		tc.conn.Close()
	}
	return nil
}

func (tc *TelnetClientStruct) Send() error {
	return tc.scan(bufio.NewScanner(tc.in), tc.conn)
}

func (tc *TelnetClientStruct) Receive() error {
	defer tc.cancel()
	return tc.scan(bufio.NewScanner(tc.conn), tc.out)
}

func (tc *TelnetClientStruct) scan(scanner *bufio.Scanner, writer io.Writer) error {
	for {
		select {
		case <-tc.ctx.Done():
			return fmt.Errorf("finished by context done")
		default:
			if !scanner.Scan() {
				return nil
			}
			text := scanner.Text()
			if _, err := writer.Write([]byte(fmt.Sprintf("%s\n", text))); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
		}
	}
}
