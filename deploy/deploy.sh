#!/usr/bin/env bash
set -e
set -u

target='testbed'
opensshdir='openssh-7.2p2'

copy_over=(prep.sh make.sh pam_randori.c common-auth)
for co in ${copy_over[@]}; do
 echo $co;
 if [ $co == 'common-auth' ]; then
	scp $co ${target}:/etc/pam.d/
else
	scp $co ${target}:.
fi
done
ssh "${target}" 'apt-get -y install build-essential libpam0g-dev telnetd\
 dpatch fakeroot devscripts equivs lintian quilt dpkg-dev dh-autoreconf\
 dh-systemd'

# setting a couple of limits right
ssh "${target}" './prep.sh'

# making sure every module can write to randori.log
ssh "${target}" 'touch /var/log/randori.log && chmod a+rw /var/log/randori.log'

# making the pam module
ssh "${target}" './make.sh'

# installing telnet daemon
ssh "${target}" 'apt-get -y install xinetd telnetd'
# TODO This is hoping it does not add an entry to /etc/inetd.conf
# we need to add this to disable reverse dns lookups
scp telnet "${target}:/etc/xinetd.d/"

# installing openssh
ssh "${target}" 'apt source openssh'

# XXX CAREFUL, fixed version number
scp "auth-pam.c" "${target}:${opensshdir}/"
ssh "${target}" 'apt-get build-dep openssh'
ssh "${target}" "cd ${opensshdir} && fakeroot debian/rules clean && fakeroot debian/rules binary"
ssh "${target}" 'dpkg --install --force-all openssh-server_*'
# we need to add this to disable reverse lookups
scp sshd_config "${target}:/etc/ssh/"
ssh "${target}" 'systemctl restart sshd.service'
