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


All of the steps to build a randori system are (will be) in a [prep & deploy](https://github.com/avuko/randori/tree/master/deploy) script.
The goal is to end up with a script you can point at an Ubuntu based system 
and have a running randori system when it is done. Or read to see which steps are required.
Because I hate chasing down undocumented dependencies at 2 in the morning.

This will all log to something nice like `ELK` or `sqlite3` by the time I'm done.


## PS

I'm still thinking of the best way to patch the golang ssh client library.
You'll need to change <godir>src/golang.org/x/crypto/ssh/common.go to contain:

```diff
var supportedCiphers = []string{
 	"aes128-ctr", "aes192-ctr", "aes256-ctr",
 	"aes128-gcm@openssh.com",
-	"arcfour256", "arcfour128",
+	"arcfour256", "arcfour128", "aes128-cbc",
 }
```

