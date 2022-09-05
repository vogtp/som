.PHONY: build
build: generate build-stater build-monitor-cdp build-visualiser build-alerter build-somctl

# add build date and time to version
curdate=$(shell date -u +%Y%m%d-%H%M)
build_date = -ldflags "-X  github.com/vogtp/som/pkg/core/cfg.BuildInfo=$(curdate)"

.PHONY: generate
generate:
	go generate ./...

.PHONY: build-somctl
build-somctl:
	go build $(build_date) -tags prod -o ./build/ ./cmd/somctl/
	
build-%: 
	go build $(build_date) -tags prod -o ./build/ ./cmd/components/$*/
	mv build/$* build/som.$* 
