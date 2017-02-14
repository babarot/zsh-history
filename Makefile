build: dep
	go build -o zhist cmd/zhist/main.go

dep:
	go get -d

install: build
	sudo install -m 0755 zhist /usr/local/bin
