#!/usr/bin/env bash
# camq    : Cheap Ass Message Queue.
# author  : iamavuko@gmail.com
# date    : 2017-06-27
# version : 0.1

# exit on error
set -e
# exit on unset variable
set -u

# Make sure bash is compatible
if [ ${BASH_VERSION:0:1} != 4 ];then
echo "Not bash, or not bash version 4"
exit 1
fi

### global variables

# all log files and all of the pipes used for randori.
LOGFILE='randori.log'
WEIRDFILE='weird.log'

# create daemons that we will track and create fifo's.
declare -A DEMONS
DEMONS=([HTTPD]=apache [TELNETD]=login [SSHD]=sshd)

	for DEMON in ${DEMONS[@]};do
	if ! [[ -p $DEMON ]];then
	mkfifo $DEMON
	fi
	done

# store old Internal Field (word) Separator for later.
OLDIFS="${IFS}"

# Set new separator. This is also set in pam_randori.c
# if someone is messing with us, they are at least not
# trying to log into our system with valid credentials

IFS=$(printf "\u2002")

# store everything in an array for quick & easy access
declare -a LINEARRAY

# Lets make sure we are not dealing with any weirdness.
while read line ; do LINEARRAY=(${line});\
if [ ${#LINEARRAY[@]} == 4 ]; then

# All is fine, lets continue with our bashy debauchery.
echo ${LINEARRAY[0]} && echo "${LINEARRAY[@]:1}" > ${LINEARRAY[0]}&
else
# Somebody is trying to pull a fast one. Go away please
echo "${LINEARRAY[@]}" >> "${WEIRDFILE}"
fi ;done <"${LOGFILE}"

# reset IFS
IFS="${OLDIFS}"
