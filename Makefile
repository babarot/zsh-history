build: dep
	go build -o zhist cmd/zhist/main.go

dep:
	go get -d
