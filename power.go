package epson

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

func (p *Projector) GetPower(ctx context.Context) (string, error) {
	p.infof("Getting power status")

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to get power: %s", err)
		}
	}()

	cmd := []byte("PWR?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return "", err
	}

	switch {
	case bytes.Contains(resp, []byte("PWR=01")):
		p.infof("Power status is 'on'")
		return "on", nil
	case bytes.Contains(resp, []byte("PWR=00")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=02")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=03")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=04")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=05")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=05")):
		fallthrough
	case bytes.Contains(resp, []byte("PWR=09")):
		p.infof("Power status is 'standby'")
		return "standby", nil
	default:
		return "", fmt.Errorf("unknown power state: %#x", resp)
	}
}

func (p *Projector) SetPower(ctx context.Context, power string) error {
	p.infof("Setting power to %v", power)

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to set power: %s", err)
		}
	}()

	cmd := []byte("PWR OFF\r")
	if power == "on" {
		cmd = []byte("PWR ON\r")
	}

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return err
	}

	if !bytes.Equal(bytes.TrimSpace(resp), []byte{0x3a}) {
		return fmt.Errorf("error from projector: %#x", resp)
	}

	time.Sleep(4000 * time.Millisecond)
	p.infof("Successfully set power status")
	return nil
}
