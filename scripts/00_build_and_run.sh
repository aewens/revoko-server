#!/bin/bash

# $1 = config file
# $2 = log file
# $3 = pid file

# Kill previous process
if [ -f $3 ]; then
    kill -9 $(cat $3)
fi

# Get scope of script
DIR=`dirname "$0"`
SRC="$DIR/.."

# Build the project to an executable
go build -o $SRC/api.o $SRC/api.go
chmod +x $SRC/api.o

# Run project using config file and logging output to file
nohup $SRC/api.o -config $1 2>&1 >> $2 &>/dev/null &

# Write PID to pid file
echo $! > $3

# Remove config file
rm $1
