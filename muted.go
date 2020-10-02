package epson

import (
	"context"
	"errors"
)

func (p *Projector) GetMutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	return nil, errors.New("not implemented")
}

func (p *Projector) SetMute(ctx context.Context, block string, muted bool) error {
	return errors.New("not implemented")
}
