SRC:=$(shell find . -name '*.go')

krew-index-tracker: $(SRC)
	go build -ldflags "-s -w" -o $@ ./app/krew-index-tracker

lint: $(SRC)
	hack/run-lint.sh

test: $(SRC)
	go test
