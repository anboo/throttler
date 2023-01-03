#!make
include .env
export $(shell sed 's/=.*//' .env)

docker-run:
	docker-compose up --build

gorun:
	go build -o ./throttler ./ && ./throttler