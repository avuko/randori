#!/usr/bin/env bash
TARGET="${1}"
echo "TELNET"

creds=(test1:test1 儿童游戏:chinese\
	chinese:儿童游戏\
       	test1:badpw\
       	baduser:test1\
       	:emptyuser emptypass:)
for cred in ${creds[@]};do
	username="${cred%:*}"
	password="${cred#*:}"
	echo ${username}:${password}
	(sleep 1; echo "${username}"; sleep 1; echo "${password}"; sleep 1)|\
		telnet ${TARGET}
done


echo SSH
echo "success:"
sshpass -p 'test1' ssh "test1@${TARGET}"
echo "fail, bad pw"
sshpass -p 'bad pw' ssh "test1@${TARGET}"
echo "fail, user=chinese UTF-8"
sshpass -p 'chinese' ssh "儿童游戏@${TARGET}"
echo "fail, pw=chinese UTF-8"
sshpass -p '儿童游戏' ssh "chinese@${TARGET}"
echo "fail, empty pw"
sshpass -p '' ssh "emptypw@${TARGET}"
