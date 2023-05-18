.PHONY: test
test:
	@bash ./scripts/test.sh

.PHONY: deps-init
deps-init:
	echo "INFO: creating dependencies..."
	@docker-compose up -d --build
	@bash ./scripts/init-dep.sh ./scripts/init-vault ./mock-schema/vault
	
.PHONY: deps-tear
deps-tear:
	@docker-compose down --volumes --remove-orphans

.PHONY: run-local
run-local:
	@go build ./cmd/user-manager/ && ./user-manager