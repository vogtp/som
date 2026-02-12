.PHONY: test build
build: generate build-stater build-monitor-cdp build-visualiser build-alerter build-somctl build-checkctl

BRANCH=$(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')

# add build date and time to version
curdate=$(shell date --iso-8601='minutes')
build_flags = -ldflags "-X github.com/vogtp/som.BuildInfo=$(curdate) -X github.com/vogtp/som.Branch=$(BRANCH)"

GO_TAGS=
# embedd_nats embedds a nats bus
#GO_TAGS=-tags=embedd_nats

GO_CMD=CGO_ENABLED=0 go 

.PHONY: install_stinger
install_stinger:
	$(GO_CMD) install golang.org/x/tools/cmd/stringer@latest

.PHONY: generate
generate: install_stinger
	go generate ./...

.PHONY: build-somctl
build-somctl:
	$(GO_CMD) build $(build_flags) $(GO_TAGS) -o ./build/ ./cmd/somctl/

.PHONY: build-checkctl
build-checkctl:
	$(GO_CMD) build $(build_flags) $(GO_TAGS) -o ./build/ ./cmd/checkctl/
	
build-%: 
	$(GO_CMD) build $(build_flags) $(GO_TAGS) -o ./build/ ./cmd/components/$*/
	mv build/$* build/som.$* 

.PHONY: test
test:
	echo "Running test (excluding packages without tests)"
	go test $(GO_TAGS) ./... | grep -v "no test files"