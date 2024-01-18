# Makefile for building the sf command

# Binary output name
BINARY_NAME=sf

# Source directory
SRC_DIR=./cmd/sf

# Output directory
BIN_DIR=./bin

# Create output directory if it doesn't exist
$(shell mkdir -p $(BIN_DIR))


# Function to extract version numbers
extract_version_parts = $(shell echo $(1) | awk -F. '{print $$1,$$2,$$3}')


# Function to increment version
increment_version = $(shell v=($(1)); \
    if [ $(2) -eq 1 ]; then \
        echo $$((v[0]+1)).0.0; \
    elif [ $(2) -eq 2 ]; then \
        echo $${v[0]}.$$((v[1]+1)).0; \
    elif [ $(2) -eq 3 ]; then \
        echo $${v[0]}.$${v[1]}.$$((v[2]+1)); \
    fi)

.PHONY: build build-all clean

GIT_TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date +%FT%T%z 2>/dev/null || echo "unknown")

# Determine the current operating system
ifeq ($(OS),Windows_NT)
	CURRENT_OS = windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		CURRENT_OS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		CURRENT_OS = darwin
	endif
endif


# Determine the current architecture
ifeq ($(shell uname -m),x86_64)
    CURRENT_ARCH = amd64
endif
ifeq ($(shell uname -m),aarch64)
    CURRENT_ARCH = arm64
endif

ifeq ($(shell uname -m),arm64)
    CURRENT_ARCH = arm64
endif


build: build-$(CURRENT_OS)-$(CURRENT_ARCH)
	@cp $(BIN_DIR)/$(BINARY_NAME)-$(CURRENT_OS)-$(CURRENT_ARCH) $(BIN_DIR)/$(BINARY_NAME)

# Debugging output
# Debugging output
debug:
	@echo "GOARCH: $(CURRENT_ARCH)"
	@echo "CURRENT_OS: $(CURRENT_OS)"
	@echo "GIT_TAG: $(GIT_TAG)"
	@echo "GIT_COMMIT: $(GIT_COMMIT)"
	@echo "BUILD_TIME: $(BUILD_TIME)"

build-all: build-linux-amd64 build-windows-amd64 build-darwin-amd64 build-linux-arm64 build-darwin-arm64

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(SRC_DIR)
	@echo "Linux AMD64 build complete."

build-windows-amd64:
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(SRC_DIR)
	@echo "Windows AMD64 build complete."

build-darwin-amd64:
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(SRC_DIR)
	@echo "macOS AMD64 build complete."

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 $(SRC_DIR)
	@echo "Linux ARM64 build complete."

build-darwin-arm64:
	@echo "Building for macOS ARM64 (M1)..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(SRC_DIR)
	@echo "macOS ARM64 (M1) build complete."

clean:
	@echo "Cleaning..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)*
	@echo "Clean complete"


release-patch: clean build-all
	@echo "GIT_TAG: $(GIT_TAG)"
	$(eval MAJOR_MINOR_PATCH = $(call extract_version_parts,$(GIT_TAG:v%=%)))
	$(eval NEW_GIT_TAG = $(call increment_version,$(MAJOR_MINOR_PATCH),3))
	@echo "New patch release version: v$(NEW_GIT_TAG)"
	@echo "gh release create v$(NEW_GIT_TAG) ./bin/*"

release-minor: clean build-all
	@echo "GIT_TAG: $(GIT_TAG)"
	$(eval MAJOR_MINOR_PATCH = $(call extract_version_parts,$(GIT_TAG:v%=%)))
	$(eval NEW_GIT_TAG = $(call increment_version,$(MAJOR_MINOR_PATCH),2))
	@echo "New minor release version: v$(NEW_GIT_TAG)"
	@echo "gh release create v$(NEW_GIT_TAG) ./bin/*"

release-major: clean build-all
	@echo "GIT_TAG: $(GIT_TAG)"
	$(eval MAJOR_MINOR_PATCH = $(call extract_version_parts,$(GIT_TAG:v%=%)))
	$(eval NEW_GIT_TAG = $(call increment_version,$(MAJOR_MINOR_PATCH),1))
	@echo "New major release version: v$(NEW_GIT_TAG)"
	@echo "gh release create v$(NEW_GIT_TAG) ./bin/*"

