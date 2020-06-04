package epson

// Logger is a logger usable by Projector.
type Logger interface {
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

func (d *Projector) debugf(format string, a ...interface{}) {
	if d.logger != nil {
		d.logger.Debugf(format, a...)
	}
}

func (d *Projector) infof(format string, a ...interface{}) {
	if d.logger != nil {
		d.logger.Infof(format, a...)
	}
}

func (d *Projector) warnf(format string, a ...interface{}) {
	if d.logger != nil {
		d.logger.Warnf(format, a...)
	}
}

func (d *Projector) errorf(format string, a ...interface{}) {
	if d.logger != nil {
		d.logger.Errorf(format, a...)
	}
}
