APPNAME = marathon-logger
VERSION=0.0.7
TESTFLAGS=-v -cover -covermode=atomic -bench=.
TEST_COVERAGE_THRESHOLD=8.0

build:
	go build -tags netgo -ldflags "-w" -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X main.APP_VERSION=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X main.APP_VERSION=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

build-all: build-mac build-linux

clean:
	rm -f ${APPNAME}
	rm -f ${APPNAME}-linux-amd64
	rm -f ${APPNAME}-darwin-amd64

all: setup
	build
	install

setup:
	go get -u github.com/wadey/gocovmerge
	glide install

test-only:
	go test ${TESTFLAGS} -coverprofile=${name}.txt github.com/ashwanthkumar/marathon-logger/${name}

test:
	go test ${TESTFLAGS} -coverprofile=main.txt github.com/ashwanthkumar/marathon-logger/

ci: test-ci

test-ci: test
	${GOPATH}/bin/gocovmerge *.txt > coverage.txt
	@go tool cover -html=coverage.txt -o coverage.html
	@go tool cover -func=coverage.txt | grep "total:" | awk '{print $$3}' | sed -e 's/%//' > cov_total.out
	@bash -c 'COVERAGE=$$(cat cov_total.out);	\
	echo "Current Coverage % is $$COVERAGE, expected is ${TEST_COVERAGE_THRESHOLD}.";	\
	exit $$(echo $$COVERAGE"<${TEST_COVERAGE_THRESHOLD}" | bc -l)'
