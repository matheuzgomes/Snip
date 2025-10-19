build:
	go build -o snip main.go
test:
	go test -v ./internal/test/...

bench:
	go test -run='^$$' -bench=. ./internal/test/...