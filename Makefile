.PHONY: generate
generate:
	go generate -x ./generate/generator/main.go

.PHONY: build
build: generate
	go build -o rhino .

.PHONY: install
install: generate
	go build -o rhino . 
	mv rhino /usr/local/bin

.PHONY: clean
clean:
	rm /usr/local/bin/rhino
