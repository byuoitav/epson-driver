package epson

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

func (p *Projector) GetPower(ctx context.Context) (string, error) {
	var power string

	work := func(conn connpool.Conn) error {
		log.L.Infof("Getting power state")

		cmd := []byte("PWR?")
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error getting power status: %v", err)
		}

		checker = strings.TrimSuffix(checker, "\r:")
		log.L.Debug("This is checker " + checker)

		switch checker {
		case "PWR=00":
			// Standby
			power = "standby"
		case "PWR=01":
			// On
			power = "on"
		case "PWR=02":
			// Warming up
			power = "warming up"
		case "PWR=03":
			// Cooling down
			power = "cooling down"
		case "PWR=04":
			// Standby (network offline)
			power = "standby (network offline)"
		case "PWR=05":
			// Standby (abnormal)
			power = "standby (abnormal)"
		case "PWR=09":
			// Standby (A/V standby)
			power = "standby"
		default:
			return fmt.Errorf("unknown power state '%s'", checker)
		}
		log.L.Debugf("received response: %s\n", checker)
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return power, err
	}

	return power, nil
}

func (p *Projector) SetPower(ctx context.Context, power string) error {
	switch power {
	case "standby":
		power = "off"
	}
	work := func(conn connpool.Conn) error {
		var cmd []byte

		switch power {
		case "on":
			cmd = []byte("PWR ON")
		case "off":
			cmd = []byte("PWR OFF")
		default:
			return fmt.Errorf("unexpected power state: %v", power)
		}
		cmd = append(cmd, 0x0d)

		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error setting the power status: %v", err)
		}

		bytes := fmt.Sprintf("%x", checker)

		if bytes != "3a" {
			return fmt.Errorf("There was an error executing the command - %s", bytes)
		}

		log.L.Infof("Power state changed: %v", power)
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	//TODO is the sleep still necessary???
	time.Sleep(25 * time.Second)
	return nil
}
