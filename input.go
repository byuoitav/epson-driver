package epson

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) GetInput(ctx context.Context) (string, error) {
	p.infof("Getting current input")

	cmd := []byte("SOURCE?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return "", err
	}

	var input string

	switch {
	case bytes.Contains(resp, []byte("SOURCE=30")):
		input = "HDMI"
	case bytes.Contains(resp, []byte("SOURCE=A0")):
		input = "DVI-D"
	case bytes.Contains(resp, []byte("SOURCE=11")):
		input = "computer"
	case bytes.Contains(resp, []byte("SOURCE=53")):
		input = "LAN"
	case bytes.Contains(resp, []byte("SOURCE=80")):
		input = "HDBaseT"
	case bytes.Contains(resp, []byte("SOURCE=B1")):
		input = "BNC"
	case bytes.Contains(resp, []byte("SOURCE=60")):
		input = "SDI"
	default:
		return "", fmt.Errorf("unknown input: %#x", resp)
	}

	p.infof("Current input is %s", input)
	return input, nil
}

func (p *Projector) SetInput(ctx context.Context, input string) error {
	p.infof("Setting input to %v", input)

	var cmd []byte
	switch input {
	case "HDMI":
		cmd = []byte("SOURCE 30\r")
	case "DVI-D":
		cmd = []byte("SOURCE A0\r")
	case "computer":
		cmd = []byte("SOURCE 11\r")
	case "LAN":
		cmd = []byte("SOURCE 53\r")
	case "HDBaseT":
		cmd = []byte("SOURCE 80\r")
	case "BNC":
		cmd = []byte("SOURCE B1\r")
	case "SDI":
		cmd = []byte("SOURCE 60\r")
	default:
		return fmt.Errorf("invalid input")
	}

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return err
	}

	if !bytes.Equal(bytes.TrimSpace(resp), []byte{0x3a}) {
		return fmt.Errorf("error from projector: %#x", resp)
	}

	p.infof("Successfully set input")
	return nil
}
