package epson

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Scale on projector is 0-243
const _volumeConversionFactor = 243.0 / 100.0

func (p *Projector) GetVolumes(ctx context.Context, blocks []string) (map[string]int, error) {
	p.infof("Getting volume")

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to get volume: %s", err)
		}
	}()

	cmd := []byte("VOL?\r")

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return nil, err
	}

	p.infof("volume response: %s", resp)

	// Trim trailing "\r:"
	str := strings.TrimSpace(strings.TrimRight(string(resp), ":"))
	split := strings.Split(str, "=")
	if len(split) != 2 {
		return nil, fmt.Errorf("unexpected response from projector: %#x", resp)
	}

	num, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("unable to convert %q to int: %w", split[0], err)
	}

	// Convert from 0-255 scale to 0-100 scale
	num = int(math.Round(float64(num) / _volumeConversionFactor))

	p.infof("Volume is %v", num)
	return map[string]int{"": num}, nil
}

func (p *Projector) SetVolume(ctx context.Context, block string, volume int) error {
	p.infof("Setting volume to %v", volume)

	// Convert from 0-100 scale to 0-255 scale
	volume = int(math.Round(float64(volume) * _volumeConversionFactor))

	var err error
	defer func() {
		if err != nil {
			p.warnf("unable to set volume: %s", err)
		}
	}()

	cmd := []byte(fmt.Sprintf("VOL %v\r", volume))

	resp, err := p.sendCommand(ctx, cmd, ':')
	if err != nil {
		return err
	}

	if !bytes.Equal(bytes.TrimSpace(resp), []byte{0x3a}) {
		return fmt.Errorf("error from projector: %#x", resp)
	}

	p.infof("Successfully set volume status")
	return nil
}
