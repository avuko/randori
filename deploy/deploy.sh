#!/usr/bin/env bash
# EXIT='yes'

set -e
set -u


target="${1}"

if [[ $(grep -q "AREYOUSURE='no'" prep.sh) ]];then
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
   else
   echo "prep.sh tells me you are sure you have"
   echo "ssh key access to the remote system"
   echo "continuing"
   sleep 1
fi

# for now, just update
ssh "${target}" 'apt-get -y update'

# update and upgrade
# ssh "${target}" 'apt-get -y update && apt-get -y upgrade'
#if [ $EXIT == 'yes' ];then
#  echo -e "\nBecause of some interaction required"
#  echo "(a digitalocean issue) I cannot run"
#  echo "update & upgrade. Please log in to"
#  echo "the remote system, run:"
#  echo -e '\napt-get -y update && apt-get -y upgrade\n'
#  echo "and then change to EXIT='no' in the ${0}"
#  echo "script. Sorry for the inconvenience."
#  exit
#fi

# prepare for golang install
ssh "${target}" 'mkdir -p ~/work'

opensshdir='openssh-7.2p2'

copy_over=(rsyslog prep.sh make.sh pam_randori.c common-auth startup.sh killall.sh \
 randorifan.go randorisink.go torissh.go toritelnet.go \
 rsyslog rsyslog.conf login.defs logrotate.conf results.sh)

for co in ${copy_over[@]}; do
 # echo $co;
 if [ $co == 'common-auth' ]; then
	scp $co ${target}:/etc/pam.d/
elif [ $co == 'killall.sh' ]; then
	scp $co ${target}:.
	# run it to stop any previous installations
	echo "Making sure there are no previous randori versions running"
	set +e
	ssh ${target} './killall.sh'
	set -e
elif [ $co == 'rsyslog' ]; then
	echo '!! WARNING: please note I am setting your logfiles up to'
	echo 'receive the results. Please verify the changes I make.'
	scp $co ${target}:/etc/logrotate.d/rsyslog
elif [ $co == 'rsyslog.conf' ]; then
	scp $co ${target}:/etc/rsyslog.conf
elif [ $co == 'login.defs' ]; then
	scp $co ${target}:/etc/login.defs
elif [ $co == 'logrotate.conf' ]; then
	scp $co ${target}:/etc/logrotate.conf
elif [[ $co == *.go ]]; then
	scp $co ${target}:work/
else
	scp $co ${target}:.
fi
done

ssh "${target}" '/etc/init.d/rsyslog restart'
ssh "${target}" 'apt-get -y install build-essential libpam0g-dev telnetd\
 dpatch fakeroot devscripts equivs lintian quilt dpkg-dev dh-autoreconf\
 dh-systemd libzmq-dev pkg-config'

# python-zmq libzmq5 python-pip

# ssh "${target}" 'pip install --upgrade pip'
# ssh "${target}" 'pip install tailer'

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
scp xinetd.conf "${target}:/etc/xinetd.conf"

# installing ftp daemon
# ssh "${target}" 'apt-get -y install vsftpd'
# scp vsftpd.conf "${target}:/etc/"
# ssh ${target} '/etc/init.d/vsftpd restart'

# XXX CAREFUL, hardcoded search/replace
ssh "${target}" "sed -i 's/^telnet/#telnet/g' /etc/inetd.conf"
ssh "${target}" 'systemctl restart xinetd.service'
# installing openssh
ssh "${target}" 'apt-get -y source openssh'

# XXX CAREFUL, fixed version number
scp "auth-pam.c" "${target}:${opensshdir}/"
ssh "${target}" 'apt-get -y build-dep openssh'
ssh "${target}" "cd ${opensshdir} && fakeroot debian/rules clean && fakeroot debian/rules binary"
ssh "${target}" 'dpkg --install --force-all openssh-server_*'
# we need to add this configuration to disable reverse lookups and allow password
# logins
scp sshd_config "${target}:/etc/ssh/"
ssh "${target}" 'systemctl restart sshd.service'

# getting go, building randori
ssh "${target}" 'curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz'
ssh "${target}" 'tar -xvf go1.8.linux-amd64.tar.gz && mv go /usr/local'
ssh "${target}" 'echo '\''export PATH=$PATH:/usr/local/go/bin'\'' >> ~/.profile'
ssh "${target}" 'echo '\''export GOPATH=$HOME/work'\'' >> ~/.profile'
ssh "${target}" 'export PATH=$PATH:/usr/local/go/bin GOPATH=$HOME/work ;go get github.com/alecthomas/gozmq'
ssh "${target}" 'export PATH=$PATH:/usr/local/go/bin GOPATH=$HOME/work ;go get golang.org/x/crypto/ssh'
# ssh "${target}" 'export PATH=$PATH:/usr/local/go/bin GOPATH=$HOME/work go get github.com/go-redis/redis'
ssh "${target}" 'sed -i '\''s/arcfour128",/arcfour128", "aes128-cbc",/'\'' /root/work/src/golang.org/x/crypto/ssh/common.go'
ssh "${target}" 'export PATH=$PATH:/usr/local/go/bin GOPATH=$HOME/work ;go build work/randorifan.go && go build work/randorisink.go'
ssh "${target}" 'export PATH=$PATH:/usr/local/go/bin GOPATH=$HOME/work ;go build work/toritelnet.go && go build work/torissh.go'
