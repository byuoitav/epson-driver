package epson

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) GetBlank(ctx context.Context) (bool, error) {
	p.infof("Getting blanked status")

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to get blanked: %s", err)
		}
	}()

	// check if it's powered off
	power, err := p.GetPower(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to check power state: %w", err)
	}

	if power == "standby" {
		return p.lastKnownBlanked, nil
	}

	cmd := []byte("MUTE?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return false, err
	}

	switch {
	case bytes.Contains(resp, []byte("MUTE=ON")):
		p.infof("Blanked status is true")
		p.lastKnownBlanked = true
		return true, nil
	case bytes.Contains(resp, []byte("MUTE=OFF")):
		p.infof("Blanked status is false")
		p.lastKnownBlanked = false
		return false, nil
	default:
		return false, fmt.Errorf("unknown blanked state: %#x", resp)
	}
}

func (p *Projector) SetBlank(ctx context.Context, blanked bool) error {
	p.infof("Setting blanked to %v", blanked)

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to set blanked: %s", err)
		}
	}()

	// check if it's powered off
	power, err := p.GetPower(ctx)
	if err != nil {
		return fmt.Errorf("unable to check power state: %w", err)
	}

	if power == "standby" {
		// pretend like it worked
		return nil
	}

	cmd := []byte("MUTE OFF\r")
	if blanked {
		cmd = []byte("MUTE ON\r")
	}

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return err
	}

	if !bytes.Equal(bytes.TrimSpace(resp), []byte{0x3a}) {
		return fmt.Errorf("error from projector: %#x", resp)
	}

	p.infof("Successfully set blanked status")
	return nil
}
