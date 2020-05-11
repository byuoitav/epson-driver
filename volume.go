package epson

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/connpool"
)

var currentVolume int

func (p *Projector) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	var volume int

	work := func(conn connpool.Conn) error {
		log.L.Infof("Getting volume")

		cmd := []byte("VOL?")
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error getting the volume: %v", err)
		}
		checker = strings.Split(checker, "=")[1]
		checker = strings.Split(checker, "\r")[0]

		//convert and divide by 12 because they have it on a scale of 0-255
		num, err := strconv.Atoi(checker)
		if err != nil {
			log.L.Warnf("Error converting to int %v\n", err)
		}

		num = num / 12

		log.L.Infof("The volume level is %v", num)

		if num > 20 || num < 0 {
			log.L.Warnf("Volume out of range: %d", num)
			return fmt.Errorf("volume out of range: %d", num)
		}

		volume = num
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return volume, err
	}
	currentVolume = volume

	return volume, nil
}

func (p *Projector) SetVolumeByBlock(ctx context.Context, block string, volume int) error {
	work := func(conn connpool.Conn) error {
		log.L.Infof("Setting volume to %d", volume)

		word := "VOL "
		bigVolume := volume*12 + 3
		newVolume := strconv.Itoa(bigVolume)
		word += newVolume
		cmd := []byte(word)
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error setting the volume: %v", err)
		}

		bytes := fmt.Sprintf("%x", checker)

		if bytes != "3a" {
			return fmt.Errorf("There was an error executing the command - %s", bytes)
		}

		log.L.Infof("Volume set to %d", volume)
		currentVolume = volume
		return nil

	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	return nil
}

func (p *Projector) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	var muted bool

	work := func(conn connpool.Conn) error {
		log.L.Infof("Getting mute")

		cmd := []byte("VOL?")
		cmd = append(cmd, 0x0d)
		checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
		if err != nil {
			return fmt.Errorf("There was an error getting the volume: %v", err)
		}
		checker = strings.Split(checker, "=")[1]
		checker = strings.Split(checker, "\r")[0]

		num, err := strconv.Atoi(checker)
		if err != nil {
			log.L.Warnf("Error converting to int %v\n", err)
		}
		muted = (num == 0)
		return nil
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return muted, err
	}

	log.L.Infof("Mute status is %s", muted)
	return muted, nil
}

func (p *Projector) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	work := func(conn connpool.Conn) error {
		log.L.Infof("Setting mute to %s", muted)
		switch muted {
		case true:
			cmd := []byte("VOL 0")
			cmd = append(cmd, 0x0d)
			checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
			if err != nil {
				return fmt.Errorf("There was an error muting: %v", err)
			}

			bytes := fmt.Sprintf("%x", checker)

			if bytes != "3a" {
				return fmt.Errorf("There was an error executing the command - %s", bytes)
			}

			return nil
		case false:
			str := "VOL " + strconv.Itoa(currentVolume*12+3)
			cmd := []byte(str)
			cmd = append(cmd, 0x0d)
			checker, err := p.writeAndRead(ctx, conn, cmd, 5*time.Second, ':')
			if err != nil {
				return fmt.Errorf("There was an error unmuting: %v", err)
			}

			bytes := fmt.Sprintf("%x", checker)

			if bytes != "3a" {
				return fmt.Errorf("There was an error executing the command - %s", bytes)
			}

			return nil
		default:
			return fmt.Errorf("unexpected mute state '%v'", muted)
		}
	}

	err := p.Pool.Do(ctx, work)
	if err != nil {
		return err
	}

	log.L.Infof("Mute status is %s", muted)

	return nil
}
