// SPDX-FileCopyrightText: 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "trashhelper.h"

static void delete_trash_file (GFile *file, gboolean del_file, gboolean del_children);

TrashHelper::TrashHelper(QObject *parent)
    : QObject(parent)
    , m_trash(g_file_new_for_uri("trash:///"))
    , m_trashMonitor(g_file_monitor_file(m_trash, G_FILE_MONITOR_NONE, NULL, NULL))
{
    g_signal_connect(m_trashMonitor, "changed", G_CALLBACK(slot_onTrashMonitorChanged), this);
}

TrashHelper::~TrashHelper()
{
    g_object_unref(m_trashMonitor);
    g_object_unref(m_trash);
}

int TrashHelper::trashItemCount()
{
    GFileInfo *info;
    gint file_count = 0;

    info = g_file_query_info(m_trash, G_FILE_ATTRIBUTE_TRASH_ITEM_COUNT, G_FILE_QUERY_INFO_NONE, NULL, NULL);
    if (info != NULL) {
        file_count = g_file_info_get_attribute_uint32(info, G_FILE_ATTRIBUTE_TRASH_ITEM_COUNT);
        g_object_unref(info);
    }

    return file_count;
}

void TrashHelper::onTrashMonitorChanged(GFileMonitor *monitor, GFile *file, GFile *other_file, GFileMonitorEvent event_type)
{
    Q_UNUSED(monitor)
    Q_UNUSED(file)
    Q_UNUSED(other_file)

    if (event_type == G_FILE_MONITOR_EVENT_ATTRIBUTE_CHANGED) {
        Q_EMIT trashAttributeChanged();
    }
}

void TrashHelper::slot_onTrashMonitorChanged(GFileMonitor *monitor, GFile *file,
                                         GFile *other_file, GFileMonitorEvent event_type,
                                         gpointer user_data)
{
    TrashHelper * that = reinterpret_cast<TrashHelper*>(user_data);
    that->onTrashMonitorChanged(monitor, file, other_file, event_type);
}

bool TrashHelper::emptyTrash()
{
    delete_trash_file(m_trash, false, true);
    return true;
}

// gio code
static void
delete_trash_file (GFile *file, gboolean del_file, gboolean del_children)
{
  GFileInfo *info;
  GFile *child;
  GFileEnumerator *enumerator;

  g_return_if_fail (g_file_has_uri_scheme (file, "trash"));

  if (del_children)
    {
      enumerator = g_file_enumerate_children (file,
                                              G_FILE_ATTRIBUTE_STANDARD_NAME ","
                                              G_FILE_ATTRIBUTE_STANDARD_TYPE,
                                              G_FILE_QUERY_INFO_NOFOLLOW_SYMLINKS,
                                              NULL,
                                              NULL);
      if (enumerator)
        {
          while ((info = g_file_enumerator_next_file (enumerator, NULL, NULL)) != NULL)
            {
              child = g_file_get_child (file, g_file_info_get_name (info));

              /* The g_file_delete operation works differently for locations
               * provided by the trash backend as it prevents modifications of
               * trashed items. For that reason, it is enough to call
               * g_file_delete on top-level items only.
               */
              delete_trash_file (child, TRUE, FALSE);

              g_object_unref (child);
              g_object_unref (info);
            }
          g_file_enumerator_close (enumerator, NULL, NULL);
          g_object_unref (enumerator);
        }
    }

  if (del_file)
    g_file_delete (file, NULL, NULL);
}
