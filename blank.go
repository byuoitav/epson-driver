package epson

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) GetBlanked(ctx context.Context) (bool, error) {
	p.infof("Getting blanked status")

	cmd := []byte("MUTE?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return false, err
	}

	switch {
	case bytes.Contains(resp, []byte("MUTE=ON")):
		p.infof("Blanked status is true")
		return true, nil
	case bytes.Contains(resp, []byte("MUTE=OFF")):
		p.infof("Blanked status is false")
		return false, nil
	default:
		return false, fmt.Errorf("unknown blanked state: %#x", resp)
	}
}

func (p *Projector) SetBlanked(ctx context.Context, blanked bool) error {
	p.infof("Setting blanked to %v", blanked)

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
