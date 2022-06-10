#!/bin/bash

go build -o web-app ./cmd/web/*.go
./web-app -production=false -cache=false -dbhost=localhost -dbport=5432 -dbname=bookings -dbuser=postgres -dbpass=postgres -dbssl=