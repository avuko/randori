/*
 * pam_randori - get remote service/clientip/username/password from
 * brute-force attacks
 *
 * Usage: add "auth required pam_randori.so"
 * into /etc/pam.d/common-auth
 * just above "auth requisite pam_deny.so"
 *
 *
 * Reload services using PAM to start getting output.
 * Perhaps needless to add, but you might want to
 * only log in with keys :)
 */

#include <stdio.h>
#include <string.h>
#include <time.h>
#include <security/pam_modules.h>
#define LOGFILE "/var/log/randori.log"

PAM_EXTERN int pam_sm_authenticate(pam_handle_t * pamh, int flags
                                   ,int argc, const char **argv)
{
    int retval;

    const void *servicename;
    const char *username;
    const void *password;
    const void *rhostname;
    FILE *log;

    time_t	now;
    struct tm	ts;
    char	timestamp[80];


    time(&now);
    // Format time, "yyyy-mm-ddThh:mm:ss+zzzz"
    ts = *localtime(&now);
    strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%S%z", &ts);

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
    fprintf(log, "%s\u2002%s\u2002%s\u2002%s\u2002%s\u2002\n", (char *) timestamp,
    		    (char *) servicename, (char *) rhostname,
    		    (char *) username, (char *) password);
    fclose( log);

    return PAM_IGNORE;

    //}
}


PAM_EXTERN int pam_sm_setcred(pam_handle_t *pamh, int flags,
                              int argc, const char **argv)
{
    return PAM_IGNORE;
}
