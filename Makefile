PROJ_DIR				:= $(realpath $(CURDIR))

# Commands
COMPOSE					:= docker-compose

# Create the list of directories for micro services, separated by a space (I.E: 'account auth messages')
MICRO_SERVICES_DIRS		= account ads

# Create binary dir path
BIN_DIR_NAME			= .bin
BIN_DIR_PATH			= $(PROJ_DIR)/$(BIN_DIR_NAME)

# Create lint variables
LINT_BIN_NAME			= golangci-lint
LINT_BIN_PATH			= $(BIN_DIR_PATH)/$(LINT_BIN_NAME)
LINT_CMD				= $(LINT_BIN_PATH) run --timeout=1m

.PHONY: compose-build
compose-build:
	$(COMPOSE) build --parallel

.PHONY: compose-run
compose-run: compose-build
	$(COMPOSE) up

.PHONY: compose-run-bg
compose-run-bg: compose-build
	$(COMPOSE) up -d

.PHONY: compose-run-db
compose-run-db:
	$(COMPOSE) up postgres_db

.PHONY: go-build
go-build: $(MICRO_SERVICES_DIRS)
	@for ms_dir in $^ ; do 																		\
  		echo "Building {$${ms_dir}} microservice ..." ; 										\
		cd $(PROJ_DIR)/$${ms_dir} && go build -o $${ms_dir} . && rm $${ms_dir} || exit 1; 		\
  		echo "Microservice {$${ms_dir}} Built !\n" ; 											\
	done

.PHONY: lint
lint: $(MICRO_SERVICES_DIRS)
	@for ms_dir in $^ ; do 									\
  		echo "Running lint in {$${ms_dir}} ..." ; 			\
		cd $(PROJ_DIR)/$${ms_dir} && $(LINT_CMD) || exit 1; \
  		echo "Microservice {$${ms_dir}} Linted !\n" ; 		\
	done

.PHONY: compose-clean
compose-clean:
	$(COMPOSE) rm -fsv