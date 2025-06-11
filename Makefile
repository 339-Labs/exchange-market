GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

TM_ABI_ARTIFACT := /Users/kit/code/go/com.web3/src/339/exchange-market/abis/uniswapv3/UniswapV3Pool.json
ABIGEN := /Users/kit/code/go/com.web3/bin/abigen

exchange-market:
	env GO111MODULE=on go build -o exchange-market -v $(LDFLAGS) ./cmd

clean:
	rm exchange-market

test:
	go test -v ./...

lint:
	golangci-lint run ./...

bindings:
	cat $(TM_ABI_ARTIFACT) \
		| $(ABIGEN) --pkg bindings \
		--abi - \
		--out bindings/uniswapv3/UniswapV3Pool.go \
		--type UniswapV3Pool \

.PHONY: \
	exchange-market \
	clean \
	test \
	bindings \
	lint