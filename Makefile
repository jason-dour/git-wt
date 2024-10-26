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
	@GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BASENAME).exe ./cmd/${BASENAME}
	@zip -q9 $(BASENAME)_v$(VERSION)_windows_amd64.zip $(BASENAME).exe
	@$(RM) $(BASENAME).exe
	@GOOS=windows GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(BASENAME).exe ./cmd/${BASENAME}
	@zip -q9 $(BASENAME)_v$(VERSION)_windows_arm64.zip $(BASENAME).exe
	@$(RM) $(BASENAME).exe
.PHONY: windows

linux:
	@GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME)_v$(VERSION)_linux_amd64.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
	@GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(BASENAME) ./cmd/${BASENAME}
	@tar -zcf $(BASENAME)_v$(VERSION)_linux_arm64.tar.gz $(BASENAME)
	@$(RM) $(BASENAME)
.PHONY: linux

macos:
	@GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BASENAME)_amd64 ./cmd/${BASENAME}
	@GOOS=darwin GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(BASENAME)_arm64 ./cmd/${BASENAME}
	@lipo -create -output ${BASENAME} ${BASENAME}_amd64 ${BASENAME}_arm64
	@tar -zcf $(BASENAME)_v$(VERSION)_macos_universal.tar.gz $(BASENAME)
	@$(RM) $(BASENAME) $(BASENAME)_a*
.PHONY: macos

release: windows linux macos
.PHONY: release

clean:
	@$(RM) $(BASENAME) $(BASENAME).v* $(BASENAME)_*
.PHONY: clean
