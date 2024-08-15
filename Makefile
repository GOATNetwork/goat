# The name of the binary to be built
BINARY_NAME := goatd

# The Go compiler command
GO := go

# Default target to run
all: build

# Build the project
build:
	# Compile the Go project with the specified ldflags for the binary
	# and place the binary in the ./bin directory.
	$(GO) build -ldflags "-s -w" -o ./bin/$(BINARY_NAME) ./cmd/goatd

# Clean up the built files
clean:
	# Remove the ./bin directory to clean up the build artifacts
	rm -rf ./bin

# Phony target to avoid conflicts with files named 'all', 'build', or 'clean'
.PHONY: all build clean