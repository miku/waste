waste: cmd/waste/main.go
	go get ./...
	go build -o $@ $?

clean:
	rm -f waste
