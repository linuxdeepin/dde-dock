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

#include <stdio.h>
#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/record.h>
#include <string.h>

#include "xrecord.h"
#include "_cgo_export.h"

static xcb_connection_t *data_disp = NULL;
static xcb_connection_t *ctrl_disp = NULL;
static xcb_record_context_t rc = 0;

// stop flag
static int stop = 0;

void event_callback(uint8_t * data);

int xrecord_grab_init()
{
	ctrl_disp = xcb_connect(NULL, NULL);
	data_disp = xcb_connect(NULL, NULL);

	if (xcb_connection_has_error(ctrl_disp)
	    || xcb_connection_has_error(data_disp)) {
		/*fprintf(stderr, "Error to open local display!\n"); */
		return 1;
	}
	const xcb_query_extension_reply_t *query_ext =
	    xcb_get_extension_data(ctrl_disp,
				   &xcb_record_id);

	if (!query_ext) {
		/*fprintf(stderr, "RECORD extension not supported on this X server!\n"); */
		return 2;
	}

	rc = xcb_generate_id(ctrl_disp);

	xcb_record_client_spec_t rcs;
	rcs = XCB_RECORD_CS_ALL_CLIENTS;

	xcb_record_range_t rr;
	memset(&rr, 0, sizeof(rr));
	rr.device_events.first = XCB_KEY_PRESS;
	rr.device_events.last = XCB_BUTTON_RELEASE;

	xcb_void_cookie_t create_cookie =
	    xcb_record_create_context_checked(ctrl_disp, rc, 0, 1, 1, &rcs,
					      &rr);
	xcb_generic_error_t *error =
	    xcb_request_check(ctrl_disp, create_cookie);
	if (error) {
		/*fprintf (stderr, "Could not create a record context!\n"); */
		free(error);
		return 3;
	}

	return 0;
}

void xrecord_grab_event_loop_start()
{
	xcb_record_enable_context_cookie_t cookie =
	    xcb_record_enable_context(data_disp, rc);

	while (!stop && data_disp != NULL) {
		xcb_record_enable_context_reply_t *reply =
		    xcb_record_enable_context_reply(data_disp, cookie, NULL);
		if (!reply)
			break;
		if (reply->client_swapped) {
			/*fprintf (stderr, "I am too lazy to implement byteswapping\n"); */
			return;
		}

		if (reply->category == 0 /* XRecordFromServer */ ) {
			uint8_t *data = xcb_record_enable_context_data(reply);
			int data_length =
			    xcb_record_enable_context_data_length(reply);
			/*printf("data length is %d\n", data_length); */

			if (data_length == sizeof(xcb_key_press_event_t)) {
				event_callback(data);
			}
		}
		free(reply);
	}
}

void xrecord_grab_finalize()
{
	stop = 1;
	xcb_record_disable_context(ctrl_disp, rc);
	xcb_record_free_context(ctrl_disp, rc);
	xcb_flush(ctrl_disp);

	xcb_disconnect(data_disp);
	xcb_disconnect(ctrl_disp);
	data_disp = NULL;
	ctrl_disp = NULL;
}

void event_callback(uint8_t * data)
{
	uint8_t event_type = data[0];
	switch (event_type) {
	case XCB_KEY_PRESS:{
			xcb_key_press_event_t *ev =
			    (xcb_key_press_event_t *) data;
			handleKeyEvent(1, ev->detail, ev->state);
			break;
		}
	case XCB_KEY_RELEASE:{
			xcb_key_release_event_t *ev =
			    (xcb_key_release_event_t *) data;
			handleKeyEvent(0, ev->detail, ev->state);
			break;
		}
	case XCB_BUTTON_PRESS:
		handleButtonEvent(1);
		break;
	case XCB_BUTTON_RELEASE:
		handleButtonEvent(0);
		break;
	}
}
