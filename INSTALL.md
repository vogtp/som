# Install SOM

## Prepair

### Create Repository

Add som as dependency `go get github.com/vogtp/som/`

Copy the following files and directories from the SOM repository:
1. cmd/
2. init/
3. .gitignore
4. som_example.yml as som.yml
5. Makefile

### Create Packages

Create a directory / go package (custom in this example)

#### Add key for secret store

Create a go file to load key for the store.  Example below.

custom/ignore_key.go:

`` 
package custom

import (
	"github.com/vogtp/som/pkg/monitor/szenario"
	"github.com/vogtp/som/szenarios"
)

func init() {
    // recommended key length: 40
	core.Keystore.Add([]byte("CHANGE_ME"))
}
``

#### Create Szenario Loader

On how to create custom szenarios see README.md.

custom/loader.go
``
package custom

import (
	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Load the szenarios and return the config
func Load() *szenario.Config {
	szConfig := szenario.New()
	userTypeWorld := szenario.MustUserType(szConfig.CreateUsertType("world", "World contains szenarios accessible without password"))
	userTypeStaf := szenario.MustUserType(szConfig.CreateUsertType("staf", "Staf contains szenarios relevant for staf members"))

	szConfig.Add(
		"google",
		&szenarios.GoogleSzenario{Base: &szenario.Base{}},
		[]*szenario.UserType{userTypeWorld, userTypeStaf},
	)
	szConfig.Add(
		"OWA", // Outlook Web Access
		&szenarios.OwaSzenario{Base: &szenario.Base{},
			OwaURL: "http://mail.MY-COMPANY.com",
		},
		[]*szenario.UserType{userTypeStaf},
	)

	return szConfig
}
``

This example load the SOM example szenario (OWA and google) and asociate them with usertypes.
Usertype are used to asociate users to szenarios.  For more information see README.md.

#### Load Szenarios

Under `cmd/` change all occuences of `szenarios.Load()` to `custom.Load()`.

#### Prometheus

Install prometheus and add the files from `init/prometheus` to `/etc/prometheus/consoles` and `/etc/prometheus/console_libraries`.

Add a prometheus job as described in `init/prometheus/job.yml`

#### Config

In the copied `som.yml` change the values to match your setup.

## Build

See Makefile: `make build`

1. Build binaries 
1.1 Components (1.1.1 or 1.1.2)
1.1.1 All binaries in cmd/components (except allinoe) recommended
1.1.2 Build allinone 

Copy the `som.*` binaries to `/srv/som`

## Add Users

Run `somctl user add` to add users.

## Install systemd services units

Copy the services files from `init/` to `/etc/systemd/system/`.
`som.monitor@.service` can be linked to `som.monitor@USERNAME.service` to start a monitor with the according user.  More than one user monitor can be started.

## Start

For all systemd services files (except `som.monitor@.service` without  username) run:

`systemctl start SERVICE`

and

`systemctl enable SERVICE`

Connect to http://localhost:8083/ and enjoy.