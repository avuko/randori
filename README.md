# randori: fend off multiple attackers

## Fully based on PAM (*P*wn *A*ll *M*alware)

<!-- ![randori](randori.gif) -->

<img src="./randori.gif" align="left" />

> Randori (乱取り) is a form of practice<br />
> in which a designated aikidoka<br />
> defends against multiple attackers<br />
> in quick succession.<br />
> [https://en.wikipedia.org/wiki/Randori](https://en.wikipedia.org/wiki/Randori)

Basically it is my http://github.com/avuko/aiki PoC on steroids.

Shoutout to `0xBF` ([ONSec-Lab](https://github.com/ONsec-Lab/scripts/tree/master/pam_steal))
for giving us pam_steal. Everything here started with that one simple, great idea.

Also thanks to [micheloosterhof](https://twitter.com/micheloosterhof) for being
approachable when I had questions and comments about [cowrie](https://github.com/micheloosterhof/cowrie).


All of the steps to build a randori system are in a [prep & deploy](https://github.com/avuko/randori/tree/master/deploy) script.
The goal is to end up with a script you can point at an Ubuntu based system
to which you have root SSH access (with a key, not with a password)
and have a working randori system when the script is done.
You could also read the deploy.sh script to see which steps are required.

Currently, just before log rotation, the /root/results.sh script will run to get
all the relevant data from /var/log/syslog and /var/log/auth.log.
I'm going to improve this in the near future with a more useful way
of logging. I have simply not settled on a good way to do this, suggestions
are more than welcome.

## Running randori

Running is a simply a matter of:

```bash
screen -S randori
./startup.sh
```

If you just want to log usernames and passwords used and not "authenticate back", don't run `startup.sh` and
simply monitor `/root/randori-clients.log` and `/var/log/randori.log`

## Analysis

I will be adding tools I use for monitoring and analysis soon.
You can already have a look at the [Kathe](https://github.com/avuko/kathe)
tool I used for my [hack.lu](https://www.youtube.com/watch?v=-i1cyxTa8AM) presentation (it is multi-purpose).

## Next steps

 - Add more protocols
 - Fix the ugly logging issues
 - Add "scan back" functionality

Additionally, the general idea behind this (use existing services)
will be used to build a honeypot for monitoring DDoS activities.

## Disclaimer

This software has been created purely for the purposes of academic research
and for the development of effective defensive techniques, and is not
intended to be used to attack systems except where explicitly authorized.
Project maintainers are not responsible or liable for misuse of the
software. Use responsibly.

## Note from the developer

Nothing can stop you from using this tool -or tools like this- for
evil/stupid. Please consider that those "boxes" you hope to be "popping"
when you use any of the tools infosec provides, belong to real people.

Have some compassion. Or if that is too complicated, just don't be an asshat.
