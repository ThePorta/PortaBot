all: bot.bin

clean:
	rm -rf bin/*

%.bin: ./cmd
	go build -o bin/$* ./cmd/$*