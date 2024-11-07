.PHONY: install
install:
	go build main.go
	mv main /usr/bin/vantage
	chmod a+rx /usr/bin/vantage

.PHONY: run
run:
	go build -o ./tmp/vantage main.go
	sudo ./tmp/vantage table

.PHONY: dev
dev:
	air
