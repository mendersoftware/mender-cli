GO ?= go
V ?=
PKGS = $(shell go list ./...)
PKGFILES = $(shell find . \( -path ./vendor -o -path ./Godeps \) -prune \
		-o -type f -name '*.go' -print)
PKGFILES_notest = $(shell echo $(PKGFILES) | tr ' ' '\n' | grep -v _test.go)

GO_TEST_TOOLS = \
	github.com/opennota/check/cmd/varcheck \
	github.com/mendersoftware/deadcode

VERSION = $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)

GO_LDFLAGS = \
	-ldflags "-X github.com/mendersoftware/mender-cli/cmd.Version=$(VERSION)"
BUILDFLAGS ?=

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
	CGO_ENABLED=0 $(GO) build $(BUILDFLAGS) $(GO_LDFLAGS) $(BUILDV) $(BUILDTAGS)

build-autocomplete-scripts: build
	@./mender-cli --generate-autocomplete

build-multiplatform:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(BUILDTAGS) $(GO_LDFLAGS) $(BUILDV) $(BUILDFLAGS) \
	     -o mender-cli.linux.amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(BUILDTAGS) $(GO_LDFLAGS) $(BUILDV) $(BUILDFLAGS) \
	     -o mender-cli.darwin.amd64

build-coverage:
	CGO_ENABLED=0 $(GO) build -cover -o mender-cli-test \
		-coverpkg $(shell echo $(PKGS) | tr  ' ' ',')

install:
	CGO_ENABLED=0 $(GO) install $(BUILDTAGS) $(GO_LDFLAGS) $(BUILDV) $(BUILDFLAGS)

install-autocomplete-scripts: build-autocomplete-scripts
	@echo "Installing Bash auto-complete script into ${DESTDIR}${PREFIX}/etc/bash_completion.d/"
	@install -d ${DESTDIR}$(PREFIX)/etc/bash_completion.d/
	@install -m 644 ./autocomplete/autocomplete.sh $(DESTDIR)$(PREFIX)/etc/bash_completion.d/
	@if which zsh >/dev/null 2>&1 ; then \
	echo "Installing zsh auto-complete script into ${DESTDIR}${PREFIX}/usr/local/share/zsh/site-functions/" && \
	install -d $(DESTDIR)$(PREFIX)/usr/local/share/zsh/site-functions/ && \
	install -m 644 ./autocomplete/autocomplete.zsh $(DESTDIR)$(PREFIX)/usr/local/share/zsh/site-functions/_mender-cli \
	; fi

clean:
	$(GO) clean
	rm -f coverage.txt coverage-tmp.txt

get-go-tools:
	set -e ; for t in $(GO_TEST_TOOLS); do \
		echo "-- go getting $$t"; \
		GO111MODULE=off go get -u $$t; \
	done

get-build-deps:
	apt-get update -qq
	apt-get install -yyq $(shell cat deb-requirements.txt)

get-deps: get-go-tools get-build-deps

test-unit:
	$(GO) test $(BUILDV) $(PKGS)

build-acceptance:
	docker compose -f tests/docker-compose.yml build acceptance $(BUILDFLAGS)

run-acceptance:
	docker compose -f tests/docker-compose.yml run acceptance

test-static:
	echo "-- checking for dead code"
	deadcode -ignore version.go:Version
	echo "-- checking with varcheck"
	varcheck .

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

.PHONY: build clean get-go-tools get-apt-deps get-deps test check \
	cover htmlcover coverage
