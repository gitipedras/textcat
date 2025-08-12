APP_NAME := textcat
BUILD_DIR := build

OS_LIST := windows linux darwin
ARCH_LIST := amd64 arm64

# Extension for executables on Windows
ifeq ($(OS),windows)
	EXT := .exe
else
	EXT :=
endif

.PHONY: all clean

all: clean build archive

clean:
	rm -rf $(BUILD_DIR)

build:
	@mkdir -p $(BUILD_DIR)
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			echo "Building for $$os/$$arch..."; \
			EXT=""; \
			if [ "$$os" = "windows" ]; then EXT=".exe"; fi; \
			out_dir=$(BUILD_DIR)/$$os-$$arch; \
			mkdir -p $$out_dir; \
			GOOS=$$os GOARCH=$$arch go build -o $$out_dir/$(APP_NAME)$${EXT} . || exit 1; \
		done; \
	done

archive:
	@zip -r build-latest.zip build
