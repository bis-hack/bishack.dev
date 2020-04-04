GO ?= go

test: lint
	@echo '  -> running test' @$(GO) test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo
.PHONY: test

setup:
	@if [[ ! -x `which go` ]]; then echo '\n  Go is not installed!'; exit; fi;
	@if [[ ! -x `which docker` ]]; then echo '\n  Docker is not installed!'; exit; fi;
	@echo
	@echo ' Setting up development environment'
	@echo
	@echo ' [Go]'
	@echo ' -> downloading go modules...'
	@$(GO) mod download
	@echo ' -> downloading golangci-lint...'
	@cd ~ && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.24.0 > /dev/null && cd - > /dev/null
	@echo ' -> downloading hot-reload program called air...'
	@cd ~ && go get -u github.com/cosmtrek/air && cd - > /dev/null
	@echo
	@echo ' [Apex/Up]'
	@echo ' -> downloading apex/up...'
	@curl -sf https://up.apex.sh/install | sh > /dev/null
	@echo
	@echo ' [Docker]'
	@echo ' -> [dynamodb] setting up service...'
	@docker rm -f dynamo > /dev/null 2>&1
	@docker run --name dynamo -d -p 8000:8000 amazon/dynamodb-local -jar DynamoDBLocal.jar -inMemory -sharedDb > /dev/null
	@echo ' -> [dynamodb] rock en roll!'
	@echo
	@echo ' [IMPORTANT!]'
	@echo ' Create a .env file and ask @penzur for the content'
	@echo
	@echo ' Done ✔'
	@echo
.PHONY: setup


lint:
	@echo '  -> running golangci_lint'
	@golangci-lint run
	@echo
.PHONY: lint

dev:
	@air
.PHONY: dev


deploy: test up clean
	@echo "  -> done ✓"
.PHONY: deploy

deploy.prod: test up.prod clean
	@echo "  -> done ✓"
.PHONY: deploy

destroy: up.json
	@up stack delete
	@rm -rf up.json
	@echo "  ✓ done"
.PHONY: destroy


up: up.json
	@echo "  -> deploying"
	@up
.PHONY: up

up.prod: up.json.prod
	@echo "  -> deploying production"
	@up production
.PHONY: up

clean:
	@rm -rf up.json
	@rm -rf ./dist/
.PHONY: clean

# parse up template
up.json:
	@echo "  -> creating up.json from template file"
	@cat up.tmpl | sed "s/\$$COGNITO_CLIENT_ID/${COGNITO_CLIENT_ID}/g" \
		| sed "s/\$$COGNITO_CLIENT_SECRET/${COGNITO_CLIENT_SECRET}/g" \
		| sed "s/\$$COGNITO_POOL_ID/${COGNITO_POOL_ID}/g" \
		| sed "s/\$$GITHUB_CLIENT_SECRET/${GITHUB_CLIENT_SECRET}/g" \
		| sed "s/\$$GITHUB_CLIENT_ID/${GITHUB_CLIENT_ID}/g" \
		| sed "s/\$$SLACK_TOKEN/${SLACK_TOKEN}/g" \
		| sed "s/\$$SESSION_KEY/${SESSION_KEY}/g" \
		| sed "s/\$$CSRF_KEY/${CSRF_KEY}/g" \
		| sed "s|\$$GITHUB_CALLBACK|${GITHUB_CALLBACK}|g" \
		| sed "s/\$$DYNAMO_TABLE_POSTS/${DYNAMO_TABLE_POSTS}/g" \
		| sed "s/\$$DYNAMO_TABLE_LIKES/${DYNAMO_TABLE_LIKES}/g" \
		> up.json
# parse up template for prod
up.json.prod:
	@echo "  -> creating up.json from template file"
	@cat up.tmpl | sed "s/\$$COGNITO_CLIENT_ID/${COGNITO_CLIENT_ID_PROD}/g" \
		| sed "s/\$$COGNITO_CLIENT_SECRET/${COGNITO_CLIENT_SECRET_PROD}/g" \
		| sed "s/\$$COGNITO_POOL_ID/${COGNITO_POOL_ID_PROD}/g" \
		| sed "s/\$$GITHUB_CLIENT_SECRET/${GITHUB_CLIENT_SECRET_PROD}/g" \
		| sed "s/\$$GITHUB_CLIENT_ID/${GITHUB_CLIENT_ID_PROD}/g" \
		| sed "s/\$$SLACK_TOKEN/${SLACK_TOKEN}/g" \
		| sed "s/\$$SESSION_KEY/${SESSION_KEY}/g" \
		| sed "s/\$$CSRF_KEY/${CSRF_KEY}/g" \
		| sed "s|\$$GITHUB_CALLBACK|${GITHUB_CALLBACK_PROD}|g" \
		| sed "s/\$$DYNAMO_TABLE_POSTS/${DYNAMO_TABLE_POSTS_PROD}/g" \
		| sed "s/\$$DYNAMO_TABLE_LIKES/${DYNAMO_TABLE_LIKES_PROD}/g" \
		> up.json
