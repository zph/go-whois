build: clean
	go build -o dist/whois_server server.go; \
		GOARCH=amd64 GOOS=linux go build -o dist/whois_server.linux_amd64 server.go; \
		go build -gcflags '-N' -o dist/whois_server.debug server.go;

lint:
	./bin/go-lint

hooks:
	cp -f hooks/* .git/hooks/

clean:
	rm -f dist/*

email:
	go build -o dist/email extract_email.go


