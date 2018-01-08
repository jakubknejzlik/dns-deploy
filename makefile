install:
	go get ./...
	go build -o dns-deploy

deploy-local:
	make install
	mv dns-deploy /usr/local/bin/
	chmod +x /usr/local/bin/dns-deploy
