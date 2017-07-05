#!/bin/sh

set -e

rm -f pam_randori.so
gcc -g -O2 -MT pam_randori_la-pam_randori.lo -MD -MP -MF\
pam_randori_la-pam_randori.Tpo -c pam_randori.c  -fPIC -DPIC -o\
pam_randori_la-pam_randori.o
gcc -shared pam_randori_la-pam_randori.o -lpam_misc -lpam -Wl,\
-soname -Wl,pam_randori.so -o pam_randori.so
rm -f pam_randori_la-pam_randori.Tpo pam_randori_la-pam_randori.o

cp pam_randori.so /lib/x86_64-linux-gnu/security/
