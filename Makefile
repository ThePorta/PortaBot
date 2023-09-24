all: bot.bin server.bin signalBuyBot.bin

clean:
	rm -rf bin/*

%.bin: ./cmd
	go build -o bin/$* ./cmd/$*