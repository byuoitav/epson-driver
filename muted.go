package epson

import (
	"context"
	"errors"
)

func (p *Projector) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	return false, errors.New("not implemented")
}

func (p *Projector) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	return errors.New("not implemented")
}
