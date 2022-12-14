NAMESPACE=adduc
NAME=woodpecker
BINARY=terraform-provider-$(NAME)
VERSION=0.0.1-dev
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -i | sed 's/x86_64/amd64/')
DIR=terraform.local/${NAMESPACE}/${NAME}/${VERSION}/$(OS)_$(ARCH)

build:
	go build -o artifacts/$(BINARY)

install: build
	mkdir -p ~/.terraform.d/plugins/$(DIR)
	cp artifacts/$(BINARY) ~/.terraform.d/plugins/$(DIR)

doc:
	tfplugindocs

test-reset: install
	cd demo && rm -rf terraform.tfstate \

test-plan: install
	cd demo && terraform plan # apply -auto-approve

test-apply: install
	cd demo && terraform apply -auto-approve

test-import: install
	cd demo && terraform import woodpecker_repository.repository jlong/repo-3