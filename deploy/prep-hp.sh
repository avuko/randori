#!/usr/bin/env bash

# ONLY CHANGE THIS IF YOU HAVE PASSWORDLESS root ACCESS
AREYOUSURE='yes'

# need to be root for this
if [ $(id -u) == 0 ];then
	echo "root, continue"
else
	echo "not root, exiting"
 	exit 1
fi


# change the root password to something impossibly hard to brute force
if [ $AREYOUSURE != yes ];then
echo "Did not get a clear \"yes\", exiting."
exit
else
	if [ -f ~/.rootpasswd ];
	then echo ".rootpasswd file already created"
	else
	PASSWORD=($(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 80 |\
	head -n 1))
	echo "root:${PASSWORD[0]}|chpasswd"
	echo "${PASSWORD[0]}" > .rootpasswd
	echo "created .rootpasswd in case you mess up"
fi
fi
# that sums up the initial prep, now for deploying the software
