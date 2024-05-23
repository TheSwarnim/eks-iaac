.PHONY: test up

# Set the default stack to "dev"
STACK ?= dev

# Run tests
test:
	cd src && go test -cover ./...

# Run pulumi up
up: test
	cd iaac && pulumi up --stack $(STACK)

# Run pulumi refresh
refresh:
	cd iaac && pulumi refresh --stack $(STACK)

# Run pulumi destroy
destroy:
	cd iaac && pulumi destroy --stack $(STACK)