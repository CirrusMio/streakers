#!/bin/bash

cd /var/www

go get github.com/nitrous-io/goop
goop install

goop exec go run server.go
