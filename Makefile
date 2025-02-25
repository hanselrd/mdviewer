BINARIES := mdviewer lobsterzips2parquet
VERSION := 0.0.1-alpha.1
MODULE := github.com/hanselrd/mdviewer

PLATFORMS := windows/amd64 \
             linux/amd64 \
             darwin/amd64 \
             darwin/arm64
BUILDS := debug release

DATE := date
ECHO := echo
FIND := find
GIT := git
GO := go
GOFUMPT := gofumpt
GOIMPORTS := goimports
GOLINES := golines
MKDIR := mkdir
RM := rm

BUILD.VERSION := $(VERSION)
BUILD.HASH := $(shell $(GIT) rev-parse HEAD)
BUILD.TIME := $(shell $(DATE) --utc "+%Y-%m-%dT%H:%M:%SZ")
BUILD.DIRTY := $(shell $(GIT) diff --quiet || $(ECHO) "1")

GCFLAGS.debug := all=-N -l
GCFLAGS.release := all=-B -C
GCFLAGS :=
LDFLAGS.debug :=
LDFLAGS.release := -s -w
LDFLAGS := -X '$(MODULE)/internal/build.Version=$(BUILD.VERSION)' \
           -X '$(MODULE)/internal/build.Hash=$(BUILD.HASH)' \
           -X '$(MODULE)/internal/build.Time=$(BUILD.TIME)' \
           -X '$(MODULE)/internal/build.Dirty=$(BUILD.DIRTY)'

CMDDIR := cmd
SOURCEDIRS := pkg internal
DATADIR := data
BINDIR := bin

SOURCES := $(shell $(FIND) $(SOURCEDIRS) -type f -name "*.go")
BINARIES2 := $(foreach PLATFORM,$(PLATFORMS), \
	$(foreach BINARY,$(BINARIES), \
		$(BINARY)_$(subst /,_,$(PLATFORM))$(if $(findstring windows,$(PLATFORM)),.exe,) \
	) \
)

.PHONY: build
build: $(addprefix build-,$(BUILDS))

define build-BUILD
.PHONY: build-$(1)
build-$(1): $(addprefix $(BINDIR)/$(1)/,$(BINARIES2))
endef
$(foreach BUILD,$(BUILDS),$(eval $(call build-BUILD,$(BUILD))))

define build-BUILD-BINARY
$(BINDIR)/$(1)/$(2): GOOS := $(word 2,$(subst _, ,$(2)))
$(BINDIR)/$(1)/$(2): GOARCH := $(word 3,$(subst _, ,$(basename $(2))))
$(BINDIR)/$(1)/$(2): GCFLAGS += $(GCFLAGS.$(1))
$(BINDIR)/$(1)/$(2): LDFLAGS += $(LDFLAGS.$(1))
$(BINDIR)/$(1)/$(2): $(SOURCES) $(shell $(FIND) $(CMDDIR)/$(word 1,$(subst _, ,$(2))) -type f -name "*.go") | $(BINDIR)/$(1)
	GOOS=$$(GOOS) GOARCH=$$(GOARCH) $(GO) build -gcflags="$$(GCFLAGS)" -ldflags="$$(LDFLAGS)" -o $$@ ./$(CMDDIR)/$(word 1,$(subst _, ,$(2)))
endef
$(foreach BUILD,$(BUILDS), \
	$(foreach BINARY,$(BINARIES2), \
		$(eval $(call build-BUILD-BINARY,$(BUILD),$(BINARY))) \
	) \
)

$(addprefix $(BINDIR)/,$(BUILDS)):
	$(MKDIR) -p $@

.PHONY: format
format:
	@# $(GO) fmt ./...
	$(GOIMPORTS) -w -local "$(MODULE)" .
	$(GOFUMPT) -w -extra .
	$(GOLINES) -w -m 80 $(CMDDIR) $(SOURCEDIRS)

.PHONY: tidy
tidy:
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: clean
clean:
	$(GO) clean
	$(RM) -rfv $(BINDIR) $(DATADIR)/*.parquet $(DATADIR)/*.csv $(DATADIR)/*.txt
