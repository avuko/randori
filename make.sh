#!/bin/sh

set -e

rm -f pam_aiki.so
gcc -g -O2 -MT pam_aiki_la-pam_aiki.lo -MD -MP -MF\
pam_aiki_la-pam_aiki.Tpo -c pam_aiki.c  -fPIC -DPIC -o\
pam_aiki_la-pam_aiki.o
gcc -shared pam_aiki_la-pam_aiki.o -lpam_misc -lpam -Wl,\
-soname -Wl,pam_aiki.so -o pam_aiki.so
rm -f pam_aiki_la-pam_aiki.Tpo pam_aiki_la-pam_aiki.o

cp pam_aiki.so /lib/x86_64-linux-gnu/security/
