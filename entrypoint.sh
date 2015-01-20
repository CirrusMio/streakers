#!/bin/bash

go get github.com/nitrous-io/goop
goop install

goop exec go run server.go
