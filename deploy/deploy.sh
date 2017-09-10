#!/usr/bin/env bash
set -e
set -u

target="${1}"

if [ "grep -q prep.sh \"AREYOUSURE='no'\"" ];then
   echo ""
   echo "You have not set the buildscript to set a"
   echo "password for root on the randori system."
   echo "If you are redeploying, this is not a problem."
   echo "Otherwise, please change AREYOURSURE='no' to"
   echo "AREYOUSURE='yes' in prep.sh"
   echo ""
   echo "Continue with setting a root password? (yes/no)"

   read continueinstall

   if [ "${continueinstall}" == 'yes' ]; then
	echo "Continuing deployment"
   else
	exit
   fi
fi

opensshdir='openssh-7.2p2'

copy_over=(rsyslog prep.sh make.sh pam_randori.c common-auth startup.sh killall.sh\
 randorifan-linux-amd64 randorisink-linux-amd64 torissh-linux-amd64 toritelnet-linux-amd64\
 rsyslog rsyslog.conf)

for co in ${copy_over[@]}; do
 echo $co;
 if [ $co == 'common-auth' ]; then
	scp $co ${target}:/etc/pam.d/
elif [ $co == 'killall.sh' ]; then
	scp $co ${target}:.
	# run it to stop any previous installations
	ssh ${target} './killall.sh'
elif [ $co == 'rsyslog' ]; then
	echo '!! WARNING: please note I am setting your logfiles up to'
	echo 'pile up for a full year. This could fill up your disks.'
	scp $co ${target}:/etc/logrotate.d/rsyslog
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

# installing ftp daemon
ssh "${target}" 'apt-get -y install vsftpd'
scp vsftpd.conf "${target}:/etc/"
ssh ${target} '/etc/init.d/vsftpd restart'

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
