/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <netlink/netlink.h>
#include <netlink/genl/genl.h>
#include <netlink/genl/family.h>
#include <netlink/genl/ctrl.h>

#include "nl80211.h"
#include "_cgo_export.h"

struct nl80211_state {
    struct nl_sock *socket;
    int nl80211_id;
};

static const char *ifmodes[NL80211_IFTYPE_MAX + 1] = {
    "unspecified",
    "IBSS",
    "managed",
    "AP",
    "AP/VLAN",
    "WDS",
    "monitor",
    "mesh point",
    "P2P-client",
    "P2P-GO",
    "P2P-device",
    "outside context of a BSS",
    "NAN",
};

static char modebuf[100];

const char *iftype_name(enum nl80211_iftype iftype)
{
    if (iftype <= NL80211_IFTYPE_MAX && ifmodes[iftype])
        return ifmodes[iftype];
    sprintf(modebuf, "Unknown mode (%d)", iftype);
    return modebuf;
}

static int nl80211_init(struct nl80211_state *state)
{
    int err;

    state->socket = nl_socket_alloc();
    if (!state->socket) {
        fprintf(stderr, "Failed to allocate netlink socket.\n");
        return -ENOMEM;
    }

    // size from iw
    nl_socket_set_buffer_size(state->socket, 8192, 8192);

    if (genl_connect(state->socket)) {
        fprintf(stderr, "Failed to connect to generic netlink.\n");
        err = -ENOLINK;
        goto out_socket_free;
    }


    state->nl80211_id = genl_ctrl_resolve(state->socket, "nl80211");
    if (state->nl80211_id < 0) {
        fprintf(stderr, "nl80211 not found.\n");
        err = -ENOENT;
        goto out_socket_free;
    }

    return 0;

out_socket_free:
    nl_socket_free(state->socket);
    return err;
}

static void nl80211_cleanup(struct nl80211_state *state)
{
    nl_socket_free(state->socket);
}

static int error_handler(struct sockaddr_nl *nla, struct nlmsgerr *err,
                         void *arg)
{
    int *ret = arg;
    *ret = err->error;
    return NL_STOP;
}

static int finish_handler(struct nl_msg *msg, void *arg)
{
    int *ret = arg;
    *ret = 0;
    return NL_SKIP;
}

static int ack_handler(struct nl_msg *msg, void *arg)
{
    int *ret = arg;
    *ret = 0;
    return NL_STOP;
}

static char *current_phy = NULL;

static int valid_handler(struct nl_msg *msg, void *arg)
{
    struct nlattr *tb_msg[NL80211_ATTR_MAX + 1];
    struct genlmsghdr *gnlh = nlmsg_data(nlmsg_hdr(msg));
    static int64_t phy_id = -1;
    int print_name = 1;

    nla_parse(tb_msg, NL80211_ATTR_MAX, genlmsg_attrdata(gnlh, 0),
              genlmsg_attrlen(gnlh, 0), NULL);

    if (tb_msg[NL80211_ATTR_WIPHY]) {
        if (nla_get_u32(tb_msg[NL80211_ATTR_WIPHY]) == phy_id) {
            print_name = 0;
        }
        phy_id = nla_get_u32(tb_msg[NL80211_ATTR_WIPHY]);
    }
    if (print_name && tb_msg[NL80211_ATTR_WIPHY_NAME]) {
        /* printf("Wiphy %s\n", nla_get_string(tb_msg[NL80211_ATTR_WIPHY_NAME])); */
        if (current_phy) {
            free(current_phy);
            current_phy = NULL;
        }
        current_phy = strdup(nla_get_string(tb_msg[NL80211_ATTR_WIPHY_NAME]));
    }

    struct nlattr *nl_mode;
    int rem_mode;
    if (tb_msg[NL80211_ATTR_SUPPORTED_IFTYPES]) {
        /* printf("\tSupported interface modes:\n"); */
        nla_for_each_nested(nl_mode, tb_msg[NL80211_ATTR_SUPPORTED_IFTYPES], rem_mode) {
            /* printf("\t\t * %s\n", iftype_name(nla_type(nl_mode))); */
            addWirelessInfo(current_phy, strdup(iftype_name(nla_type(nl_mode))));
        }
    }

    return NL_SKIP;
}

int
wireless_info_query()
{
    struct nl80211_state nlstate;
    int err = nl80211_init(&nlstate);
    if (err) {
        fprintf(stderr, "Failed to init nl\n");
        return -1;
    }

    struct nl_cb *cb = NULL;
    struct nl_msg *msg = NULL;

    msg = nlmsg_alloc();
    if (!msg) {
        fprintf(stderr, "Failed to allocate netlink message\n");
        goto out;
    }

    // NL_CB_DEBUG, NL_CB_DEFAULT
    cb = nl_cb_alloc(NL_CB_DEFAULT);
    if (!cb) {
        fprintf(stderr, "Failed to allocate netlink callbacks\n");
        goto out;
    }

    nl_cb_set(cb, NL_CB_VALID, NL_CB_CUSTOM, valid_handler, NULL);

    genlmsg_put(msg, 0, 0, nlstate.nl80211_id, 0, NLM_F_DUMP,
                NL80211_CMD_GET_WIPHY, 0);

    err = nl_send_auto_complete(nlstate.socket, msg);
    if (err < 0) {
        goto out;
    }
    err = 1;

    nl_cb_err(cb, NL_CB_CUSTOM, error_handler, &err);
    nl_cb_set(cb, NL_CB_FINISH, NL_CB_CUSTOM, finish_handler, &err);
    nl_cb_set(cb, NL_CB_ACK, NL_CB_CUSTOM, ack_handler, &err);

    while (err > 0) {
        nl_recvmsgs(nlstate.socket, cb);
    }

    if (current_phy) {
        free(current_phy);
        current_phy = NULL;
    }
out:
    if (cb) {
      nl_cb_put(cb);
    }

    if (msg) {
      nlmsg_free(msg);
    }

    nl80211_cleanup(&nlstate);
    return 0;
}
