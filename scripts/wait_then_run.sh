#!/bin/sh

# Checks a connection can be made to <host> <port> then runs <command>
# Makes up to <retries> attempts, with a 5 second delay between attempts
#
# Usage:
# ./wait_then_run.sh <host> <port> <retries>
#
# Example:
# ./wait_then_run.sh localhost 5432 30
#
# number of retries defaults to 10

set -e

HOST=$1
PORT=$2

if [ "$3" != "0" ]; then
    RETRIES=$3
else
    RETRIES=10
fi

while [ $RETRIES -ge 0 ]; do
    if nc -z $HOST $PORT; then
        RESULT=$?
    else
        RESULT=$?
    fi

    if [ $RESULT -eq 0 ]; then
        echo "$HOST:$PORT is ready!"
        sleep 5
        break
    elif [ $RETRIES -eq 0 ]; then
        echo "Exhausted attempts, exiting"
        break
    else
        RETRIES=$(($RETRIES-1))
        echo "Waiting for $HOST:$PORT, $RETRIES remaining attempts..."
        sleep 10
    fi
done
