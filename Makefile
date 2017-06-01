build: gopath gobin install

install: 
	go install

gopath: 
	export GOPATH="$(PWD)"

gobin: 
	export GOBIN="$(PWD)/bin"