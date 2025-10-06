build:
	go build -o snip main.go

release:
	goreleaser release --clean