package epson

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) GetAudioVideoInputs(ctx context.Context) (map[string]string, error) {
	p.infof("Getting current input")

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to get input: %s", err)
		}
	}()

	// check if it's powered off
	power, err := p.GetPower(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to check power state: %w", err)
	}

	if power == "standby" {
		// pretend like the default input is HDBaseT?
		return map[string]string{"": p.lastKnownInput}, nil
	}

	cmd := []byte("SOURCE?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return nil, err
	}

	var input string

	switch {
	case bytes.Contains(resp, []byte("SOURCE=30")):
		input = "hdmi1"
	case bytes.Contains(resp, []byte("SOURCE=A0")):
		input = "hdmi2"
	case bytes.Contains(resp, []byte("SOURCE=C0")):
		input = "hdmi3"
	case bytes.Contains(resp, []byte("SOURCE=10")):
		input = "computer"
	case bytes.Contains(resp, []byte("SOURCE=52")):
		input = "usb1"
	case bytes.Contains(resp, []byte("SOURCE=54")):
		input = "usb2"
	case bytes.Contains(resp, []byte("SOURCE=53")):
		input = "lan"
	case bytes.Contains(resp, []byte("SOURCE=56")):
		input = "screenmirroring1"
	case bytes.Contains(resp, []byte("SOURCE=80")):
		input = "hdbaset"
	case bytes.Contains(resp, []byte("SOURCE=B1")):
		input = "bnc"
	case bytes.Contains(resp, []byte("SOURCE=60")):
		input = "sdi"
	default:
		return nil, fmt.Errorf("unknown input: %#x", resp)
	}

	p.infof("Current input is %s", input)
	p.lastKnownInput = input
	return map[string]string{"": input}, nil
}

func (p *Projector) SetAudioVideoInput(ctx context.Context, output, input string) error {
	p.infof("Setting input to %v", input)

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to set input: %s", err)
		}
	}()

	var cmd []byte
	switch input {
	case "hdmi1":
		cmd = []byte("SOURCE 30\r")
	case "hdmi2":
		cmd = []byte("SOURCE A0\r")
	case "hdmi3":
		cmd = []byte("SOURCE C0\r")
	case "computer":
		cmd = []byte("SOURCE 10\r")
	case "usb1":
		cmd = []byte("SOURCE 52\r")
	case "usb2":
		cmd = []byte("SOURCE 54\r")
	case "lan":
		cmd = []byte("SOURCE 53\r")
	case "screenmirroring1":
		cmd = []byte("SOURCE 56\r")
	case "hdbaset":
		cmd = []byte("SOURCE 80\r")
	case "bnc":
		cmd = []byte("SOURCE B1\r")
	case "sdi":
		cmd = []byte("SOURCE 60\r")
	default:
		return fmt.Errorf("invalid input")
	}

	// check if it's powered off
	power, err := p.GetPower(ctx)
	if err != nil {
		return fmt.Errorf("unable to check power state: %w", err)
	}

	if power == "standby" {
		// pretend like it worked
		return nil
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
