#include "appicon.h"

#undef signals
extern "C" {
  #include <gtk/gtk.h>
}
#define signals public

AppIcon::AppIcon(QWidget *parent,Qt::WindowFlags f) :
    QLabel(parent)
{
    this->setParent(parent);
    this->setWindowFlags(f);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->setAlignment(Qt::AlignCenter);
}

AppIcon::AppIcon(QString iconPath, QWidget *parent, Qt::WindowFlags f) :
    QLabel(parent)
{
    this->setParent(parent);
    this->setWindowFlags(f);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->setAlignment(Qt::AlignCenter);

    QString sysIconPath = getSysIcon(iconPath);
    QPixmap iconPixmap;
    if (sysIconPath != "")
    {
        iconPixmap.load(sysIconPath);
    }
    else
    {
        iconPixmap.load(iconPath);
    }
    this->setPixmap(iconPixmap);
}

void AppIcon::setIcon(const QString &iconPath)
{
    QString sysIconPath = getSysIcon(iconPath);
    QPixmap iconPixmap(this->width(),this->height());
    if (sysIconPath != "")
    {
        iconPixmap.load(sysIconPath);
    }
    else
    {
        iconPixmap.load(iconPath);
    }
    this->setPixmap(iconPixmap);
}

QString AppIcon::getSysIcon(const QString &iconName, int size)
{
    return "";
//    char *name = iconName.toUtf8().data();
//    GtkIconTheme* theme;

//    if (g_path_is_absolute(name))
//        return iconName;
//    g_return_val_if_fail(name != NULL, NULL);

//    int pic_name_len = strlen(name);
//    char* ext = strrchr(name, '.');
//    if (ext != NULL) {
//        if (g_ascii_strcasecmp(ext+1, "png") == 0 || g_ascii_strcasecmp(ext+1, "svg") == 0 || g_ascii_strcasecmp(ext+1, "jpg") == 0) {
//            pic_name_len = ext - name;
//            g_debug("Icon name should an absoulte path or an basename without extension");
//        }
//    }

//    char* pic_name = g_strndup(name, pic_name_len);
//    theme = gtk_icon_theme_get_default();

//    GtkIconInfo* info = gtk_icon_theme_lookup_icon(theme, pic_name, size, GTK_ICON_LOOKUP_GENERIC_FALLBACK);
//    g_free(pic_name);
//    if (info) {
//        char* path = g_strdup(gtk_icon_info_get_filename(info));
//#if GTK_MAJOR_VERSION >= 3
//        g_object_unref(info);
//#elif GTK_MAJOR_VERSION == 2
//        gtk_icon_info_free(info);
//#endif
//        return QString(path);
//    } else {
//        return NULL;
//    }
}
