#include <stdio.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>
#include <security/pam_ext.h>
#include <systemd/sd-bus.h>
#include <errno.h>
#include <strings.h>
#include <syslog.h>

#define AUTHORITY_DBUS_SERVICE "com.deepin.daemon.Authority"
#define AUTHORITY_DBUS_PATH "/com/deepin/daemon/Authority"
#define AUTHORITY_DBUS_INTERFACE "com.deepin.daemon.Authority"

static int
has_cookie(pam_handle_t *pamh, sd_bus *bus, const char *username, int *has)
{
    sd_bus_error err = SD_BUS_ERROR_NULL;
    sd_bus_message *reply = NULL;
    int ret = 0;

    ret = sd_bus_call_method(bus, AUTHORITY_DBUS_SERVICE, AUTHORITY_DBUS_PATH,
                             AUTHORITY_DBUS_INTERFACE,
                             "HasCookie", &err, &reply, "s", username);
    if (ret < 0) {
        pam_syslog(pamh, LOG_ERR, "fail to call HasCookie: %s, %s", err.name, err.message);
        goto finish;
    }

    ret = sd_bus_message_read_basic(reply, 'b', has);
    if (ret < 0) {
        pam_syslog(pamh, LOG_ERR, "fail to read HasCookie reply: %s", strerror(errno));
        goto finish;
    }

finish:
    sd_bus_error_free(&err);
    sd_bus_message_unref(reply);

    return ret < 0 ? 1 : 0;
}

static int
check_cookie(pam_handle_t *pamh, sd_bus *bus, const char *username,
    const char *cookie, int *result)
{
    sd_bus_error err = SD_BUS_ERROR_NULL;
    sd_bus_message *reply = NULL;
    int ret = 0;

    ret = sd_bus_call_method(bus, AUTHORITY_DBUS_SERVICE, AUTHORITY_DBUS_PATH,
                            AUTHORITY_DBUS_INTERFACE,
                             "CheckCookie", &err, &reply,
                             "ss", username, cookie);
    if (ret < 0) {
        pam_syslog(pamh, LOG_ERR, "fail to call CheckCookie: %s, %s", err.name, err.message);
        goto finish;
    }

    ret = sd_bus_message_read_basic(reply, 'b', result);
    if (ret < 0) {
        pam_syslog(pamh, LOG_ERR, "fail to read CheckCookie reply: %s", strerror(errno));
        goto finish;
    }

finish:
    sd_bus_error_free(&err);
    sd_bus_message_unref(reply);

    return ret < 0 ? 1 : 0;
}


PAM_EXTERN int
pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc,
                                   const char **argv) {

    int arg_idx;
    int debug = 0;
    for (arg_idx = 0; arg_idx < argc; arg_idx++) {
        if (strcasecmp(argv[arg_idx], "debug") == 0 ) {
            debug = 1;
        }
    }

    const char *username;
    int ret;
    ret = pam_get_user(pamh, &username, NULL);
    if (ret != PAM_SUCCESS) {
        return PAM_SERVICE_ERR;
    }

    // connect to the bus
    sd_bus *bus = NULL;
    ret = sd_bus_open_system(&bus);
    if (ret < 0) {
        pam_syslog(pamh, LOG_ERR, "failed to open system bus: %s", strerror(errno));
        return PAM_SERVICE_ERR;
    }

    int has;
    ret = has_cookie(pamh, bus, username, &has);
    if (ret != 0) {
        return PAM_SERVICE_ERR;
    }

    if (debug) {
        pam_syslog(pamh, LOG_DEBUG, "has_cookie: %d", has);
    }
    if (!has) {
        return PAM_AUTH_ERR;
    }

    const char *cookie;
    ret = pam_get_authtok(pamh, PAM_AUTHTOK, &cookie, NULL);
    if (ret != PAM_SUCCESS) {
        return PAM_SERVICE_ERR;
    }

    if (cookie == NULL) {
        return PAM_AUTH_ERR;
    }

    int check_result;
    ret = check_cookie(pamh, bus, username, cookie, &check_result);
    if (ret != 0) {
        return PAM_SERVICE_ERR;
    }

    if (debug) {
        pam_syslog(pamh, LOG_DEBUG, "check_result: %d", check_result);
    }
    if (check_result) {
        return PAM_SUCCESS;
    }
    return PAM_AUTH_ERR;
}

PAM_EXTERN int
pam_sm_setcred (pam_handle_t *pamh, int flags,
		int argc, const char **argv)
{
  return PAM_IGNORE;
}
