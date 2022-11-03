package visualiser

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
)

// Dumper writes all files to disk
type Dumper struct {
	hcl       hcl.Logger
	outFolder string
}

// NewDumper registers a Dumper on the event bus
func NewDumper() {
	bus := core.Get().Bus()
	d := Dumper{
		hcl:       bus.GetLogger().Named("dumper"),
		outFolder: fmt.Sprintf("%s/dump/", viper.GetString(cfg.DataDir)),
	}
	if err := core.EnsureOutFolder(d.outFolder); err != nil {
		d.hcl.Warnf("there is no outfolder: %v", err)
	}
	bus.Szenario.Handle(d.handleSzenarioEvt)
	d.hcl.Infof("Will save dumps to %s", d.outFolder)
}

func (d *Dumper) handleSzenarioEvt(e *msg.SzenarioEvtMsg) {
	for _, f := range e.Files {
		name := fmt.Sprintf("%s/%s.%s", d.outFolder, f.Name, f.Type.Ext)
		d.hcl.Infof("Writing %s", name)
		if err := os.WriteFile(name, f.Payload, 0644); err != nil {
			d.hcl.Warnf("cannot write file: %v", err)
		}
	}
	// TODO add time
}
