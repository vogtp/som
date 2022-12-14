package szenarios

import (
	"github.com/vogtp/som/pkg/monitor/szenario"
)

func init() {
	// IMPORTANT this should be changed and done in a file which is not checked in to the repository
	// core.Keystore.Add([]byte("CHANGE_ME"))
}

// Load the szenarios and return the config
func Load() *szenario.Config {
	szConfig := szenario.New()
	userTypeWorld := szenario.MustUserType(szConfig.CreateUserType("world", "World contains szenarios accessible without password"))
	userTypeStaf := szenario.MustUserType(szConfig.CreateUserType("staf", "Staf contains szenarios relevant for staf members"))

	szConfig.Add(
		"google",
		&GoogleSzenario{Base: &szenario.Base{},
			Search: "SOM",
		},
		[]*szenario.UserType{userTypeWorld, userTypeStaf},
	)
	szConfig.Add(
		"OWA",
		&OwaSzenario{
			Base: &szenario.Base{
				LoginRetry: 4,
			},
			OwaURL: "http://mail.example.com",
		},
		[]*szenario.UserType{userTypeStaf},
	)
	szConfig.Add(
		"crasher",
		&CrasherSzenario{Base: &szenario.Base{}},
		[]*szenario.UserType{},
	)
	return szConfig
}
