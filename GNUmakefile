CHANGE_HOST?=localhost
TEST?=$$(go list ./... |grep -v 'vendor')
SWEEP_DIR?=./change
SWEEP?=$(CHANGE_ADDRESS)
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: testacc bin

bin:
	go install
	scripts/install_plugin.sh

tools:
	@echo "==> Installing external tools..."
	GO111MODULE=on go get -u golang.org/x/lint/golint
	GO111MODULE=on go get -u gotest.tools/gotestsum
	GO111MODULE=on go get -u github.com/gordonklaus/ineffassign
	GO111MODULE=on go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go get -u github.com/katbyte/terrafmt
	GO111MODULE=on go get -u github.com/bflad/tfproviderdocs

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint: tools fmtcheck vet depscheck docs
	@echo "==> Checking source code against linters..."
	golint -set_exit_status $$(find . -type d | grep -v vendor)
	ineffassign .

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m

test: lint
	go test -v $(TEST)

testacc: lint
	@echo "==> Running acceptance tests..."
	TF_ACC=1 CHANGE_HOST=$(CHANGE_HOST) go test $(TEST)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/)

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

docs: docs-lint docscheck

docs-lint:
	@echo "==> Checking docs against linters..."
	@misspell -error -source=text docs/ || (echo; \
		echo "Unexpected misspelling found in docs files."; \
		echo "To automatically fix the misspelling, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli docs/ || (echo; \
		echo "Unexpected issues found in docs Markdown files."; \
		echo "To apply any automatic fixes, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@terrafmt diff ./docs --check --pattern '*.md' --quiet || (echo; \
		echo "Unexpected differences in docs HCL formatting."; \
		echo "To see the full differences, run: terrafmt diff ./docs --pattern '*.md'"; \
		echo "To automatically fix the formatting, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)

docs-lint-fix:
	@echo "==> Applying automatic docs linter fixes..."
	@misspell -w -source=text docs/
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/
	@terrafmt fmt ./docs --pattern '*.md'

docscheck:
	@tfproviderdocs check \
		-allowed-resource-subcategories-file docs/allowed-subcategories.txt \
		-require-resource-subcategory
	@misspell -error -source text CHANGELOG.md

.PHONY: build tools fmt fmtcheck lint sweep test testacc vet depscheck docs docs-lint docs-lint-fix docscheck change-image
