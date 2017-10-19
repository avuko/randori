#!/usr/bin/env bash
# workersnr = 16

echo "Making sure previous run is killed"

./killall.sh

# start sink. Currently logs to "log-of-randori"
#echo "starting sink"
./randorisink &
sleep 3

# spin up ftp workers
# echo "starting ftp"
# for i in {1..16}; do ./toriftp & done
# sleep 3

# spin up telnet workers
#echo "starting telnet"
for i in {1..16}; do ./toritelnet & done

# spin up ssh workers
#echo "starting ssh"
for i in {1..16}; do ./torissh & done
sleep 3


# start tailing the PAM randori log and feeding the fan
# echo "starting fan"
tail -F -n 0 /var/log/randori.log | ./randorifan &

# echo "everything started"

