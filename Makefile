PACKAGE_NAME := shier
PROJECT_DIR := $(PWD)

install-dep:
	go get -u github.com/shuLhan/go-bindata/...
	dep ensure -v

test:
	# this project still have no tests, coming soon :p
	go test -v -cover -race ./...

generate-migration:
	# rm assets/sql/tags # remove generated file from ctags
	go-bindata -pkg assets -o assets/bindata.go assets/sql/...

build:
	rm -f out/shier # always remove existing binary
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/shierdb $(PACKAGE_NAME)