build:
	export GO15VENDOREXPERIMENT=1
	go build
	go test -cover

test:
	export GO15VENDOREXPERIMENT=1
	go test

coverage:
	export GO15VENDOREXPERIMENT=1
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
