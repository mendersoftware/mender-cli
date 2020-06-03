GO ?= go
GOFMT ?= gofmt
V ?=
PKGS = $(shell go list ./...)
PKGFILES = $(shell find . \( -path ./vendor -o -path ./Godeps \) -prune \
		-o -type f -name '*.go' -print)
PKGFILES_notest = $(shell echo $(PKGFILES) | tr ' ' '\n' | grep -v _test.go)
GOCYCLO ?= 15

TOOLS = \
	github.com/fzipp/gocyclo \
	github.com/opennota/check/cmd/varcheck \
	github.com/mendersoftware/deadcode

VERSION = $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)

GO_LDFLAGS = \
	-ldflags "-X main.Version=$(VERSION)"

ifeq ($(V),1)
BUILDV = -v
endif

TAGS =
ifeq ($(LOCAL),1)
TAGS += local
endif

ifneq ($(TAGS),)
BUILDTAGS = -tags '$(TAGS)'
endif

build:
	CGO_ENABLED=0 $(GO) build $(GO_LDFLAGS) $(BUILDV) $(BUILDTAGS)

build-multiplatform:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GO_LDFLAGS) $(BUILDV) $(BUILDTAGS) \
	     -o mender-cli.linux.amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GO_LDFLAGS) $(BUILDV) $(BUILDTAGS) \
	     -o mender-cli.darwin.amd64

install:
	CGO_ENABLED=0 $(GO) install $(GO_LDFLAGS) $(BUILDV) $(BUILDTAGS)

clean:
	$(GO) clean
	rm -f coverage.txt coverage-tmp.txt

get-tools:
	set -e ; for t in $(TOOLS); do \
		echo "-- go getting $$t"; \
		GO111MODULE=off go get -u $$t; \
	done

check: test extracheck

test:
	$(GO) test $(BUILDV) $(PKGS)

extracheck:
	echo "-- checking if code is gofmt'ed"
	if [ -n "$$($(GOFMT) -d $(PKGFILES))" ]; then \
		echo "-- gofmt check failed"; \
		/bin/false; \
	fi
	echo "-- checking with govet"
	$(GO) vet $(PKGS)
	echo "-- checking for dead code"
	deadcode -ignore version.go:Version
	echo "-- checking with varcheck"
	varcheck .
	echo "-- checking cyclometric complexity > $(GOCYCLO)"
	gocyclo -over $(GOCYCLO) $(PKGFILES_notest)

cover: coverage
	$(GO) tool cover -func=coverage.txt

htmlcover: coverage
	$(GO) tool cover -html=coverage.txt

coverage:
	rm -f coverage.txt
	echo 'mode: set' > coverage.txt
	set -e ; for p in $(PKGS); do \
		rm -f coverage-tmp.txt;  \
		$(GO) test -coverprofile=coverage-tmp.txt $$p ; \
		if [ -f coverage-tmp.txt ]; then \
			cat coverage-tmp.txt | grep -v 'mode:' >> coverage.txt || /bin/true; \
		fi; \
	done
	rm -f coverage-tmp.txt

.PHONY: build clean get-tools test check \
	cover htmlcover coverage
