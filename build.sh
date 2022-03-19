#!/bin/bash

cd collector && go build -o bin/lmc ./cmd
cd ../client/metrics-dashboard && npm run build