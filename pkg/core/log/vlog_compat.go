package log

// Based on vlog code
// so probably a Apache 2.0 License

import (
	"bytes"
	"fmt"
	"io"

	"log/slog"

	"github.com/suborbital/vektor/vlog"
)

// VlogCompat returns a vlog logger backed by slog
func VlogCompat(l *slog.Logger) *vlog.Logger {
	p := vlogProducer{log: l}
	return vlog.New(p, vlog.WithWriter(io.Discard))
}

type vlogProducer struct {
	log *slog.Logger
}

// ErrorString prints a string as an error
func (d vlogProducer) ErrorString(msgs ...interface{}) string {
	d.log.Error(joinInterfaces(msgs...))
	return ""
}

// Error prints a string as an error
func (d vlogProducer) Error(err error) string {
	d.log.Error("Error", Error, err)
	return ""
}

// Warn prints a string as an warning
func (d vlogProducer) Warn(msgs ...interface{}) string {
	d.log.Warn(joinInterfaces(msgs...))
	return ""
}

// Info prints a string as an info message
func (d vlogProducer) Info(msgs ...interface{}) string {
	d.log.Info(joinInterfaces(msgs...))
	return ""
}

// Debug prints a string as debug output
func (d vlogProducer) Debug(msgs ...interface{}) string {
	d.log.Debug(joinInterfaces(msgs...))
	return ""
}

// Trace prints a function name and returns a function to be deferred, logging the completion of a function
func (d vlogProducer) Trace(fnName string) (string, func() string) {
	traceFunc := func() string {
		d.log.Debug(fmt.Sprintf("END %s", fnName), "func", fnName)
		return ""
	}

	d.log.Debug(fmt.Sprintf("START %s", fnName), "func", fnName)
	return "", traceFunc
}

func joinInterfaces(msgs ...interface{}) string {
	msg := ""

	for _, m := range msgs {
		switch elem := m.(type) {
		case string:
			msg += fmt.Sprintf(" %s", elem)
		case uint, uint8, uint16, uint32, int, int8, int16, int32, int64, float32, float64, complex64, complex128:
			buf := &bytes.Buffer{}
			fmt.Fprint(buf, elem)
			msg += " " + buf.String()
		case vlog.SafeStringer:
			msg += " " + elem.SafeString()
		default:
			msg += fmt.Sprintf(" %v", elem)
		}
	}

	// get rid of that first space
	return msg[1:]
}
