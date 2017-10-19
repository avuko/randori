#!/usr/bin/env bash

fgrep -a TORI /var/log/syslog | fgrep -a -v 'result=ERROR' >> randori-tori.log
fgrep -aB1 'client software version' /var/log/auth.log | fgrep -a -A1 'Connection from ' | cut -d ':' -f 5,6 | sed 's/Connection from /|/g;s/ port /|/g;s/client software version /|/g' | awk -F'|' '{print $2}' | tr '\n' '|' | sed 's/||/\n/g;s/|$/\n/'  |sort -u >> randori-sshclients.log
