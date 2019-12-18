package epson

import (
	"context"
	"fmt"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

func (p *Projector) GetInput(ctx context.Context) (string, error) {
	var input string

	work := func(conn connpool.Conn) error {

		cmd := []byte("SOURCE?")
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error getting the input: %v", err)
		}

		switch checker {
		case "SOURCE=30\r:":
			input = "HDMI"
		case "SOURCE=A0\r:":
			input = "DVI-D"
		case "SOURCE=11\r:":
			input = "computer"
		case "SOURCE=53\r:":
			input = "LAN"
		case "SOURCE=80\r:":
			input = "HDBaseT"
		case "SOURCE=B1\r:":
			input = "BNC"
		case "SOURCE=60\r:":
			input = "SDI"
		default:
			return fmt.Errorf("unknown source response '%s'", checker)
		}
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return input, err
	}

	return input, nil
}

func (p *Projector) SetInput(ctx context.Context, input string) error {
	var str string
	switch input {
	case "HDMI":
		str = "30"
	case "DVI-D":
		str = "A0"
	case "computer":
		str = "11"
	case "LAN":
		str = "53"
	case "HDBaseT":
		str = "80"
	case "BNC":
		str = "B1"
	case "SDI":
		str = "60"
	default:
		return fmt.Errorf("unknown source input '%s'", input)
	}

	work := func(conn connpool.Conn) error {

		cmd := []byte(fmt.Sprintf("SOURCE %s", str))
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error setting the input: %v", err)
		}

		bytes := fmt.Sprintf("%x", checker)

		if bytes != "3a" {
			return fmt.Errorf("There was an error executing the command - %s", bytes)
		}

		log.L.Infof("input changed: %v", input)
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	//TODO remove?
	time.Sleep(25 * time.Second)
	return nil
}
