#!make
include .env
export $(shell sed 's/=.*//' .env)

run:
	go run ./*.go