build:
	go build

release:
	go run github.com/goreleaser/goreleaser --rm-dist
	./dist/gas_linux_amd64/gas -version
