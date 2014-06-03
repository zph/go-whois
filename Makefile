lint:
	./bin/go-lint

hooks:
	cp -f hooks/* .git/hooks/

build: clean
	go build -o dist/whois server.go; \
		GOARCH=amd64 GOOS=linux go build -o dist/whois.linux_amd64 server.go; \
		go build -gcflags '-N' -o dist/whois.debug server.go;

clean:
	rm -f dist/*

email:
	go build -o dist/email extract_email.go
