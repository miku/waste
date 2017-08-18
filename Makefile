waste: cmd/waste/main.go
	go build -o $@ $?

clean:
	rm -f waste
