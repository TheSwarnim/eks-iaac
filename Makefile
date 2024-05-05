.PHONY: test up

# Set the default stack to "dev"
STACK ?= dev

# Run tests
test:
	cd src && go test -cover ./...

# Run pulumi up
up: test
	cd iaac && pulumi up --stack $(STACK)