#!make
include .env
export $(shell sed 's/=.*//' .env)

docker-run:
	docker-compose up --build

go-run:
	go build -o ./throttler ./ && ./throttler