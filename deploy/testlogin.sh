#!/usr/bin/env bash
TARGET='10.21.247.142'
echo "TELNET"
echo "success:"
(sleep 1; echo test1; sleep 1; echo test1; sleep 1;echo "result: $?")|\
telnet $TARGET ;echo "result: $?"
echo "fail, bad pw"
(sleep 1; echo chinese; sleep 1; echo "儿童游戏"; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, user=chinese UTF-8"
(sleep 1; echo "儿童游戏"; sleep 1; echo "chines"; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, pw=chinese UTF-8"
(sleep 1; echo test1; sleep 1; echo badpw; sleep 1;echo "result: $?")|\
telnet $TARGET ; echo "result: $?"
echo "fail, bad user"
(sleep 1; echo baduser; sleep 1; echo test1; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, bad pw/user"
(sleep 1; echo baduser; sleep 1; echo badpw; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, empty user"
(sleep 1; echo ""; sleep 1; echo emptyuser; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, empty pw"
(sleep 1; echo emptypw; sleep 1; echo ""; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, empty user/pw"
(sleep 1; echo ""; sleep 1; echo ""; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, cr"
(sleep 1; echo "cr"; sleep 1; echo -e "\r\rcr"; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
echo "fail, nl"
(sleep 1; echo "nl"; sleep 1; echo -e "\n\nnl"; sleep 1;)|\
telnet $TARGET ; echo "result: $?"
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
