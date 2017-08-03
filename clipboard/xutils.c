/*
 * Copyright © 2004 Red Hat, Inc.
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notice appear in all copies and that both that
 * copyright notice and this permission notice appear in supporting
 * documentation, and that the name of Red Hat not be used in advertising or
 * publicity pertaining to distribution of the software without specific,
 * written prior permission.  Red Hat makes no representations about the
 * suitability of this software for any purpose.  It is provided "as is"
 * without express or implied warranty.
 *
 * RED HAT DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE, INCLUDING ALL
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO EVENT SHALL RED HAT
 * BE LIABLE FOR ANY SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
 * OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN 
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 * Author:  Matthias Clasen, Red Hat, Inc.
 */

#include <stdlib.h>

#include "xutils.h"

Atom XA_ATOM_PAIR;
Atom XA_CLIPBOARD_MANAGER;
Atom XA_CLIPBOARD;
Atom XA_DELETE;
Atom XA_INCR;
Atom XA_INSERT_PROPERTY;
Atom XA_INSERT_SELECTION;
Atom XA_MANAGER;
Atom XA_MULTIPLE;
Atom XA_NULL;
Atom XA_SAVE_TARGETS;
Atom XA_TARGETS;
Atom XA_TIMESTAMP;

unsigned long SELECTION_MAX_SIZE = 0;


void
init_atoms (Display *display)
{
  unsigned long max_request_size;
  
  if (SELECTION_MAX_SIZE > 0)
    return;

  XA_ATOM_PAIR = XInternAtom (display, "ATOM_PAIR", False);
  XA_CLIPBOARD_MANAGER = XInternAtom (display, "CLIPBOARD_MANAGER", False);
  XA_CLIPBOARD = XInternAtom (display, "CLIPBOARD", False);
  XA_DELETE = XInternAtom (display, "DELETE", False);
  XA_INCR = XInternAtom (display, "INCR", False);
  XA_INSERT_PROPERTY = XInternAtom (display, "INSERT_PROPERTY", False);
  XA_INSERT_SELECTION = XInternAtom (display, "INSERT_SELECTION", False);
  XA_MANAGER = XInternAtom (display, "MANAGER", False);
  XA_MULTIPLE = XInternAtom (display, "MULTIPLE", False);
  XA_NULL = XInternAtom (display, "NULL", False);
  XA_SAVE_TARGETS = XInternAtom (display, "SAVE_TARGETS", False);
  XA_TARGETS = XInternAtom (display, "TARGETS", False);
  XA_TIMESTAMP = XInternAtom (display, "TIMESTAMP", False);
  
  max_request_size = XExtendedMaxRequestSize (display);
  if (max_request_size == 0)
    max_request_size = XMaxRequestSize (display);
  
  SELECTION_MAX_SIZE = max_request_size - 100;
  if (SELECTION_MAX_SIZE > 262144)
    SELECTION_MAX_SIZE =  262144;
}

typedef struct 
{
  Window window;
  Atom timestamp_prop_atom;
} TimeStampInfo;

static Bool
timestamp_predicate (Display *display,
		     XEvent  *xevent,
		     XPointer arg)
{
  TimeStampInfo *info = (TimeStampInfo *)arg;

  if (xevent->type == PropertyNotify &&
      xevent->xproperty.window == info->window &&
      xevent->xproperty.atom == info->timestamp_prop_atom)
    return True;

  return False;
}

Time
get_server_time (Display *display,
		 Window   window)
{
  unsigned char c = 'a';
  XEvent xevent;
  TimeStampInfo info;

  info.timestamp_prop_atom = XInternAtom  (display, "_TIMESTAMP_PROP", False);
  info.window = window;

  XChangeProperty (display, window,
		   info.timestamp_prop_atom, info.timestamp_prop_atom,
		   8, PropModeReplace, &c, 1);

  XIfEvent (display, &xevent,
	    timestamp_predicate, (XPointer)&info);

  return xevent.xproperty.time;
}

