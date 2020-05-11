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
		log.L.Infof("Getting blank status")

		cmd := []byte("MUTE?")
		cmd = append(cmd, 0x0d)
		enter := []byte{0x45, 0x53, 0x43, 0x2F, 0x56, 0x50, 0x2E, 0x6E, 0x65, 0x74, 0x10, 0x03, 0x00, 0x00, 0x00, 0x00}

		_, err := p.writeAndRead(ctx, conn, enter, 5*time.Second, ' ')
		if err != nil {
			log.L.Warnf("there was an error sending the ESC/VP.net command: %v", err)
			return err
		}

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
		log.L.Infof("Setting blank status to %s", blanked)
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

	log.L.Infof("Blank status set to %s", blanked)
	return nil
}
