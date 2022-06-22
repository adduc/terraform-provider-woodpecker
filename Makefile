HOSTNAME?=github.com
NAMESPACE=adduc
NAME=woodpecker
BINARY=terraform-provider-$(NAME)
VERSION=0.0.1-dev
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -i | sed 's/x86_64/amd64/')
DIR=${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/$(OS)_$(ARCH)

build:
	go build -o bin/$(BINARY)

install: build
	mkdir -p ~/.terraform.d/plugins/$(DIR)
	cp bin/$(BINARY) ~/.terraform.d/plugins/$(DIR)

test: install
	cd demo && rm -r .terraform.lock.hcl .terraform && terraform init && terraform plan
