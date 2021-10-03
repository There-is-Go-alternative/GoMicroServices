PROJ_DIR				:= $(realpath $(CURDIR))

# Create the list of directories for micro services, separated by a space (I.E: 'account auth messages')
MICRO_SERVICES_DIRS		= account ads

# Create binary dir path
BIN_DIR_NAME			= .bin
BIN_DIR_PATH			= $(PROJ_DIR)/$(BIN_DIR_NAME)

# Create lint variables
LINT_BIN_NAME			= golangci-lint
LINT_BIN_PATH			= $(BIN_DIR_PATH)/$(LINT_BIN_NAME)
LINT_CMD				= $(LINT_BIN_PATH) run --timeout=1m

.PHONY: lint
lint: $(MICRO_SERVICES_DIRS)
	@for ms_dir in $^ ; do \
  		echo "Running lint in $${ms_dir} ..." ; 			\
		cd $(PROJ_DIR)/$${ms_dir} && $(LINT_CMD) || exit 1; \
  		echo "Microservice $${ms_dir} Linted !\n" ; 		\
	done
