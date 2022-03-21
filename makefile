.PHONY: all ui collector database

all: collector ui database

collector:
	cd collector && go build -o bin/lmc ./cmd

ui:
	cd client/metrics-dashboard && npm run build

database:
	sudo docker-compose build

