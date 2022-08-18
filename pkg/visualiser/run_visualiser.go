package visualiser

import (
	"github.com/vogtp/som/pkg/bridger"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/visualiser/webstatus"
)

// Run the visualiser
func Run(name string, coreOpts ...core.Option) (func(), error) {
	_, close := core.New(name, coreOpts...)
	bridger.RegisterPrometheus()
	webstatus.New()
	return close, nil
}
