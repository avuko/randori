# Randori: Like Aiki. With a couple of Dans under its belt.

## Fully based on PAM. Because PAM is actually an acronym for "*P*wn *A*ll *M*achines"

Randori (乱取り) is a term used in Japanese martial arts to describe free-style
practice. In the Aikikai style of aikido, it refers to a form of practice in which
a designated aikidoka defends against multiple attackers in quick succession without
knowing how they will attack or in what order. [https://en.wikipedia.org/wiki/Randori]

Basically it is http://github.com/avuko/aiki on steroids.



First of all, shoutout to `0xBF` (ONSec-Lab) for giving us
https://github.com/ONsec-Lab/scripts/tree/master/pam_steal.
All of the below is based on that simple, great idea.

Also thanks to @micheloosterhof for being approachable when I had questions
and comments about cowrie (https://github.com/micheloosterhof/cowrie).

## PAM module

This PAM module will log to /var/log/aiki.log all services, remote hosts, usernames and
passwords. In a future release, all of this will be logged to a message queue for further
processing. For now, it is just a regular low-interaction honeypot gathering credentials.

```c
/*
 * pam_aiki - get remote service/clientip/username/password from
 * brute-force attacks
 *
 * Usage: add "auth required pam_aiki.so" into /etc/pam.d/common-auth
 * just above "auth requisite pam_deny.so"
 * 
 * 
 * Reload services using PAM to start getting output.
 * Perhaps needless to add, but you might want to 
 * only log in with keys :)
 */

#include <stdio.h>
#include <string.h>
#include <security/pam_modules.h>
#define LOGFILE "/var/log/aiki.log"

PAM_EXTERN int pam_sm_authenticate(pam_handle_t * pamh, int flags
                                   ,int argc, const char **argv)
{
    int retval;
    
    const void *servicename;
    const char *username;
    const void *password;
    const void *rhostname;
    FILE *log;
    
    /* get the name of the calling PAM_SERVICE. */
    retval=pam_get_item(pamh, PAM_SERVICE, &servicename);

    /* get the RHOST ip address. */
    retval=pam_get_item(pamh, PAM_RHOST, &rhostname);

    retval = pam_get_user(pamh, &username, NULL);

    retval = pam_get_item(pamh, PAM_AUTHTOK, &password);

    /* As opposed to the original pam_steal, I DO care about
     * non-existing user passwords.
     * Perhaps we should drop attempts without a password later
    */
    //if (password != NULL) {
    log = fopen (LOGFILE, "a");
    fprintf(log, "%s\u2002%s\u2002%s\u2002%s\n", (char *) servicename,
            (char *) rhostname, (char *) username, (char *) password);
    fclose( log);

    return PAM_IGNORE;

    //}
}


PAM_EXTERN int pam_sm_setcred(pam_handle_t *pamh, int flags,
                              int argc, const char **argv)
{
    return PAM_IGNORE;
}
```

Run `./make.sh`

*make.sh*

```bash
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
```

Add `pam_aiki.so` to `/etc/pam.d/common-auth` just above the
`pam_deny` module:

```diff
diff common-auth /etc/pam.d/common-auth
19,20d18
< # XXX adding pam_aiki
< auth 	required 			pam_aiki.so
26d23
<
```

`/etc/pam.d/common-auth` should now look something like this:

```bash
# here's the fallback if no module succeeds
# XXX adding pam_aiki
auth 	required 			pam_aiki.so
auth	requisite			pam_deny.so
[...]
```

## OpenSSH rebuild

Yes, this is necessary.  OpenSSH, instead of keeping the original password,
throws out a (rather haphazardly chosen?) string:

```bash
grep -nR INCORRECT auth-pam.c
822:	/* t char junk[] = "\b\n\r\177INCORRECT"; */
```

In order not to mess with the original code to much (and because I'm already way
out of my comfort zone writing/editing all this C), I made a simple change only:

```c
	/* XXX avuko: 2017-19-06T17:00:00 Tweak to return the password
     * entered for a non existing account */

	/* t char junk[] = "\b\n\r\177INCORRECT"; */
	char *ret = NULL;
	size_t i, l = wire_password != NULL ? strlen(wire_password) : 0;

	if (l >= INT_MAX)
		fatal("%s: password length too long: %zu", __func__, l);

	ret = malloc(l + 1);
	for (i = 0; i < l; i++)
		ret[i] = wire_password[i % (sizeof(wire_password) - 1)];
		/* ret[i] = junk[i % (sizeof(junk) - 1)]; */
	ret[i] = '\0';
	return ret;
}
```

Yes, that is very, very likely a fully unnecessary loop. But then there is
`CVE-2016-6210-2`, so lets leave well enough alone.

## Apache

get apache2 and the apache pam module

```bash
apt install apache2
sudo apt-get install libapache2-mod-authnz-external pwauth
sudo apt-get install libapache2-mod-authz-unixgroup
sudo a2enmod authnz_external authz_unixgroup
```

Edit `/etc/apache2/mods-enabled/authnz_pam.conf`

```apache
<Location />
 AuthType Basic
 AuthName "admin"
 AuthBasicProvider PAM
 AuthPAMService apache
 Require pam-account apache
</Location>
```

Edit/create `/etc/pam.d/apache`

```email
@include common-auth
```

And that was all there was to it. Next, we need a lot of finishing touches with
`<VirtualHost>` magic.


