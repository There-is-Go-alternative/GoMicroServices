PROJ_DIR				:= $(realpath $(CURDIR))

# Compose
COMPOSE					:= docker-compose
COMPOSE_COMMON_FILE		:= docker-compose.yml
COMPOSE_TEST_FILE		:= docker-compose.test.yml

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

.PHONY: compose-build-test
compose-build-test:
	$(COMPOSE) -f $(COMPOSE_COMMON_FILE) -f $(COMPOSE_TEST_FILE) build

.PHONY: go-build
go-build: $(MICRO_SERVICES_DIRS)
	@for ms_dir in $^ ; do 																							\
  		echo "Building {$${ms_dir}} microservice ..." ; 															\
		cd $(PROJ_DIR)/$${ms_dir} && go generate ./... && go build -o $${ms_dir} . && rm $${ms_dir} || exit 1; 		\
  		echo "Microservice {$${ms_dir}} Built !\n" ; 																\
	done

.PHONY: lint
lint: $(MICRO_SERVICES_DIRS)
	@for ms_dir in $^ ; do 										\
  		echo "Running lint in {$${ms_dir}} ..." ; 				\
		cd $(PROJ_DIR)/$${ms_dir} &&  $(LINT_CMD) || exit 1; 	\
  		echo "Microservice {$${ms_dir}} Linted !\n" ; 			\
	done

.PHONY: account-migration
account-migration:
	cd account && go run github.com/prisma/prisma-client-go db push --schema infra/database/schema.prisma
	cd account && go run github.com/prisma/prisma-client-go generate --schema infra/database/schema.prisma

.PHONY: compose-clean
compose-clean:
	$(COMPOSE) rm -fsv

.PHONY: docker-build
docker-build:
	@for ms_dir in $^ ; do 																		\
		echo "Running lint in {$${ms_dir}} ..." ; 												\
		cd $(PROJ_DIR)/$${ms_dir} && docker build -t $${ms_dir}:latest -f build/Dockerfile .	\
		echo "Microservice {$${ms_dir}} Linted !\n" ; 											\
	done

	@echo "y" | docker system prune -a --volumes

.PHONY: docker-clean
docker-clean:
	@echo "y" | docker system prune -a --volumes
