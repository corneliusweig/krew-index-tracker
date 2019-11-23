SRC:=$(shell find . -name '*.go')

krew-index-tracker: $(SRC)
	go build -ldflags "-s -w" -o $@ ./app/krew-index-tracker

krew-index-tracker-http: $(SRC)
	go build -ldflags "-s -w" -o $@ ./app/http

lint: $(SRC)
	hack/run-lint.sh

test: $(SRC)
	go test
