RM := rm -f

MAKEFILE = $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))

MODULE = $(shell awk '/^module/{print $$2}' go.mod)
BASENAME = $(lastword $(subst /, , $(MODULE)))
VERSION = $(shell cat VERSION)
LDFLAGS = "-X $(MODULE)/internal/cmn.Version=$(VERSION) -X $(MODULE)/internal/cmn.Basename=$(BASENAME)"

all: build

.PHONY: all

build:
	@go build -ldflags $(LDFLAGS) ./cmd/${BASENAME}
.PHONY: build

windows:
	@GOOS=windows
	@GOARCH=386
	@go build -ldflags $(LDFLAGS) -o $(BASENAME).exe ./cmd/${BASENAME}
	@zip -q9 $(BASENAME).v$(VERSION).32bit.windows.zip $(BASENAME).exe
	@$(RM) $(BASENAME).exe
	@GOARCH=amd64
	@go build -ldflags $(LDFLAGS) -o $(BASENAME).exe ./cmd/${BASENAME}
	@zip -q9 $(BASENAME).v$(VERSION).64bit.windows.zip $(BASENAME).exe
	@$(RM) $(BASENAME).exe
.PHONY: windows

linux:
	@GOOS=linux
	@GOARCH=386
	@go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME).v$(VERSION).32bit.linux.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
	@GOARCH=amd64
	@go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME).v$(VERSION).64bit.linux.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
.PHONY: linux

macos:
	@GOOS=darwin
	@GOARCH=amd64
	@go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME).v$(VERSION).64bit.macos.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
	@GOARCH=arm64
	@go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME).v$(VERSION).64bit.macos.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
.PHONY: macos

release: windows linux macos
.PHONY: release

clean:
	@$(RM) $(BASENAME) $(BASENAME).v*
.PHONY: clean
