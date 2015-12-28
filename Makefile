
all:
	go build -o image-mining

clean:
	go clean
	rm -rf image-mining
