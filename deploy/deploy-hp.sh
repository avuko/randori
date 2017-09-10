#!/usr/bin/env bash
set -e
set -u

target="${1}"
opensshdir='openssh-7.2p2'

copy_over=(prep-hp.sh make.sh pam_randori.c common-auth rsyslog rsyslog.conf)

for co in ${copy_over[@]}; do
 echo $co;
 if [ $co == 'common-auth' ]; then
	scp $co ${target}:/etc/pam.d/
elif [ $co == 'rsyslog' ]; then
	echo "!! WARNING: please note I'm setting your logfiles up to"
	echo "pile up for a full year. This could fill up your disks."
	scp $co ${target}:/etc/logrotate.d/rsyslog
elif [ $co == 'rsyslog.conf' ]; then
	scp $co ${target}:/etc/rsyslog.conf
	ssh ${target} '/etc/init.d/rsyslog restart'
else
	scp $co ${target}:.
fi
done
ssh "${target}" 'apt-get update && apt-get upgrade'
ssh "${target}" 'apt-get -y install build-essential telnetd\
 dpatch fakeroot devscripts equivs lintian quilt dpkg-dev dh-autoreconf\
 dh-systemd libpam0g-dev'


# setting a couple of limits right and changing root pw
ssh "${target}" './prep-hp.sh'

scp "${target}:.rootpasswd" "${target}-rootpasswd"
# making sure every module can write to randori.log
ssh "${target}" 'touch /var/log/randori.log && chmod a+rw /var/log/randori.log'

# making the pam module
ssh "${target}" './make.sh'

# installing telnet daemon
ssh "${target}" 'apt-get -y install xinetd telnetd'
# we need to add this configuration to disable reverse dns lookups
scp telnet "${target}:/etc/xinetd.d/"


# XXX CAREFUL, hardcoded replace
ssh "${target}" "sed -i 's/^telnet/#telnet/g' /etc/inetd.conf"
ssh "${target}" 'systemctl restart xinetd.service'
# installing openssh
ssh "${target}" 'apt source openssh'

# XXX CAREFUL, fixed version number
scp "auth-pam.c" "${target}:${opensshdir}/"
ssh "${target}" 'apt-get build-dep openssh'
ssh "${target}" "cd ${opensshdir} && fakeroot debian/rules clean && fakeroot debian/rules binary"
ssh "${target}" 'dpkg --install --force-all openssh-server_*'
# we need to add this configuration to disable reverse lookups and allow password
# logins
scp sshd_config "${target}:/etc/ssh/"
ssh "${target}" 'systemctl restart sshd.service'
