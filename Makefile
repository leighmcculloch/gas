build:
	go build

release:
	go run goreleaser --rm-dist
