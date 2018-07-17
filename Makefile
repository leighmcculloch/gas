release:
	rm -fr dist
	goreleaser

setup:
	go get github.com/goreleaser/goreleaser
