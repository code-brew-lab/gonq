PROJECT_NAME ?= gonq

GOOS ?= $(shell uname)
GOARCH ?= $(shell arch)
CGO_ENABLED ?= 0

BIN_NAME = "$(PROJECT_NAME)"

SRC_DIR = ./cmd
BIN_DIR = ./bin

format:
	@echo "Formatting..."
	@gofmt -s -w .
	@echo "Done!"

install:
	@echo "Installing dependencies..."
	@go mod download && go mod verify
	@echo "Done!"

tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Done!"

clean-cache:
	@echo "Cleaning build and package cache..."
	@go clean -modcache
	@go clean -cache
	@echo "Done!"

build:
	@echo "Building api..."
	@echo "GOOS         = $(GOOS)"
	@echo "GOARCH       = $(GOARCH)"
	@echo "CGO_ENABLED  = $(CGO_ENABLED)"
	@go build -o $(BIN_DIR)/$(BIN_NAME) -v -x $(SRC_DIR)/${PROJECT_NAME}/.
	@echo "Built!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf $(BIN_DIR)/*
	@echo "Cleaned!"

help:
	@echo "-------------------------------------------------------------------------"
	@echo "                           ${PROJECT_NAME} makefile"
	@echo "-------------------------------------------------------------------------"
	@echo ""
	@echo "Make available targets:"
	@echo "  format          : Format all project tree."
	@echo ""
	@echo "  install         : Installs the projects dependencies."
	@echo "  tidy            : Tidy projects dependencies."
	@echo "  clean-cache     : Cleans project cache."
	@echo ""
	@echo "  build           : Builds the binary for the specified GOOS and GOARCH."
	@echo ""
	@echo "  clean           : Cleans up generated files related to the binary."
	@echo ""
	@echo "Usage:"
	@echo "  make <target> [VARIABLE=value]"
	@echo ""
	@echo "Variables available for override (with defaults):"
	@echo "PROJECT_NAME         = $(PROJECT_NAME)"
	@echo "GOOS                 = $(GOOS)"
	@echo "GOARCH               = $(GOARCH)"
	@echo "CGO_ENABLED          = $(CGO_ENABLED)"
	@echo "-------------------------------------------------------------------------"
