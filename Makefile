
BASE_DIR ?= $(shell cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

USER_MANAGER_ENV_PATH ?="$(BASE_DIR)/user-manger-service/infra/docker/.env"

SKAFFOLD_PATH ?= "$(BASE_DIR)/infra/skaffold/skaffold.yaml"

SKAFFOLD_INFRA_PATH ?= "$(BASE_DIR)/infra/skaffold/skaffold.infra.yaml"

.PHONY: docker_dev

docker_dev:
	@docker-compose -f docker-compose.yml up

.PHONY: docker_start_rabbitmq

docker_start_rabbitmq:
	@docker-compose start rabbitmq

.PHONY: restart

restart:
	@docker-compose restart

.PHONY: down

down:
	@docker-compose down --volumes

.PHONY: stop

stop:
	@docker-compose stop

.PHONY: docker_infra

docker_infra:
	@docker-compose start

.PHONY: k8s_run

k8s_run:	
	@skaffold -f $(SKAFFOLD_PATH) run

.PHONY: k8s_create_infra

k8s_create_infra:
	@skaffold -f $(SKAFFOLD_INFRA_PATH) run

.PHONY: k8s_setup

k8s_setup:
	./infra/script/create-secrets.sh
	./infra/script/setup.sh

.PHONY: k8s_del_app

k8s_del_app:
	@skaffold -f $(SKAFFOLD_PATH) delete

.PHONY: k8s_del_infra

k8s_del_infra:
	@skaffold -f $(SKAFFOLD_INFRA_PATH) delete

.PHONY: k8s_del_secret

k8s_del_secret:
	./infra/script/delete-secrets.sh

.PHONY: test

test:	
	go list -f '{{.Dir}}/...' -m | xargs -L1 go test -timeout 120s -race -coverprofile=coverage.out -tags=integration
