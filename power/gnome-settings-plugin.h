/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 8 -*-
 *
 * Copyright (C) 2002-2005 Paolo Maggi
 * Copyright (C) 2007      William Jon McCann <mccann@jhu.edu>
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
 * Foundation, Inc., 59 Temple Place, Suite 330,
 * Boston, MA 02111-1307, USA.
 */

#ifndef __GNOME_SETTINGS_PLUGIN_H__
#define __GNOME_SETTINGS_PLUGIN_H__

#include <glib-object.h>
#include <gmodule.h>

G_BEGIN_DECLS
#define GNOME_TYPE_SETTINGS_PLUGIN              (gnome_settings_plugin_get_type())
#define GNOME_SETTINGS_PLUGIN(obj)              (G_TYPE_CHECK_INSTANCE_CAST((obj), GNOME_TYPE_SETTINGS_PLUGIN, GnomeSettingsPlugin))
#define GNOME_SETTINGS_PLUGIN_CLASS(klass)      (G_TYPE_CHECK_CLASS_CAST((klass),  GNOME_TYPE_SETTINGS_PLUGIN, GnomeSettingsPluginClass))
#define GNOME_IS_SETTINGS_PLUGIN(obj)           (G_TYPE_CHECK_INSTANCE_TYPE((obj), GNOME_TYPE_SETTINGS_PLUGIN))
#define GNOME_IS_SETTINGS_PLUGIN_CLASS(klass)   (G_TYPE_CHECK_CLASS_TYPE ((klass), GNOME_TYPE_SETTINGS_PLUGIN))
#define GNOME_SETTINGS_PLUGIN_GET_CLASS(obj)    (G_TYPE_INSTANCE_GET_CLASS((obj),  GNOME_TYPE_SETTINGS_PLUGIN, GnomeSettingsPluginClass))

typedef struct
{
        GObject parent;
} GnomeSettingsPlugin;

typedef struct
{
        GObjectClass parent_class;

        /* Virtual public methods */
        void            (*activate)                     (GnomeSettingsPlugin *plugin);
        void            (*deactivate)                   (GnomeSettingsPlugin *plugin);
} GnomeSettingsPluginClass;

GType            gnome_settings_plugin_get_type           (void) G_GNUC_CONST;

void             gnome_settings_plugin_activate           (GnomeSettingsPlugin *plugin);
void             gnome_settings_plugin_deactivate         (GnomeSettingsPlugin *plugin);

#define GSD_DBUS_NAME "org.gnome.SettingsDaemon"
#define GSD_DBUS_PATH "/org/gnome/SettingsDaemon"
#define GSD_DBUS_BASE_INTERFACE "org.gnome.SettingsDaemon"

/*
 * Utility macro used to register plugins
 *
 * use: GNOME_SETTINGS_PLUGIN_REGISTER (PluginName, plugin_name)
 */
#define GNOME_SETTINGS_PLUGIN_REGISTER(PluginName, plugin_name)                \
typedef struct {                                                               \
        PluginName##Manager *manager;                                          \
} PluginName##PluginPrivate;                                                   \
typedef struct {                                                               \
	GnomeSettingsPlugin    parent;                                         \
	PluginName##PluginPrivate *priv;                                       \
} PluginName##Plugin;                                                          \
typedef struct {                                                               \
	GnomeSettingsPluginClass parent_class;                                 \
} PluginName##PluginClass;                                                     \
GType plugin_name##_plugin_get_type (void) G_GNUC_CONST;                       \
G_MODULE_EXPORT GType register_gnome_settings_plugin (GTypeModule *module);    \
                                                                               \
        G_DEFINE_DYNAMIC_TYPE (PluginName##Plugin,                             \
                               plugin_name##_plugin,                           \
                               GNOME_TYPE_SETTINGS_PLUGIN)                     \
                                                                               \
G_MODULE_EXPORT GType                                                          \
register_gnome_settings_plugin (GTypeModule *type_module)                      \
{                                                                              \
        plugin_name##_plugin_register_type (type_module);                      \
                                                                               \
        return plugin_name##_plugin_get_type();                                \
}                                                                              \
                                                                               \
static void                                                                    \
plugin_name##_plugin_class_finalize (PluginName##PluginClass * plugin_name##_class) \
{                                                                              \
}                                                                              \
                                                                               \
static void                                                                    \
plugin_name##_plugin_init (PluginName##Plugin *plugin)                         \
{                                                                              \
        plugin->priv = G_TYPE_INSTANCE_GET_PRIVATE ((plugin),                  \
                plugin_name##_plugin_get_type(), PluginName##PluginPrivate);   \
        g_debug (#PluginName " initializing");                                 \
        plugin->priv->manager = plugin_name##_manager_new ();                  \
}                                                                              \
                                                                               \
static void                                                                    \
plugin_name##_plugin_finalize (GObject *object)                                \
{                                                                              \
        PluginName##Plugin *plugin;                                            \
        g_return_if_fail (object != NULL);                                     \
        g_return_if_fail (G_TYPE_CHECK_INSTANCE_TYPE (object, plugin_name##_plugin_get_type())); \
        g_debug ("PluginName## finalizing");                                   \
        plugin = G_TYPE_CHECK_INSTANCE_CAST ((object), plugin_name##_plugin_get_type(), PluginName##Plugin); \
        g_return_if_fail (plugin->priv != NULL);                               \
        if (plugin->priv->manager != NULL)                                     \
                g_object_unref (plugin->priv->manager);                        \
        G_OBJECT_CLASS (plugin_name##_plugin_parent_class)->finalize (object);        \
}                                                                              \
                                                                               \
static void                                                                    \
impl_activate (GnomeSettingsPlugin *plugin)                                    \
{                                                                              \
        GError *error = NULL;                                                  \
        PluginName##Plugin *plugin_cast;                                       \
        g_debug ("Activating %s plugin", G_STRINGIFY(plugin_name));            \
        plugin_cast = G_TYPE_CHECK_INSTANCE_CAST ((plugin), plugin_name##_plugin_get_type(), PluginName##Plugin); \
        if (!plugin_name##_manager_start (plugin_cast->priv->manager, &error)) { \
                g_warning ("Unable to start %s manager: %s", G_STRINGIFY(plugin_name), error->message); \
                g_error_free (error);                                          \
        }                                                                      \
}                                                                              \
                                                                               \
static void                                                                    \
impl_deactivate (GnomeSettingsPlugin *plugin)                                  \
{                                                                              \
        PluginName##Plugin *plugin_cast;                                       \
        plugin_cast = G_TYPE_CHECK_INSTANCE_CAST ((plugin), plugin_name##_plugin_get_type(), PluginName##Plugin); \
        g_debug ("Deactivating %s plugin", G_STRINGIFY (plugin_name));         \
        plugin_name##_manager_stop (plugin_cast->priv->manager); \
}                                                                              \
                                                                               \
static void                                                                    \
plugin_name##_plugin_class_init (PluginName##PluginClass *klass)               \
{                                                                              \
        GObjectClass           *object_class = G_OBJECT_CLASS (klass);         \
        GnomeSettingsPluginClass *plugin_class = GNOME_SETTINGS_PLUGIN_CLASS (klass); \
                                                                               \
        object_class->finalize = plugin_name##_plugin_finalize;                \
        plugin_class->activate = impl_activate;                                \
        plugin_class->deactivate = impl_deactivate;                            \
        g_type_class_add_private (klass, sizeof (PluginName##PluginPrivate));  \
}

G_END_DECLS

#endif  /* __GNOME_SETTINGS_PLUGIN_H__ */
