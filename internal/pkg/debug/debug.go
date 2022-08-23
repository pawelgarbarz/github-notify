package debug

import "errors"

const NoDebug = 0
const Important = 1
const Detailed = 2
const SuperDetailed = 3

var ErrDebugLevelNotSupported = errors.New("debug level not supported")

type DebugInterface interface {
	Level() int
}

type debug struct {
	level int
}

func NewDebug() *debug {
	return &debug{
		level: NoDebug,
	}
}

func (d *debug) SetLevel(level int) error {
	if level < NoDebug || level > SuperDetailed {
		return ErrDebugLevelNotSupported
	}

	d.level = level

	return nil
}

func (d *debug) Level() int {
	return d.level
}
