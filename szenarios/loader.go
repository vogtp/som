package szenarios

import (
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Load the szenarios and return the config
func Load() *szenario.Config {
	szConfig := szenario.New()
	userTypeWorld := szenario.MustUserType(szConfig.CreateUsertType("world", "World contains szenarios accessable without password"))
	userTypeStaf := szenario.MustUserType(szConfig.CreateUsertType("staf", "Staf contains szenarios relevant for staf members"))

	szConfig.Add(
		"google",
		&googleSzenario{&szenario.Base{}},
		[]*szenario.UserType{userTypeWorld, userTypeStaf},
	)
	szConfig.Add(
		"OWA",
		&owaSzenario{Base: &szenario.Base{},
			owaUrl: "http://mail.example.com",
		},
		[]*szenario.UserType{userTypeStaf},
	)
	szConfig.Add(
		"crasher",
		&crasherSzenario{&szenario.Base{}},
		[]*szenario.UserType{},
	)
	return szConfig
}
