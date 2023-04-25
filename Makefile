# executables
############
ASDF_EXEC						:= $(shell command -v asdf 2>/dev/null)
GO_EXEC					:= $(shell command -v go 2>/dev/null)
INSTALL_CMD					:= $(shell command -v brew 2>/dev/null)
PRE_COMMIT_EXEC			:= $(shell command -v pre-commit 2>/dev/null)
GOLINT_EXEC		:= $(shell command -v golint 2>/dev/null)

ifndef VERBOSE
.SILENT:
endif

# Main targets
##############
.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: ## Check that required binaries are installed
ifndef INSTALL_CMD
	$(error Check that brew is installed and available in your PATH)
else
install: .pre-commit .asdf .go .golint
endif

.PHONY: update
update: ## Update external dependencies
ifndef PRE_COMMIT_EXEC
	$(error Unable to update as required binaries are missing. Please run `make install`...)
else
	$(info [*] checking for external pre-commit hook updates...)
	$(PRE_COMMIT_EXEC) autoupdate
	$(info [*] running sanity check of all hooks...)
	$(PRE_COMMIT_EXEC) run --all-files
endif

.PHONY: run-all
run-all: ## Runs all pre-commit hooks
ifndef PRE_COMMIT_EXEC
	$(error Unable to update as required binaries are missing. Please run `make install`...)
else
	$(PRE_COMMIT_EXEC) run --all-files
endif

# Build project
################
.PHONY: .build
build:
	go build -o autommit-dev ./cmd/autommit/main.go

# Helpers
#########
.PHONY: .asdf
.asdf:
ifndef ASDF_EXEC
	$(info [*] installing asdf...)
	$(INSTALL_CMD) install asdf &>/dev/null || echo "$(CCRED) unable to install asdf via $(INSTALL_CMD). Please install manually..."
else
	$(info [*] asdf already installed. Nothing to do...)
	exit 0
endif

.PHONY: .go
.go:
ifndef GO_EXEC
	$(info [*] installing golang...)
	$(INSTALL_CMD) install go &>/dev/null || echo "$(CCRED) unable to install go via $(INSTALL_CMD). Please install manually..."
else
	$(info [*] go already installed. Nothing to do...)
	exit 0
endif

.PHONY: .golint
.golint:
ifndef GOLINT_EXEC
	$(info [*] installing golint...)
	$(INSTALL_CMD) install golint &>/dev/null || echo "$(CCRED) unable to install golint via $(INSTALL_CMD). Please install manually..."
else
	$(info [*] golint already installed. Nothing to do...)
	exit 0
endif
