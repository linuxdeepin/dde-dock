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
#ifndef X_UTILS_H
#define X_UTILS_H

#include <X11/Xlib.h>


extern Atom XA_ATOM_PAIR;
extern Atom XA_CLIPBOARD_MANAGER;
extern Atom XA_CLIPBOARD;
extern Atom XA_DELETE;
extern Atom XA_INCR;
extern Atom XA_INSERT_PROPERTY;
extern Atom XA_INSERT_SELECTION;
extern Atom XA_MANAGER;
extern Atom XA_MULTIPLE;
extern Atom XA_NULL;
extern Atom XA_SAVE_TARGETS;
extern Atom XA_TARGETS;
extern Atom XA_TIMESTAMP;

extern unsigned long SELECTION_MAX_SIZE;

void init_atoms      (Display *display);

Time get_server_time (Display *display,
		      Window   window);

#endif /* X_UTILS_H */
