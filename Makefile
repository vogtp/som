.PHONY: build
build: generate build-stater build-monitor-cdp build-visualiser build-alerter build-somctl

# add build date and time to version
curdate=$(shell date --iso-8601='minutes')
build_flags = -ldflags "-X  github.com/vogtp/som.BuildInfo=$(curdate)"

GO_CMD=CGO_ENABLED=0 go

.PHONY: install_stinger
install_stinger:
	$(GO_CMD) install golang.org/x/tools/cmd/stringer@latest

.PHONY: generate
generate: install_stinger
	go generate ./...

.PHONY: build-somctl
build-somctl:
	$(GO_CMD) build $(build_flags) -tags prod -o ./build/ ./cmd/somctl/
	
build-%: 
	$(GO_CMD) build $(build_flags) -tags prod -o ./build/ ./cmd/components/$*/
	mv build/$* build/som.$* 
