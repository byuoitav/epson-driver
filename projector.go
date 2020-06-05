package epson

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type Projector struct {
	address string
	pool    *connpool.Pool

	logger Logger

	lastKnownInput   string
	lastKnownBlanked bool
}

func NewProjector(addr string, opts ...Option) *Projector {
	options := options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	p := &Projector{
		address: addr,
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
		logger:           options.logger,
		lastKnownInput:   "hdbaset",
		lastKnownBlanked: false,
	}

	p.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		dial := net.Dialer{}
		conn, err := dial.DialContext(ctx, "tcp", p.address+":3629")
		if err != nil {
			return nil, err
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(5 * time.Second)
		}

		conn.SetDeadline(deadline)

		// send "ESC/VP.net" in order to allow other commands
		cmd := []byte{0x45, 0x53, 0x43, 0x2F, 0x56, 0x50, 0x2E, 0x6E, 0x65, 0x74, 0x10, 0x03, 0x00, 0x00, 0x00, 0x00}

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			conn.Close()
			return conn, fmt.Errorf("unable to write new connection string: %w", err)
		case n != len(cmd):
			conn.Close()
			return conn, fmt.Errorf("unable to write new connection string: wrote %d/%d bytes", n, len(cmd))
		}

		resp := make([]byte, len(cmd))

		// read the same thing back
		n, err = conn.Read(resp)
		switch {
		case err != nil:
			conn.Close()
			return conn, fmt.Errorf("unable to read new connection string: %w", err)
		case n != len(cmd):
			conn.Close()
			return conn, fmt.Errorf("unable to read new connection string: read %d/%d bytes", n, len(cmd))
		}

		return conn, nil
	}

	return p
}

func (p *Projector) sendCommand(ctx context.Context, cmd []byte, delim byte) ([]byte, error) {
	var resp []byte

	err := p.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		conn.SetWriteDeadline(deadline)
		p.debugf("Sending command: %#x", cmd)

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write command: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write command: wrote %d/%d bytes", n, len(cmd))
		}

		resp, err = conn.ReadUntil(delim, deadline)
		if err != nil {
			return fmt.Errorf("unable to read response: %w", err)
		}

		p.debugf("Got response: %#x", resp)
		return nil
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (p *Projector) GetInfo(ctx context.Context) (interface{}, error) {
	return false, errors.New("not implemented")
}

func (p *Projector) GetActiveSignal(ctx context.Context, port string) (bool, error) {
	return false, errors.New("not implemented")
}
