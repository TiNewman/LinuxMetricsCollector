#!/bin/bash

# if running the script from another directory, navigates to the project root dir
cd "$(dirname -- "${BASH_SOURCE[0]}")"

# build the collector component with go build
cd collector && go build -o bin/lmc ./cmd

# build the UI component with npm
cd ../client/metrics-dashboard && npm run build

# build the db docker container
cd ../../ && sudo docker-compose build