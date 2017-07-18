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
# default limits are not up for the task of dealing with bots
if $(grep -q 'nofile[[:space:]]*500000' /etc/security/limits.conf);then
	echo "nofile already set"
else
cat >> /etc/security/limits.conf <<EOF
*         hard    nofile      500000
*         soft    nofile      500000
root      hard    nofile      500000
root      soft    nofile      500000
EOF
fi

if $(grep -q 'fs.file-max = 2097152' /etc/sysctl.conf);then
 echo "fs.file-max already set"
else
echo "fs.file-max = 2097152" >> /etc/sysctl.conf
fi

# that sums up the initial prep, now for deploying the software
