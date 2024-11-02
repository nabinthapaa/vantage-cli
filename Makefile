.PHONY: install
install:
	go build main.go
	mv main /usr/bin/vantage
	chmod a+rx /usr/bin/vantage

.PHONY: dev
dev:
	air
