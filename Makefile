NAMESPACE=adduc
NAME=woodpecker
BINARY=terraform-provider-$(NAME)
VERSION=0.0.1-dev
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -i | sed 's/x86_64/amd64/')
DIR=terraform.local/${NAMESPACE}/${NAME}/${VERSION}/$(OS)_$(ARCH)

build:
	go build -o artifacts/$(BINARY)

# install the plugin locally
# @todo deprecate in favor of using provider_installation.dev_overrides
install: build
	mkdir -p ~/.terraform.d/plugins/$(DIR)
	cp artifacts/$(BINARY) ~/.terraform.d/plugins/$(DIR)

# generate docs
doc:
	tfplugindocs

# reset test environment
reset:
	.ci/reset.sh

# run unit tests
test:
	echo "@todo"

# run acceptance tests
testacc:
	bash -c 'set -a; source .ci/.env; env TF_ACC=1 go test -v ./...'