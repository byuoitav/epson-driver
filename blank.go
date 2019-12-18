package epson

import (
	"context"
	"fmt"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

func (p *Projector) GetBlanked(ctx context.Context) (bool, error) {
	var blanked bool

	work := func(conn connpool.Conn) error {

		cmd := []byte("MUTE?")
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error getting blanked: %v", err)
		}
		switch checker {
		case "MUTE=ON\r:":
			blanked = true
		case "MUTE=OFF\r:":
			blanked = false
		default:
			return fmt.Errorf("unknown blanked state '%s'", checker)
		}
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return blanked, err
	}

	return blanked, nil
}

func (p *Projector) SetBlanked(ctx context.Context, blanked bool) error {
	work := func(conn connpool.Conn) error {
		var str string

		switch blanked {
		case true:
			str = "ON"
		case false:
			str = "OFF"
		default:
			return fmt.Errorf("unexpected blank state '%v'", blanked)
		}

		cmd := []byte(fmt.Sprintf("MUTE %s", str))
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error setting blank status: %v", err)
		}

		bytes := fmt.Sprintf("%x", checker)

		if bytes != "3a" {
			return fmt.Errorf("There was an error executing the command - %s", bytes)
		}

		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	log.L.Infof("blanking screen")
	return nil
}
