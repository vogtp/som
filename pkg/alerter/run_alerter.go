package alerter

import "github.com/vogtp/som/pkg/core"

// Run the alerter
func Run(name string, coreOpts ...core.Option) (func(), error) {
	_, close := core.New(name, coreOpts...)
	parseConfig()
	NewMailer()
	NewTeams()
	return close, nil
}

func parseConfig() {
	parseDestinations()
	parseRules()
}
