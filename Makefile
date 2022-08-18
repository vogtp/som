.PHONY: build
build: generate build-stater build-monitor-cdp build-visualiser build-alerter build-somctl

# add build date and time to version
curdate=$(shell date -u +%Y%m%d-%H%M)
build_date = -ldflags "-X  gitlab-int.its.unibas.ch/vogtp/autochrome/pkg/core/cfg.BuildInfo=$(curdate)"

.PHONY: generate
generate:
	rm pkg/visualiser/webstatus/README.md 
	ln README.md pkg/visualiser/webstatus/README.md 
	go generate ./...

.PHONY: build-stater
build-stater:
	go build $(build_date) -tags prod,inituserstore -o ./build/ ./cmd/components/stater/ 
	mv build/stater build/som.stater

.PHONY: build-somctl
build-somctl:
	go build $(build_date) -tags prod -o ./build/ ./cmd/somctl/
	
build-%: 
	go build $(build_date) -tags prod -o ./build/ ./cmd/components/$*/
	mv build/$* build/som.$* 