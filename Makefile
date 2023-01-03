#!make
include .env
export $(shell sed 's/=.*//' .env)

run:
	go build -o ./throttler ./ && ./throttler