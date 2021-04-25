MAKEFLAGS += --silent

.PHONY: run

no_targets__:

run/dev:
	CompileDaemon -build "go build -v -o ./bin ./main.go" \
	              -command "./bin/main" -exclude-dir=.git \
	              -exclude-dir=bin -log-prefix=false

run:
	make build
	./bin/main

build:
	echo "Compiling source code.."
	go build -v -o ./bin ./main.go
	echo "Done!"

fmt:
	echo "Formatting source code.."
	go fmt ./...

.PHONY: test-with-exit-code
test-with-exit-code:
	ENV=test go test -v -race -timeout 30s ./...

.PHONY: test
test:
	touch coverage.out | \
	ENV=test go test -v -coverprofile=coverage.out -covermode=count \
	                 -json -cover -timeout 30s ./... | tparse -all

.PHONY: test/cover
test/cover:
	go tool cover -html=coverage.out -o coverage.html | \
	open ./coverage.html

dev/migrate:
	tern migrate -c config/development.conf -m store/migrations
test/migrate:
	tern migrate -c config/test.conf -m store/migrations

dev/rollback:
	tern migrate -c config/development.conf -m store/migrations -d -1
test/rollback:
	tern migrate -c config/test.conf -m store/migrations -d -1

dev/drop:
	tern migrate -c config/development.conf -m store/migrations -d 0
test/drop:
	tern migrate -c config/test.conf -m store/migrations -d 0

.DEFAULT_GOAL := run/dev

list:
		sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | sort"
