/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 8 -*-
 *
 * Copyright (C) 2007 William Jon McCann <mccann@jhu.edu>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 *
 */

#ifndef __GSD_CLIPBOARD_MANAGER_H
#define __GSD_CLIPBOARD_MANAGER_H

#include <glib-object.h>

G_BEGIN_DECLS

#define GSD_TYPE_CLIPBOARD_MANAGER         (gsd_clipboard_manager_get_type ())
#define GSD_CLIPBOARD_MANAGER(o)           (G_TYPE_CHECK_INSTANCE_CAST ((o), GSD_TYPE_CLIPBOARD_MANAGER, GsdClipboardManager))
#define GSD_CLIPBOARD_MANAGER_CLASS(k)     (G_TYPE_CHECK_CLASS_CAST((k), GSD_TYPE_CLIPBOARD_MANAGER, GsdClipboardManagerClass))
#define GSD_IS_CLIPBOARD_MANAGER(o)        (G_TYPE_CHECK_INSTANCE_TYPE ((o), GSD_TYPE_CLIPBOARD_MANAGER))
#define GSD_IS_CLIPBOARD_MANAGER_CLASS(k)  (G_TYPE_CHECK_CLASS_TYPE ((k), GSD_TYPE_CLIPBOARD_MANAGER))
#define GSD_CLIPBOARD_MANAGER_GET_CLASS(o) (G_TYPE_INSTANCE_GET_CLASS ((o), GSD_TYPE_CLIPBOARD_MANAGER, GsdClipboardManagerClass))

typedef struct GsdClipboardManagerPrivate GsdClipboardManagerPrivate;

typedef struct
{
        GObject                     parent;
        GsdClipboardManagerPrivate *priv;
} GsdClipboardManager;

typedef struct
{
        GObjectClass   parent_class;
} GsdClipboardManagerClass;

GType                   gsd_clipboard_manager_get_type            (void);

GsdClipboardManager *       gsd_clipboard_manager_new                 (void);
gboolean                gsd_clipboard_manager_start               (GsdClipboardManager *manager,
                                                               GError         **error);
void                    gsd_clipboard_manager_stop                (GsdClipboardManager *manager);

int start_clip_manager();
int stop_clip_manager();
G_END_DECLS

#endif /* __GSD_CLIPBOARD_MANAGER_H */
