#!/usr/bin/env bash
workersnr='16'

echo "Making sure previous run is killed"

./killall.sh

# start sink. Currently logs to "log-of-randori"
echo "starting sink"
./randorisink-linux-amd64 &
sleep 3

# spin up telnet workers
echo "starting telnet"
for i in {1..${workersnr}}; do ./toritelnet-linux-amd64 & done

# spin up ssh workers
echo "starting ssh"
for i in {1..${workersnr}}; do ./torissh-linux-amd64 & done
sleep 3

# start tailing the PAM randori log and feeding the fan
echo "starting fan"
tail -f -n 0 /var/log/randori.log | ./randorifan-linux-amd64 &

echo "everything started"
