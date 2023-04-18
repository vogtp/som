package visualiser

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
	"golang.org/x/exp/slog"
)

// Dumper writes all files to disk
type Dumper struct {
	log       *slog.Logger
	outFolder string
}

// NewDumper registers a Dumper on the event bus
func NewDumper() {
	bus := core.Get().Bus()
	d := Dumper{
		log:       bus.GetLogger().With(log.Component, "dumper"),
		outFolder: fmt.Sprintf("%s/dump/", viper.GetString(cfg.DataDir)),
	}
	if err := core.EnsureOutFolder(d.outFolder); err != nil {
		d.log.Warn("there is no outfolder", log.Error, err)
	}
	bus.Szenario.Handle(d.handleSzenarioEvt)
	d.log.Info("Will save dumps to " + d.outFolder)
}

func (d *Dumper) handleSzenarioEvt(e *msg.SzenarioEvtMsg) {
	for _, f := range e.Files {
		name := fmt.Sprintf("%s/%s.%s", d.outFolder, f.Name, f.Type.Ext)
		d.log.Info("Writing %s" + name)
		if err := os.WriteFile(name, f.Payload, 0644); err != nil {
			d.log.Warn("cannot write file", log.Error, err)
		}
	}
	// TODO add time
}
