#!/usr/bin/env bash

echo "Making sure previous run is killed"

./killall.sh

# start sink. Currently logs to "log-of-randori"
echo "starting sink"
./randorisink-linux-amd64 &
sleep 3

# spin up 4 telnet workers
echo "starting telnet"
for i in {1..4}; do ./toritelnet-linux-amd64 & done

# spin up 4 ssh workers
echo "starting ssh"
for i in {1..4}; do ./torissh-linux-amd64 & done
sleep 3

# start tailing the PAM randori log and feeding the fan
echo "starting fan"
tail -f -n 0 /var/log/randori.log | ./randorifan-linux-amd64 &

echo "everything started"
