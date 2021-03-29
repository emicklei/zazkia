GO=go
DOCKER=docker
TARGET=zazkia
K6=$(DOCKER) run -i loadimpact/k6 run

.PHONY: generate build test clean

all: build

generate: bindata.go

bindata.go:
	$(GO) generate
	
build: $(TARGET)

$(TARGET): bindata.go *.go
	$(GO) build 

test: test1 test2 test3 test4

test1:
	@./basicweb -port 8080 -dir ./tests & echo $$! > /tmp/basicweb.pid
	@sed 's/__IP__/'$$(hostname -i)'/g' tests/test1.js | $(K6) -
	@cat /tmp/basicweb.pid | xargs kill
	@rm -f /tmp/basicweb.pid

test2:
	@./basicweb -port 8080 -dir ./tests > /dev/null 2>&1 & echo $$! > /tmp/basicweb.pid
	@./zazkia -c -p 9191 -f ./tests/test2.json > /dev/null 2>&1 & echo $$! > /tmp/zazkia.pid
	@sed 's/__IP__/'$$(hostname -i)'/g' tests/test2.js | $(K6) -
	@cat /tmp/basicweb.pid | xargs kill
	@cat /tmp/zazkia.pid | xargs kill
	@rm -f /tmp/basicweb.pid /tmp/zazkia.pid

test3:
	@./basicweb -port 8080 -dir ./tests > /dev/null 2>&1 & echo $$! > /tmp/basicweb.pid
	@./zazkia -c -p 9191 -f ./tests/test3.json > /dev/null 2>&1 & echo $$! > /tmp/zazkia.pid
	@sed 's/__IP__/'$$(hostname -i)'/g' tests/test3.js | $(K6) -
	@cat /tmp/basicweb.pid | xargs kill
	@cat /tmp/zazkia.pid | xargs kill
	@rm -f /tmp/basicweb.pid /tmp/zazkia.pid

test4:
	@./basicweb -port 8080 -dir ./tests > /dev/null 2>&1 & echo $$! > /tmp/basicweb.pid
	@./zazkia -c -p 9191 -f ./tests/test4.json > /dev/null 2>&1 & echo $$! > /tmp/zazkia.pid
	@sed 's/__IP__/'$$(hostname -i)'/g' tests/test4.js | $(K6) -
	@cat /tmp/basicweb.pid | xargs kill
	@cat /tmp/zazkia.pid | xargs kill
	@rm -f /tmp/basicweb.pid /tmp/zazkia.pid

clean:
	$(DOCKER) system prune -f -a
