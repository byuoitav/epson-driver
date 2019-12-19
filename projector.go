package epson

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

type Projector struct {
	Address string
	Pool    *connpool.Pool
}

var (
	_defaultTTL   = 45 * time.Second
	_defaultDelay = 400 * time.Millisecond
)

type options struct {
	ttl   time.Duration
	delay time.Duration
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithTTL(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.ttl = t
	})
}

func WithDelay(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.delay = t
	})
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
		Address: addr,
		Pool: &connpool.Pool{
			TTL:   options.ttl,
			Delay: options.delay,
		},
	}

	p.Pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		dial := net.Dialer{}
		conn, err := dial.DialContext(ctx, "tcp", p.Address+":3629")
		if err != nil {
			return nil, err
		}

		// read the NOKEY line
		pconn := connpool.Wrap(conn)

		//sending "ESC/VP.net" in order to allow other commands
		cmd := []byte{0x45, 0x53, 0x43, 0x2F, 0x56, 0x50, 0x2E, 0x6E, 0x65, 0x74, 0x10, 0x03, 0x00, 0x00, 0x00, 0x00}

		s, err := p.writeAndRead(ctx, pconn, cmd, 5*time.Second, ' ')
		if err != nil {
			log.L.Warnf("there was an error sending the ESC/VP.net command: %v", err)
			return nil, err
		}
		log.L.Infof("connection string: %s\n", s)
		return conn, nil
	}

	return p
}

func (p *Projector) GetInfo(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func (p *Projector) GetActiveSignal(ctx context.Context, port string) (bool, error) {
	return false, nil
}
