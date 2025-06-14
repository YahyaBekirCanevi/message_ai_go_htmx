#!/bin/bash
# export CGO_ENABLED=1
go run main.go migration.go . "$@" # "$@" passes any arguments to the go run command