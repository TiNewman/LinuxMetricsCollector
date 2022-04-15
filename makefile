.PHONY: all ui collector database test

all: collector ui database

collector:
	cd collector && go build -o bin/lmc ./cmd

ui:
	cd client/metrics-dashboard && npm run build

database:
	sudo docker-compose build

test: 
	./ttar -C collector/pkg/process/ -x -f process.ttar && cd collector/pkg/process && go test -v
	./ttar -C collector/pkg/cpu/ -x -f cpu.ttar && cd collector/pkg/cpu && go test -v
	cd collector/pkg/http/websocket && go test -v -p 1