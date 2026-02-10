package visualiser

import (
	"github.com/vogtp/som/pkg/bridger"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/visualiser/webstatus"
)

// Run the visualiser
func Run(name string, coreOpts ...core.Option) (func(), error) {
	cfg.Parse()
	_, close,err := core.New(name, coreOpts...)
	if err != nil {
		return nil, err
	}
	bridger.RegisterPrometheus()
	webstatus.New()
	return close, nil
}
