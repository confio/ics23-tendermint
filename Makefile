.PHONY: test

# make sure we turn on go modules
export GO111MODULE := on

# build:
# 	go build -mod=readonly ./cmd/testgen-iavl

test:
	go test -mod=readonly .

# testgen:
# 	go run -mod=readonly ./cmd/testgen-iavl