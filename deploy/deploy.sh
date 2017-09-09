#!/usr/bin/env bash
set -e
set -u

target="${1}"
opensshdir='openssh-7.2p2'

copy_over=(rsyslog prep.sh make.sh pam_randori.c common-auth startup.sh killall.sh\
 randorifan-linux-amd64 randorisink-linux-amd64 torissh-linux-amd64 toritelnet-linux-amd64)

for co in ${copy_over[@]}; do
 echo $co;
 if [ $co == 'common-auth' ]; then
	scp $co ${target}:/etc/pam.d/
elif [ $co == 'rsyslog' ]; then
	echo "!! WARNING: please note I'm setting your logfiles up to"
	echo "pile up for a full year. This could fill up your disks."
	scp $co ${target}:/etc/logrotate.d/rsyslog
	ssh ${target} '/etc/init.d/rsyslog restart'
elif [ $co == 'rsyslog.conf' ]; then
	scp $co ${target}:/etc/rsyslog.conf
	ssh ${target} '/etc/init.d/rsyslog restart'
else
	scp $co ${target}:.
fi
done
ssh "${target}" 'apt-get -y install build-essential libpam0g-dev telnetd\
 dpatch fakeroot devscripts equivs lintian quilt dpkg-dev dh-autoreconf\
 dh-systemd python-zmq libzmq5 python-pip libzmq-dev'

ssh "${target}" 'pip install --upgrade pip'
ssh "${target}" 'pip install tailer'


# setting a couple of limits right
ssh "${target}" './prep.sh'

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
