#!/bin/bash

go build -o api.o api.go
chmod +x api.o
./api.o -config $1
