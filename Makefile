GO ?= go

test:
	@echo '  -> running test'
	@$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo
.PHONY: test

dev: up.json
	@up start
.PHONY: dev


deploy: test up clean
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


clean:
	@rm -rf up.json
	@rm -rf ./dist/
.PHONY: clean

# parse up template
up.json:
	@echo "  -> creating up.json from template file"
	@cat up.tmpl | sed "s/\$$COGNITO_CLIENT_ID/${COGNITO_CLIENT_ID}/g" \
		| sed "s/\$$COGNITO_CLIENT_SECRET/${COGNITO_CLIENT_SECRET}/g" \
		| sed "s/\$$GITHUB_CLIENT_SECRET/${GITHUB_CLIENT_SECRET}/g" \
		| sed "s/\$$GITHUB_CLIENT_ID/${GITHUB_CLIENT_ID}/g" \
		| sed "s/\$$SLACK_TOKEN/${SLACK_TOKEN}/g" \
		| sed "s/\$$SESSION_KEY/${SESSION_KEY}/g" \
		| sed "s/\$$CSRF_KEY/${CSRF_KEY}/g" \
		> up.json
