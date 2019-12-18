package epson

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

func (p *Projector) writeAndRead(ctx context.Context, conn connpool.Conn, cmd []byte, timeout time.Duration, delim byte) (string, error) {
	conn.SetWriteDeadline(time.Now().Add(timeout))

	n, err := conn.Write(cmd)
	switch {
	case err != nil:
		return "", err
	case n != len(cmd):
		return "", fmt.Errorf("wrote %v/%v bytes of command 0x%x", n, len(cmd), cmd)
	}

	b, err := conn.ReadUntil(delim, timeout)
	if err != nil {
		return "", err
	}

	log.L.Debugf("Response from command: 0x%x", b)
	return strings.TrimSpace(string(b)), nil
}
