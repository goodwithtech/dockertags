.PHONY: \
  lint \
  vet \
  fmt \
  pretest \
  test \

PKGS = $(shell go list ./... | grep -v /vendor/)
SRCS = $(shell git ls-files '*.go')
GO := GO111MODULE=on go
GO_OFF := GO111MODULE=off go

lint:
	$(GO_OFF) get -u golang.org/x/lint/golint
	$(foreach file,$(PKGS),golint -set_exit_status $(file) || exit;)
vet:
	echo $(PKGS) | xargs env $(GO) vet || exit;

fmt:
	$(foreach file,$(SRCS),gofmt -s -d $(file);)

pretest: lint vet fmt

test:
	$(foreach file,$(PKGS),$(GO) test -v $(file) || exit;)